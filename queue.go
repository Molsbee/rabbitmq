package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Queue interface {
	QueueName() string
	ExchangeName() string
	RoutingKey() string
	Publish(msg struct{}) error
	ListenerFunc(consumerTag string, deliveryFunc func(d amqp.Delivery))
}

// RabbitQueue is a struct that contains all the necessary information to publish to an exchange with a routing key
// and have it end up in the queue defined by NewRabbitQueue.  All construction of struct should go through function
// NewRabbitQueue as it has built in behavior for defining the queue, exchange and binding them together with the
// routing key
type RabbitQueue struct {
	queueName       string
	exchangeName    string
	routingKey      string

	rabbitMQ        *Rabbit // TODO: Handle in a more simplified or restrictive manor

	publishChannel  AmqpChannel
}

// newRabbitQueue constructs a RabbitQueue which contains all the necessary information for accessing a queue.
// This method will create a channel that is used by the Publish method and declares the
// Queue    - [durable: true, autoDelete: false, exclusive: false, noWait: false, args: amqp.Table]
// Exchange - [type: direct, durable: true, autoDelete: false, internal: false, noWait: false, args: nil]
// Binding  - [noWait: false, args nil]
func newRabbitQueue(queue, exchange, key string, args amqp.Table, rabbit *Rabbit) (*RabbitQueue, error) {
	channel, err := rabbit.channel()
	if err != nil {
		log.Println("failed to establish (publish) channel from rabbitmq connection")
		return nil, err
	}


	if err = channel.ExchangeDeclare(exchange, "direct", true, false, false, false, nil); err != nil {
		log.Printf("failed to declared the exchange %s\n", exchange)
		return nil, err
	}

	_, err = channel.QueueDeclare(queue, true, false, false, false, args)
	if err != nil {
		log.Printf("failed to declare persistent queue %s", queue)
		return nil, err
	}

	if err := channel.QueueBind(queue, key, exchange, false, nil); err != nil {
		log.Printf("failed to bind exchange %s to queue %s with routing key %s", exchange, queue, key)
		return nil, err
	}

	rabbitQueue := &RabbitQueue{
		queueName: queue,
		exchangeName: exchange,
		routingKey: key,
		rabbitMQ: rabbit,
		publishChannel: channel,
	}

	go rabbitQueue.reconnectChannel()

	return rabbitQueue, nil
}

func (r *RabbitQueue) reconnectChannel() {
	for {
		log.Println("message: registering notify channel for close events")
		errorChan := make(chan *amqp.Error)
		r.publishChannel.NotifyClose(errorChan)

		select {
		case <- errorChan:
			log.Println("message: lost connection to RabbitMQ")
			for {
				channel, err := r.rabbitMQ.channel()
				if err != nil {
					log.Println("message: sleeping before attempting to re-establish publish channel")
					time.Sleep(10 * time.Second)
					continue
				}
				r.publishChannel = channel
				break
			}
		case <-r.rabbitMQ.shutdown:
			log.Println("message: shutting down reconnection strategy")
		}
	}
}

// QueueName convenience method for retriving the name of the queue declared
func (r *RabbitQueue) QueueName() string {
	return r.queueName
}

// ExchangeName convenience method for retrieving the name of the exchange declared
func (r *RabbitQueue) ExchangeName() string {
	return r.exchangeName
}

// RoutingKey convenience method for retrieving the key binding the queue and exchange
func (r *RabbitQueue) RoutingKey() string {
	return r.routingKey
}

// Publish will convert a message (struct) to a json byte array and construct a amqp.Publishing struct with content
// type set to "application/json;charset=UTF-8", timestamp set for when executed, and body set to the byte array message.
// The message will then be published to RabbitMQ with mandatory and immediate set to false.
//
// Error
// JSON: Unable to marshal message
// RabbitMQ: Network connection issue or server resource constraint
func (r *RabbitQueue) Publish(msg interface{}) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		ContentType: "application/json;charset=UTF-8",
		Timestamp:   time.Now(),
		Body:        jsonMsg,
	}

	return r.publishChannel.Publish(r.exchangeName, r.routingKey, false, false, message)
}

// ListenerFunc is a generic mechanism for re-establishing RabbitMQ delivery channel in the event it's shutdown
// and calling a provided function with a delivery item to be acted upon.  When Shutdown is called on RabbitMQ it will
// also cleanly close the loop managing the listener, providing a clean method for Shutting down all goroutines
// TODO: Refactor this to call the configured version
func (r *RabbitQueue) ListenerFunc(consumerTag string, fn func(d amqp.Delivery)) {
	go func() {
		log.Println("setting up [listener] to process messages and delegate to function")
		for {
		START:
			log.Println("establishing [listener] delivery channel for processing queue messages")
			channel, err := r.rabbitMQ.channel()
			if err != nil {
				log.Printf("failed to create [listener] channel")
				time.Sleep(10 * time.Second)
				continue
			}

			deliveries, err := channel.Consume(r.queueName, consumerTag, false, false, false, false, nil)
			if err != nil {
				log.Printf("failed to establish a delivery channel with rabbitmq, error: %v\n", err)
				time.Sleep(10 * time.Second)
				continue
			}

			for {
				select {
				case d, ok := <-deliveries:
					if !ok {
						goto START
					} // Checks the health of channel if closed goto START label
					log.Println("executing function with delivery item")
					fn(d)
					log.Println("finished executing function with delivery item")
				case <-r.rabbitMQ.shutdown:
					log.Println("shutting down listener")
					channel.Cancel(consumerTag, false)
					return
				}
			}
		}
	}()
}

// ListenerConfiguration contains fields used within the Consume call to RabbitMQ.
// AutoAck
// 	When true server will acknowledge deliveries to consumer prior to writing to the network
//  	The consumer should not call Delivery.Ack
//
// Exclusive
//	When true the server will ensure that there is only one consumer for the queue
//	When false the server will distribute deliveries across multiple consumers.
//
// NoLocal
//	When true the server will not deliver publishing sent from the same connection to the consumer
//	It's advisable to use separate connections for Channel.Publish and Channel.Consumer so not to have TCP
// 	pushback on publishing affect the ability to consume messages
//
// NoWait
//	When true do not wait on the server to confirm the request and immediately begin deliveries.
//
// Args
type ListenerConfiguration struct {
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

func (r *RabbitQueue) ConfiguredListenerFunc(consumerTag string, config ListenerConfiguration, fn func(d amqp.Delivery)) {
	go func() {
		log.Println("message: setting up listener to process messages and delegate to function")
		for {
			START:
			log.Println("message: establishing delivery channel")
			channel, err := r.rabbitMQ.channel()
			if err != nil {
				log.Printf("message: failed to create channel")
				time.Sleep(10 * time.Second)
				continue
			}

			deliveries, err := channel.Consume(r.queueName, consumerTag, false, false, false, false, nil)
			if err != nil {
				log.Printf("message: failed to establish a delivery channel with RabbitMQ, error: %v\n", err)
				time.Sleep(10 * time.Second)
				continue
			}

			for {
				select {
				case d, ok := <-deliveries:
					if !ok {
						goto START
					} // Checks the health of channel if closed goto START label
					log.Println("message: executing function with delivery item")
					fn(d)
					log.Println("message: finished executing function with delivery item")
				case <-r.rabbitMQ.shutdown:
					log.Println("message: shutting down listener")
					return
				}
			}
		}
	}()
}
