package rabbitmq

import (
	"fmt"
	"time"
	amqp "github.com/rabbitmq/amqp091-go"
	"context"
	"sync/atomic"
)

type PoolConfig struct {
	URL string
	ConnectionTimeout time.Duration
	MinConnections int
	MaxConnections int
}

func DefaultConfig(url string) *PoolConfig {
	return &PoolConfig{
		URL: url,
		ConnectionTimeout: 5*time.Second,
		MinConnections: 2,
		MaxConnections: 10,
	}
}

func (c *PoolConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("URL cannot be empty!")
	}
	if c.MinConnections < 1 {
		return fmt.Errorf("MinConnections cannot be less than 1")
	}
	if c.MaxConnections > 99 {
		return fmt.Errorf("MaxConnections cannot be more than 99")
	}
	if c.MaxConnections < c.MinConnections {
		return fmt.Errorf("MaxConnections cannot be less than MinConnections")
	}
	return nil
}

type ConnectionPool struct {
	config *PoolConfig
	pool chan *amqp.Connection
	currentSize atomic.Int32
}

type connectionResult struct {
	conn *amqp.Connection
	err error
}

func NewConnectionPool(config *PoolConfig) (*ConnectionPool, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	pool := &ConnectionPool{
		config: config,
		pool: make(chan *amqp.Connection, config.MaxConnections),
	}
	ctx := context.Background()
	for i := 0; i < config.MinConnections; i++ {
		conn, err := CreateConnection(ctx, config.URL, config.ConnectionTimeout)
		if err != nil {
			pool.CloseAll()
			return nil, fmt.Errorf("failed to create initial connection %d: %w", i, err)
		}
		pool.pool <- conn
		pool.currentSize.Add(1)
	}
	return pool, nil
}

func (p* ConnectionPool) CloseAll() {
	close(p.pool)

	for conn := range p.pool {
		conn.Close()
	}
}

func (p *ConnectionPool) Acquire(ctx context.Context) (*amqp.Connection, error) {
	select {
	case conn := <-p.pool:
		return conn, nil
	default:

	}

	currentSize := p.currentSize.Load()
	if currentSize < int32(p.config.MaxConnections) {
		newSize := p.currentSize.Add(1)
		if newSize > int32(p.config.MaxConnections) {
			p.currentSize.Add(-1)
		} else {
			conn, err := CreateConnection(ctx, p.config.URL, p.config.ConnectionTimeout)
			if err != nil {
				p.currentSize.Add(-1)
				return nil, err
			}
			return conn, nil
		}
	}
	select {
	case conn := <-p.pool:
			// Someone released a connection!
			return conn, nil
	case <-ctx.Done():
			// Context cancelled or timed out
			return nil, fmt.Errorf("acquire timeout: %w", ctx.Err())
	}
}

func (p *ConnectionPool) Release(conn *amqp.Connection) {
	if conn == nil {
		return
	}
	if conn.IsClosed() {
		p.currentSize.Add(-1)
		return
	}
	select {
	case p.pool <- conn:
	default:
		conn.Close()
		p.currentSize.Add(-1)
	}
}

func (p *ConnectionPool) GetCurrentSize() int32 {
	return p.currentSize.Load()
}

func CreateConnection(ctx context.Context, url string, timeout time.Duration) (*amqp.Connection, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resultChan := make(chan connectionResult, 1)

	go func() {
		conn, err := amqp.DialConfig(url, amqp.Config{
			Heartbeat: 10 * time.Second,
			Locale: "en_US",
		})
		resultChan <- connectionResult{
			conn: conn,
			err: err,
		}
	}()

	select {
	case res := <-resultChan:
		if res.err != nil {
			return nil, fmt.Errorf("Failed to obtain RabbitMQ Connection: %w", res.err)
		}
		return res.conn, nil
	case <-timeoutCtx.Done():
		return nil, fmt.Errorf("Timeout while waiting for RabbitmMQ Connection: %w", timeoutCtx.Err())
	}
}
