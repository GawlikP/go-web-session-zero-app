package pool

import (
	"context"
	"testing"
	"time"

	rbit "session-zero-app/pkg/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TestCoonnectionPool_AcquireAndChannel(t *testing.T) {
	config := rbit.DefaultConfig("amqp://admin:admin@localhost:5672/")
	pool, err := rbit.NewConnectionPool(config)
	if err != nil {
		t.Fatalf("Failed to create pool, %v", err)
	}
	defer pool.CloseAll()

	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("Failed to acquire connection: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	if ch.IsClosed() {
		t.Error("Channel should not be closed")
	}
	t.Log("Successfully acquired connection and opened channel!")
}

func TestConnectionPool_AcquireAndReuseConnection(t *testing.T) {
	ctx := context.Background()
	config := rbit.DefaultConfig("amqp://admin:admin@localhost:5672/")
	connections := []*amqp.Connection{}
	pool, err := rbit.NewConnectionPool(config)
	if err != nil {
		t.Fatalf("Failed to create pool, %v", err)
	}
	defer pool.CloseAll()
	for i := 0; i < config.MinConnections; i++ {
		conn, err := pool.Acquire(ctx)
		if err != nil {
			t.Errorf("Cannot acquire the connection from the pool: %v\n", err)
			return
		}
		connections = append(connections, conn)
	}
	// Accquire new connection that should reuse already existing one without
	// increassing current size
	sizeBefore := pool.GetCurrentSize()
	pool.Release(connections[0])
	reusedConn, err := pool.Acquire(ctx)
	if err != nil {
		t.Errorf("Cannot acquire the reused connection from the pool: %v\n", err)
		return
	}
	ch, err := reusedConn.Channel()
	if err != nil {
		t.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()
	sizeAfter := pool.GetCurrentSize()
	if sizeAfter != sizeBefore {
		t.Errorf("The size after receiving third connection, after releasing previous one\n was expected to not change the size of the pool")
	}
}

func TestConnectionPool_ShouldLimitAtMax(t *testing.T) {
	ctx := context.Background()
	config := rbit.DefaultConfig("amqp://admin:admin@localhost:5672/")
	connections := []*amqp.Connection{}
	pool, err := rbit.NewConnectionPool(config)
	if err != nil {
		t.Fatalf("Failed to create pool, %v", err)
	}
	defer pool.CloseAll()
	for i := 0; i < config.MaxConnections; i++ {
		conn, err := pool.Acquire(ctx)
		if err != nil {
			t.Errorf("Cannot acquire the connection from the pool: %v\n", err)
			return
		}
		connections = append(connections, conn)
	}
	for i := 0; i < 2; i++ {
		tctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		conn, err := pool.Acquire(tctx)
		cancel()
		if err == nil {
			t.Errorf("Additional Connection %d acquired! That's bad!", i)
			connections = append(connections, conn)
		}
	}
	if len(connections) > config.MaxConnections {
		t.Errorf("The number of current connections exceeds the max!")
	}
}
