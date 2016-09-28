package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

// This file only contains a TestChannel with specific methods that are known to be used supporting mocking
type TestChannel struct {
	*amqp.Channel
	mock.Mock
}

func (t *TestChannel) Consume(queueName, consumerTag string, autoAck, exclusive, noLocal, noWait bool, arguments amqp.Table) (<-chan amqp.Delivery, error) {
	args := t.Called(queueName, consumerTag, autoAck, exclusive, noLocal, noWait, arguments)
	return args.Get(0).(chan amqp.Delivery), args.Error(1)
}

func (t *TestChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	args := t.Called(exchange, key, mandatory, immediate, msg)
	return args.Error(0)
}

func (t *TestChannel) ExchangeDeclare(exchange, exchangeType string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	arguments := t.Called(exchange, exchangeType, durable, autoDelete, internal, noWait, args)
	return arguments.Error(0)
}

func (t *TestChannel) QueueDeclare(queueName string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	arguments := t.Called(queueName, durable, autoDelete, exclusive, noWait, args)
	return arguments.Get(0).(amqp.Queue), arguments.Error(1)
}

func (t *TestChannel) QueueBind(queueName, routingKey, exchangeName string, noWait bool, args amqp.Table) error {
	arguments := t.Called(queueName, routingKey, exchangeName, noWait, args)
	return arguments.Error(0)
}

func (t *TestChannel) Cancel(consumerTag string, noWait bool) error {
	args := t.Called(consumerTag, noWait)
	return args.Error(0)
}
