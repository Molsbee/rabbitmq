package queue

import "github.com/streadway/amqp"

type Queue interface {
	Name() string
	Durable() bool
	AutoDelete() bool
	Exclusive() bool
	NoWait() bool
	Args() amqp.Table
}

type queue struct {
	name       string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
	args       amqp.Table
}

func (q *queue) Name() string {
	return q.name
}

func (q *queue) Durable() bool {
	return q.durable
}

func (q *queue) AutoDelete() bool {
	return q.autoDelete
}

func (q *queue) Exclusive() bool {
	return q.exclusive
}

func (q *queue) NoWait() bool {
	return q.noWait
}

func (q *queue) Args() amqp.Table {
	return q.args
}

type QueueBuilder interface {
	Durable(bool) QueueBuilder
	AutoDelete(bool) QueueBuilder
	Exclusive(bool) QueueBuilder
	NoWait(bool) QueueBuilder
	Args(amqp.Table) QueueBuilder
	Build() Queue
}

type queueBuilder struct {
	name       string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
	args       amqp.Table
}

func Builder(name string) QueueBuilder {
	return &queueBuilder{
		name:       name,
		durable:    true,
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
		args:       nil,
	}
}

func (b *queueBuilder) Durable(durable bool) QueueBuilder {
	b.durable = durable
	return b
}

func (b *queueBuilder) AutoDelete(autoDelete bool) QueueBuilder {
	b.autoDelete = autoDelete
	return b
}

func (b *queueBuilder) Exclusive(exclusive bool) QueueBuilder {
	b.exclusive = exclusive
	return b
}

func (b *queueBuilder) NoWait(noWait bool) QueueBuilder {
	b.noWait = noWait
	return b
}

func (b *queueBuilder) Args(args amqp.Table) QueueBuilder {
	b.args = args
	return b
}

func (b *queueBuilder) Build() Queue {
	return &queue{
		name:       b.name,
		durable:    b.durable,
		autoDelete: b.autoDelete,
		exclusive:  b.exclusive,
		noWait:     b.noWait,
		args:       b.args,
	}
}
