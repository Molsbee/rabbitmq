package rabbitmq

import (
	"github.com/stretchr/testify/mock"
	"github.com/streadway/amqp"
)

type TestConnection struct {
	mock.Mock
	*amqp.Connection
}

func (c *TestConnection) Close() error {
	args := c.Called()
	return args.Error(0)
}