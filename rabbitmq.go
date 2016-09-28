package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

// RabbitMQ data
type RabbitMQ struct {
	AmqpEndpoint string
	ConsumerTag  string
	Connection   AmqpConnection

	// TODO: Figure out how to support Publish channel and Consumer Channel
	Channel  AmqpChannel
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

	go rabbitMQ.reconnectionStrategy()

	return &rabbitMQ, nil
}

func (r *RabbitMQ) establishConnection() error {
	log.Printf("dialing %s\n", r.AmqpEndpoint)
	var err error
	r.Connection, err = amqp.Dial(r.AmqpEndpoint)
	if err != nil {
		log.Printf("message: failed to esablish a connection to endpoint %s, error: %+v", r.AmqpEndpoint, err)
		return err
	}

	r.Channel, err = r.Connection.Channel()
	if err != nil {
		log.Printf("message: failed to establish a channel from the connection, error: %+v", err)
		return err
	}

	return nil
}

func (r *RabbitMQ) reconnectionStrategy() {
	for {
		log.Println("message: registering notify channel for close events in case of RabbitMQ failure")
		errorChan := make(chan *amqp.Error)
		r.Connection.NotifyClose(errorChan)

		select {
		case err := <-errorChan:
			log.Println("message: lost connection to RabbitMQ all channels have been closed", err)
			for {
				if err := r.establishConnection(); err != nil {
					log.Println("message: sleeping before attempting to re-establish connection")
					time.Sleep(10 * time.Second)
					continue
				}
				break
			}
		case <-r.shutdown:
			log.Println("message: shutting down reconnection strategy")
			return
		}
	}
}

// Disconnect will shut down all resources defined by RabbitMQ struct including RabbitQueues that have been created
// through NewRabbitQueue.  After those resources have been shut down then the channel/connection will be closed
// from RabbitMQ
func (r *RabbitMQ) Disconnect() {
	close(r.shutdown) // Shutdown goroutine for loop that handler reconnection
	r.Channel.Cancel(r.ConsumerTag, false)
	r.Connection.Close()
}

// NewRabbitQueue wraps the exchange and queue declaration and their binding. A pointer is returned to a RabbitQueue struct
// which contains all the necessary information for communicating with queue.
// Assumptions
// Exchange [type: direct, durable: true, autoDelete: false, internal: false, noWait: false, args: nil]
// Queue [durable: true, autoDelete: false, exclusive: false, noWait: false, args: (provided as parameter)]
// QueueBind [noWait: false, args: nil]
func (r *RabbitMQ) NewRabbitQueue(queueName string, exchangeName string, routingKey string, args amqp.Table) (*RabbitQueue, error) {
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
		QueueName:    queue.Name,
		ExchangeName: exchangeName,
		RoutingKey:   routingKey,
		RabbitMQ:     r,
	}, nil
}
