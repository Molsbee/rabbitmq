package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	amqpQueue    = amqp.Queue{Name: "test.queue"}
	exchangeName = "minion.test"
	routingKey   = "test"
	consumerTag  = "test"
)

func TestPush_SuccessfulPublish(t *testing.T) {
	// arrange
	connection := new(TestConnection)

	channel := new(TestChannel)
	channel.On("ExchangeDeclare", exchangeName, "direct", true, false, false, false, mock.Anything).Return(nil)
	channel.On("QueueDeclare", amqpQueue.Name, true, false, false, false, mock.Anything).Return(amqpQueue, nil)
	channel.On("QueueBind", amqpQueue.Name, routingKey, exchangeName, false, mock.Anything).Return(nil)
	channel.On("Publish", exchangeName, routingKey, false, false, mock.Anything).Return(nil)

	rabbitmq := createTestRabbitMQ(connection, channel)
	queue, _ := rabbitmq.NewRabbitQueue(amqpQueue.Name, exchangeName, routingKey, nil)

	message := struct {
		Test  string
		Stuff string
	}{"this", "test"}

	// act
	err := queue.Push(message)

	// assert
	assert.Nil(t, err)
}

func TestPush_FailedPublish(t *testing.T) {
	// arrange
	connection := new(TestConnection)

	channel := new(TestChannel)
	channel.On("ExchangeDeclare", exchangeName, "direct", true, false, false, false, mock.Anything).Return(nil)
	channel.On("QueueDeclare", amqpQueue.Name, true, false, false, false, mock.Anything).Return(amqpQueue, nil)
	channel.On("QueueBind", amqpQueue.Name, routingKey, exchangeName, false, mock.Anything).Return(nil)
	channel.On("Publish", exchangeName, routingKey, false, false, mock.Anything).Return(fmt.Errorf("Failure"))

	rabbitmq := createTestRabbitMQ(connection, channel)
	queue, _ := rabbitmq.NewRabbitQueue(amqpQueue.Name, exchangeName, routingKey, nil)

	message := struct {
		Test  string
		Stuff string
	}{"this", "test"}

	// act
	err := queue.Push(message)

	// assert
	assert.NotNil(t, err)
}

func TestListen(t *testing.T) {
	// arrange
	connection := new(TestConnection)
	connection.On("Close").Return(nil)

	channel := new(TestChannel)
	deliveryChannel := make(chan amqp.Delivery, 5)
	channel.On("ExchangeDeclare", exchangeName, "direct", true, false, false, false, mock.Anything).Return(nil)
	channel.On("QueueDeclare", amqpQueue.Name, true, false, false, false, mock.Anything).Return(amqpQueue, nil)
	channel.On("QueueBind", amqpQueue.Name, routingKey, exchangeName, false, mock.Anything).Return(nil)
	channel.On("Consume", amqpQueue.Name, consumerTag, false, false, false, false, mock.Anything).Return(deliveryChannel, nil)
	channel.On("Cancel", consumerTag, true).Return(nil)

	rabbitmq := createTestRabbitMQ(connection, channel)
	defer rabbitmq.Disconnect()

	rabbitQueue, err := rabbitmq.NewRabbitQueue(amqpQueue.Name, exchangeName, routingKey, nil)
	assert.Nil(t, err)

	// act
	executed := make(chan bool)
	rabbitQueue.Listen(func(d amqp.Delivery) {
		t.Log("Listen function called with delivery item")
		assert.NotNil(t, d)
		executed <- true
	})

	deliveryChannel <- amqp.Delivery{}

	// assert
	called := <-executed
	assert.True(t, called)
}

func createTestRabbitMQ(connection *TestConnection, channel *TestChannel) *RabbitMQ {
	return &RabbitMQ{
		AmqpEndpoint: "amqp://guest:guest@localhost:5672/",
		ConsumerTag:  consumerTag,
		Connection:   connection,
		Channel:      channel,
		shutdown:     make(chan bool),
	}
}
