package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Job struct {
	Successful bool
	Output     string
}

func TestPublish_Success(t *testing.T) {
	// arrange
	publishChannel := &TestChannel{}
	connection := &TestConnection{}
	queue := createTestQueue(connection, publishChannel)

	publishChannel.On("Publish", queue.ExchangeName(), queue.RoutingKey(), false, false, mock.AnythingOfType("amqp.Publishing")).Return(nil)

	// act
	err := queue.Publish(Job{Successful: true, Output: "Stuff happened"})

	// assert
	assert.Nil(t, err)
}

func TestPublish_Failed(t *testing.T) {
	// arrange
	publishChannel := &TestChannel{}
	connection := &TestConnection{}
	queue := createTestQueue(connection, publishChannel)

	publishChannel.On("Publish", queue.ExchangeName(), queue.RoutingKey(), false, false, mock.AnythingOfType("amqp.Publishing")).Return(fmt.Errorf("Failure"))

	// act
	err := queue.Publish(Job{Successful: false, Output: "Stuff failed"})

	// assert
	assert.NotNil(t, err)
}

func createTestQueue(connection *TestConnection, channel *TestChannel) *RabbitQueue {
	return &RabbitQueue{
		queueName:    "test",
		exchangeName: "test.exchange",
		routingKey:   "test",
		rabbitMQ: &Rabbit{
			endpoint:   "amqp://guest:guest@localhost:5672/",
			connection: connection,
			shutdown:   make(chan bool),
		},
		publishChannel: channel,
	}
}
