package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Rabbit struct {
	endpoint   string
	connection AmqpConnection
	shutdown   chan bool
}

// Connect creates a new connection and returns a struct which wraps the underlying amqp connection.
// The method also starts up a go routine which is responsible for attempting to re-establish the connection in the
// event of a network/server failure.
func Connect(url string) (*Rabbit, error) {
	rabbitMQ := Rabbit{
		endpoint: url,
		shutdown: make(chan bool),
	}

	err := rabbitMQ.establishConnection()
	if err != nil {
		return nil, err
	}

	rabbitMQ.setupReconnectionStrategy()

	return &rabbitMQ, nil
}

// establishConnection will attempt to dial the connection provided if successful update Rabbit struct with connection
// otherwise return error
func (r *Rabbit) establishConnection() (err error) {
	log.Printf("dialing rabbitmq - %s\n", r.endpoint)
	r.connection, err = amqp.Dial(r.endpoint)
	if err != nil {
		log.Printf("failed to esablish a connection to endpoint %s - error: %+v", r.endpoint, err)
	}

	return
}

// setupReconnectionStrategy will start a go routine which will infinitely try to establish a connection with rabbitmq
func (r *Rabbit) setupReconnectionStrategy() {
	go func() {
		for {
			log.Println("registering notify channel with rabbitmq connection in case of failure")
			errorChan := make(chan *amqp.Error)
			r.connection.NotifyClose(errorChan)

			select {
			case err := <-errorChan:
				log.Println("lost connection to rabbitmq all channels have been closed", err)
				for {
					log.Println("attmeping to re-establish rabbitmq connection")
					if err := r.establishConnection(); err != nil {
						log.Println("failed to re-establish connection sleeping before re-attempting")
						time.Sleep(10 * time.Second)
						continue
					}
					break
				}
			case <-r.shutdown:
				log.Println("shutting down reconnection strategy for rabbitmq connection")
				return
			}
		}
	}()
}

// channel provides a safe mechanism for getting a channel from the amqp connection if you directly get a channel
// from the connection a panic will occur if the connection is dead.
func (r *Rabbit) channel() (channel AmqpChannel, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
			err = fmt.Errorf("failed to establish channel")
		}
	}()
	channel, err = r.connection.Channel()
	return
}

// Disconnect will shut down all resources defined by RabbitMQ struct including RabbitQueues that have been created
// through NewRabbitQueue.  After those resources have been shut down then the channel/connection will be closed
// from RabbitMQ
func (r *Rabbit) Disconnect() {
	close(r.shutdown) // Shutdown goroutine for loop that handler reconnection
	//r.Channel.Cancel(r.ConsumerTag, false)
	r.connection.Close()
}

// NewRabbitQueue wraps the exchange and queue declaration and their binding. A pointer is returned to a RabbitQueue struct
// which contains all the necessary information for communicating with queue.
// Assumptions
// Exchange [type: direct, durable: true, autoDelete: false, internal: false, noWait: false, args: nil]
// Queue [durable: true, autoDelete: false, exclusive: false, noWait: false, args: (provided as parameter)]
// QueueBind [noWait: false, args: nil]
func (r *Rabbit) NewRabbitQueue(queueName string, exchangeName string, routingKey string, args amqp.Table) (*RabbitQueue, error) {
	return newRabbitQueue(queueName, exchangeName, routingKey, args, r)
}
