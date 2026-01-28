package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewConsumer(url, queue string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	// Configure QoS
	if err := ch.Qos(10, 0, false); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		queue:   queue,
	}, nil
}

func (c *Consumer) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}

func (c *Consumer) Consume(handler func([]byte) error) error {
	msgs, err := c.channel.Consume(
		c.queue,
		"",    // consumer tag
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	fmt.Printf("✅ Consumer started on queue: %s\n", c.queue)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				return nil // Channel closed
			}
			c.handleMessage(msg, handler)

		case <-sigChan:
			fmt.Println("\n🛑 Shutting down consumer...")
			return nil
		}
	}
}

func (c *Consumer) handleMessage(msg amqp.Delivery, handler func([]byte) error) {
	fmt.Printf("📥 Received: %s\n", string(msg.Body))

	// Call the handler
	err := handler(msg.Body)

	if err != nil {
		fmt.Printf("❌ Handler error: %v\n", err)
		c.handleError(msg, err)
		return
	}

	// Success - acknowledge
	if err := msg.Ack(true); err != nil {
		log.Printf("Failed to ack message: %v", err)
	} else {
		fmt.Println("   ✅ Acknowledged\n")
	}
}

func (c *Consumer) handleError(msg amqp.Delivery, err error) {
	// Get retry information from headers
	retryCount := 0
	if msg.Headers != nil {
		if count, ok := msg.Headers["x-retry-count"].(int32); ok {
			retryCount = int(count)
		}
	}

	maxRetries := 3

	if retryCount < maxRetries {
		// Requeue with incremented retry count
		fmt.Printf("   🔄 Requeuing (retry %d/%d)\n\n", retryCount+1, maxRetries)
		
		// For proper retry with delay, you'd republish to a delay exchange
		// For now, just requeue (immediate retry)
		msg.Nack(false, true)
	} else {
		// Max retries exceeded
		fmt.Printf("   ⚠️  Max retries exceeded, rejecting\n\n")
		msg.Reject(false) // Goes to dead letter queue if configured
	}
}

func main() {
	consumer, err := NewConsumer(
		"amqp://admin:admin@localhost:5672/",
		"test.queue",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	// Define message handler
	handler := func(body []byte) error {
		// Your business logic here
		fmt.Printf("   Processing: %s\n", string(body))
		
		// Simulate work
		time.Sleep(100 * time.Millisecond)
		
		// Simulate occasional errors
		// if rand.Float32() < 0.2 {
		//     return fmt.Errorf("processing failed")
		// }
		
		return nil
	}

	// Start consuming
	if err := consumer.Consume(handler); err != nil {
		log.Fatal(err)
	}
}
