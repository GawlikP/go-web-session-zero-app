package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Publish(ctx context.Context, p *ConnectionPool, exchange string, routingKey string, body []byte) error {
	conn, err := p.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("%w, %w", ErrConnectionFailed, err)
	}
	defer p.Release(conn)
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("%w, %w", ErrChannelFailed, err)
	}
	defer ch.Close()

	if err := ch.Confirm(false); err != nil {
		return fmt.Errorf("%w, %w", ErrCannotEnableConfirmation, err)
	}

	// Get confirmation channel
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	err = ch.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return fmt.Errorf("%w, %w", ErrFailedToPublish, err)
	}

	confirmed := <-confirms
	if !confirmed.Ack {
		return ErrFailedAck
	}
	return nil
}
