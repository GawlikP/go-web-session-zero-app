package rabbitmq

import (
	"context"
	"fmt"
)

const (
	AuthLoginQueue = "auth.login.attempts"
	AuthCommandsExchange = "auth.commands"
	AuthLoginRoute  = "login.attempt"
	AuthRegisterQueue = "auth.register.attempts"
	AuthRegisterRoute  = "register.attempt"
)

func InitializeTopology(ctx context.Context, p *ConnectionPool) error {
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

	err = ch.ExchangeDeclare(
		AuthCommandsExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("%w, %w", ErrCannotDeclareExchange, err)
	}

	_, err = ch.QueueDeclare(
		AuthLoginQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%w, %w", ErrCannotDeclareQueue, err)
	}
	
	err = ch.QueueBind(
		AuthLoginQueue,
		AuthLoginRoute,
		AuthCommandsExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%w, %w", ErrCannotBindQueue, err)
	}

	_, err = ch.QueueDeclare(
		AuthRegisterQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%w, %w", ErrCannotDeclareQueue, err)
	}
	
	err = ch.QueueBind(
		AuthRegisterQueue,
		AuthRegisterRoute,
		AuthCommandsExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%w, %w", ErrCannotBindQueue, err)
	}
	return nil
}
