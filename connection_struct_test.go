package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

type TestConnection struct {
	*amqp.Connection
	mock.Mock
}

func (c *TestConnection) Channel() (*amqp.Channel, error) {
	args := c.Called()
	return args.Get(0).(*amqp.Channel), args.Error(1)
}

func (c *TestConnection) Close() error {
	args := c.Called()
	return args.Error(0)
}
