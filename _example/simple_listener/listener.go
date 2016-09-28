package main

import (
	"flag"
	"log"
	"github.com/molsbee/rabbitmq"
	"github.com/streadway/amqp"
)

var (
	rabbit = flag.String("amqp", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	consumerTag = flag.String("consumer-tag", "simple-listener", "AMQP consumer tag")
)

func main() {
	flag.Parse()

	rabbitMQ, err := rabbitmq.Connect(*rabbit, *consumerTag)
	if err != nil {
		log.Println("Failled to esablish a connection or channel with RabbitMQ", err)
	}

	queue, err := rabbitMQ.NewRabbitQueue("test", "test.exchange", "test", nil)
	if err != nil {
		log.Println("Failed to declare queue, exchange and bind them", err)
	}

	queue.ListenerFunc(func(d amqp.Delivery) {
		log.Println(d.Body)
		d.Ack(false)
	})

	block := make(chan bool)
	<- block
}