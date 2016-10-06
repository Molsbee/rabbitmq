package config

import "github.com/streadway/amqp"

type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type QueueBuilder interface {
	Durable(bool) QueueBuilder
	AutoDelete(bool) QueueBuilder
	Exclusive(bool) QueueBuilder
	NoWait(bool) QueueBuilder
	Args(amqp.Table) QueueBuilder
	Build() *Queue
}

type queueBuilder struct {
	name       string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
	args       amqp.Table
}

func QueueBuilder(name string) QueueBuilder {
	return &queueBuilder{
		name:       name,
		durable:    true,
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
		args:       nil,
	}
}

func (q *queueBuilder) Durable(durable bool) QueueBuilder {
	q.durable = durable
	return q
}

func (q *queueBuilder) AutoDelete(autoDelete bool) QueueBuilder {
	q.autoDelete = autoDelete
	return q
}

func (q *queueBuilder) Exclusive(exclusive bool) QueueBuilder {
	q.exclusive = exclusive
	return q
}

func (q *queueBuilder) NoWait(noWait bool) QueueBuilder {
	q.noWait = noWait
	return q
}

func (q *queueBuilder) Args(args amqp.Table) QueueBuilder {
	q.args = args
	return q
}

func (q *queueBuilder) Build() *Queue {
	return &Queue{
		Name:       q.name,
		Durable:    q.durable,
		AutoDelete: q.autoDelete,
		Exclusive:  q.exclusive,
		NoWait:     q.noWait,
		Args:       q.args,
	}
}
