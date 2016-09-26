package rabbitmq

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

// RabbitMQ data
type RabbitMQ struct {
	AmqpEndpoint string
	ConsumerTag  string
	Connection AmqpConnection
	Channel    AmqpChannel
	shutdown chan bool
}

// Connect creates a new RabbitMQ connection and returns a struct which wraps the connection and channel along with
// meta information like the Endpoint and ConsumerTag provided.  The method also starts up a go routine with a for loop
// that will continually try to re-establish a connection with RabbitMQ in the event of network/server failure.
func Connect(url string, consumerTag string) (*RabbitMQ, error) {
	rabbitMQ := RabbitMQ{
		AmqpEndpoint: url,
		ConsumerTag:  consumerTag,
		shutdown:     make(chan bool),
	}

	err := rabbitMQ.establishConnection()
	if err != nil {
		return nil, err
	}

	go rabbitMQ.reconnectStrategy()

	return &rabbitMQ, nil
}

func (r *RabbitMQ) establishConnection() (err error) {
	r.Connection, err = amqp.Dial(r.AmqpEndpoint)
	if err != nil {
		return
	}

	r.Channel, err = r.Connection.Channel()
	if err != nil {
		return
	}

	return nil
}

func (r *RabbitMQ) reconnectStrategy() {
	for {
		log.Info("Registering notify channel in case of RabbitMQ failure")
		errorChan := make(chan *amqp.Error)
		r.Channel.NotifyClose(errorChan)

		select {
		case err := <-errorChan:
			if err != nil {
				log.WithField("error", err).Error("Close channel called with error")
			}
			for {
				if err := r.establishConnection(); err != nil {
					log.WithField("error", err).Error("Failed to reconnect")
					time.Sleep(10 * time.Second)
					continue
				}
				break
			}
		case <-r.shutdown:
			log.Info("Shutting down connect for loop")
			return
		}
	}
}

// Disconnect will shut down all resources defined by RabbitMQ struct including RabbitQueues that have been created
// through NewRabbitQueue.  After those resources have been shut down then the channel/connection will be closed
// from RabbitMQ
func (r *RabbitMQ) Disconnect() {
	close(r.shutdown) // Shutdown goroutine for loop that handler reconnection
	r.Channel.Cancel(r.ConsumerTag, true)
	r.Connection.Close()
}

// RabbitQueue is a struct that contains all the necessary information to publish to an exchange with a routing key
// and have it end up in the queue defined by NewRabbitQueue.  All construction of struct should go through function
// NewRabbitQueue as it has built in behavior for defining the queue, exchange and binding them together with the
// routing key
type RabbitQueue struct {
	QueueName string
	Exchange  string
	Key       string
	RabbitMQ  *RabbitMQ
}

// NewRabbitQueue wraps the exchange and queue declaration and their binding. A point is returned to a RabbitQueue struct
// which contains all the necessary information for communicating with queue.
func (r *RabbitMQ) NewRabbitQueue(queueName string, exchangeName string, routingKey string, args amqp.Table) (*RabbitQueue, error) {
	// TODO: see this article for ensuring published messages are committed to the broker (confirm select)
	// https://agocs.org/blog/2014/08/19/rabbitmq-best-practices-in-go/
	if err := r.Channel.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil); err != nil {
		return nil, err
	}

	queue, err := r.Channel.QueueDeclare(queueName, true, false, false, false, args)
	if err != nil {
		return nil, err
	}

	if err := r.Channel.QueueBind(queue.Name, routingKey, exchangeName, false, nil); err != nil {
		return nil, err
	}

	return &RabbitQueue{
		QueueName: queue.Name,
		Exchange:  exchangeName,
		Key:       routingKey,
		RabbitMQ:  r,
	}, nil
}

// Push converts a struct to a json byte array and construct an amqp.Publishing struct with json as body and
// content type set to application/json.  After message has been construct it is published to the exchange with routing
// key to make it to the proper queue.
func (r *RabbitQueue) Push(msg interface{}) error {
	jsonMessage, _ := json.Marshal(msg)

	message := amqp.Publishing{
		ContentType: "application/json;charset=UTF-8",
		Body:        jsonMessage,
	}
	return r.RabbitMQ.Channel.Publish(r.Exchange, r.Key, false, false, message)
}

// Listen is a generic mechanism for re-establishing RabbitMQ delivery channel in the event it's shutdown
// and will call the function provided with the amqp.Delivery item.  When shutdown is called on the RabbitMQ struct
// it will also close the for loop that's controlling the delivery channel, providing a mechanism to cleanly shutdown
// RabbitMQ
func (r *RabbitQueue) Listen(fn func(d amqp.Delivery)) {
	go func() {
		log.Info("Setting up listener to process messages and delegate to function")
		for {
			START:
			log.Info("Establishing delivery channel")
			deliveries, err := r.RabbitMQ.Channel.Consume(r.QueueName, r.RabbitMQ.ConsumerTag, false, false, false, false, nil)
			if err != nil {
				log.WithField("error", err).Error("Failed to establish a delivery channel from RabbitMQ")
				time.Sleep(10 * time.Second)
				continue
			}

			for {
				select {
				case d, ok := <-deliveries:
					if !ok {
						goto START
					}
					log.Info("Executing function with delivery item")
					fn(d)
					log.Info("Finished executing function with delivery item")
				case <-r.RabbitMQ.shutdown:
					log.Info("Shutting down listener")
					return
				}
			}
		}
	}()
}