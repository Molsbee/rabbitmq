package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"time"
)

// RabbitQueue is a struct that contains all the necessary information to publish to an exchange with a routing key
// and have it end up in the queue defined by NewRabbitQueue.  All construction of struct should go through function
// NewRabbitQueue as it has built in behavior for defining the queue, exchange and binding them together with the
// routing key
type RabbitQueue struct {
	QueueName    string
	ExchangeName string
	RoutingKey   string

	RabbitMQ *RabbitMQ // TODO: Handle in a more simplified or restrictive manor
}

// Publish will convert a message (interface) to a json byte array and construct a amqp.Publishing struct with content
// type set to "application/json;charset=UTF-8", timestamp set for when executed, and body set to the byte array message.
// The message will then be published to RabbitMQ with mandatory and immediate set to false.
//
// Error
// JSON: Unable to marshal message
// RabbitMQ: Network connection issue or server resource constraint
func (r *RabbitQueue) Publish(msg struct{}) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		ContentType: "application/json;charset=UTF-8",
		Timestamp:   time.Now(),
		Body:        jsonMsg,
	}

	return r.RabbitMQ.Channel.Publish(r.ExchangeName, r.RoutingKey, false, false, message)
}

// ListenerFunc is a generic mechanism for re-establishing RabbitMQ delivery channel in the event it's shutdown
// and calling a provided function with a delivery item to be acted upon.  When Shutdown is called on RabbitMQ it will
// also cleanly close the loop managing the listener, providing a clean method for Shutting down all goroutines
// TODO: Refactor this to call the configured version
func (r *RabbitQueue) ListenerFunc(fn func(d amqp.Delivery)) {
	go func() {
		log.Println("Setting up listener to process messages and delegate to function")
		for {
		START:
			log.Println("Establishing delivery channel")
			deliveries, err := r.RabbitMQ.Channel.Consume(r.QueueName, r.RabbitMQ.ConsumerTag, false, false, false, false, nil)
			if err != nil {
				log.Printf("Message: Failed to establish a delivery channel with RabbitMQ Error: %v\n", err)
				time.Sleep(10 * time.Second)
				continue
			}

			for {
				select {
				case d, ok := <-deliveries:
					if !ok {
						goto START
					} // Checks the health of channel if closed goto START label
					log.Println("Executing function with delivery item")
					fn(d)
					log.Println("Finished executing function with delivery item")
				case <-r.RabbitMQ.shutdown:
					log.Println("Shutting down listener")
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

func (r *RabbitQueue) ConfiguredListenerFunc(config ListenerConfiguration, fn func(d amqp.Delivery)) {
	go func() {
		log.Println("Setting up listener to process messages and delegate to function")
		for {
		START:
			log.Println("Establishing delivery channel")
			deliveries, err := r.RabbitMQ.Channel.Consume(r.QueueName, r.RabbitMQ.ConsumerTag, false, false, false, false, nil)
			if err != nil {
				log.Printf("Message: Failed to establish a delivery channel with RabbitMQ Error: %v\n", err)
				time.Sleep(10 * time.Second)
				continue
			}

			for {
				select {
				case d, ok := <-deliveries:
					if !ok {
						goto START
					} // Checks the health of channel if closed goto START label
					log.Println("Executing function with delivery item")
					fn(d)
					log.Println("Finished executing function with delivery item")
				case <-r.RabbitMQ.shutdown:
					log.Println("Shutting down listener")
					return
				}
			}
		}
	}()
}
