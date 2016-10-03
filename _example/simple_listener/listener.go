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

	rabbitMQ, err := rabbitmq.Connect(*rabbit)
	if err != nil {
		log.Fatalf("Failled to esablish a connection or channel with RabbitMQ - error %+v", err)
	}

	queue, err := rabbitMQ.NewRabbitQueue("test", "test.exchange", "test", nil)
	if err != nil {
		log.Fatal("Failed to declare queue, exchange and bind them - error %+v", err)
	}

	queue.ListenerFunc(*consumerTag, func(d amqp.Delivery) {
		log.Println(d.Body)
		d.Ack(false)
	})

	block := make(chan bool)
	<- block
}