package rabbitmq

import "errors"

var (
	ErrConnectionFailed = errors.New("Cannot obtain a connection from defined pool")
	ErrChannelFailed = errors.New("Cannot obtain a channel from defined connection")
	ErrCannotDeclareExchange = errors.New("Failed to declare the exchange")
	ErrCannotDeclareQueue = errors.New("Failed to declare the queue")
	ErrCannotBindQueue = errors.New("Cannot Bind Queue")
	ErrCannotEnableConfirmation = errors.New("Cannot require confirmation from the channel")
	ErrFailedToPublish = errors.New("Failed To publish the message on the exchange")
	ErrFailedAck= errors.New("Failed To Acknowldege on the exchange")
)
