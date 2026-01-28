package main

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	fmt.Println("✅ Connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create channel: %v", err)
	}
	defer ch.Close()
	fmt.Println("✅ Channel opened")

	// Enable publisher confirms
	if err := ch.Confirm(false); err != nil {
		log.Fatalf("Failed to put channel in confirm mode: %v", err)
	}
	fmt.Println("✅ Publisher confirms enabled")

	// Get confirmation channel
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	exchange := "test.exchange"
	routingKey := "test.key"

	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Message #%d sent at %s", i, time.Now().Format("15:04:05"))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		err = ch.PublishWithContext(
			ctx,
			exchange,
			routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte(message),
				Timestamp:    time.Now(),
				DeliveryMode: amqp.Persistent,
			},
		)
		cancel()

		if err != nil {
			log.Printf("❌ Failed to publish message %d: %v", i, err)
			continue
		}

		// Wait for confirmation from RabbitMQ
		confirmed := <-confirms
		if confirmed.Ack {
			fmt.Printf("✅ Confirmed: %s\n", message)
		} else {
			fmt.Printf("❌ Message not confirmed: %s\n", message)
		}

		time.Sleep(1 * time.Second)
	}

	fmt.Println("✅ All messages published and confirmed!")
}
