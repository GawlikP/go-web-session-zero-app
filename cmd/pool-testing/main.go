package main

import (
	"fmt"
	"context"
	"time"
	rbit "session-zero-app/pkg/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	config := rbit.DefaultConfig("amqp://admin:admin@localhost:5672/")
	pool, err := rbit.NewConnectionPool(config)
	if err != nil {
		fmt.Printf("Failed to create pool: %v\n", err)
		return
	}
	defer pool.CloseAll()

	fmt.Printf("Testing acquire & release\n")
	testBasicAqquireRelease(pool)
	// fmt.Printf("Testing Reuse connections\n")
	// testReuseConnection(pool)
	fmt.Printf("Testing the max pool size\n")
	testMaxConnections(pool)
}

func testBasicAqquireRelease(p *rbit.ConnectionPool) {
	ctx := context.Background()
	conn, err := p.Acquire(ctx)
	if err != nil {
		fmt.Printf("Cannot acquire the connection from the pool: %v\n", err)
		return
	}
	defer p.Release(conn)
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Cannot receive the channel on a Rabbitmq connection: %v\n", err)
		return
	}
	ch.Close()
}

func testMaxConnections(p *rbit.ConnectionPool) {
	ctx := context.Background()
	connections := []*amqp.Connection{}
	for i := 0; i < 10; i++ {
		conn, err := p.Acquire(ctx)
		if err != nil {
			fmt.Printf("Cannot acquire the connection from the pool: %v\n", err)
			return
		}
		connections = append(connections, conn)
		fmt.Printf("Acquired connection :%d", i)
	}
	fmt.Printf("Acquired max connections 10/10")

	for i := 0; i < 2; i++ {
		tctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		conn, err := p.Acquire(tctx)
		cancel()
		if err != nil {
			fmt.Printf("Additional Connection %d timed out! Great!: %v", i, err)
		} else {
			fmt.Printf("Additional Connection %d aqquired! Thats bad!", i)
			connections = append(connections, conn)
		}
	}
	fmt.Printf("Acquired connections in total: %d, should be 10!", len(connections))
	for _, connection := range(connections) {
		p.Release(connection)
	}
	fmt.Print("Finished!")
}
