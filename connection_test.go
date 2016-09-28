package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

type TestConnection struct {
	mock.Mock
	*amqp.Connection
}

func (c *TestConnection) Close() error {
	args := c.Called()
	return args.Error(0)
}
