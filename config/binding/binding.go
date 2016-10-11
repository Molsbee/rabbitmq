package binding

import (
	"github.com/molsbee/rabbitmq/config/exchange"
	"github.com/molsbee/rabbitmq/config/queue"
	"github.com/streadway/amqp"
)

type Binding interface {
	QueueName() string
	Key() string
	ExchangeName() string
	NoWait() bool
	Args() amqp.Table
}

type binding struct {
	queue    string
	key      string
	exchange string
	noWait   bool
	args     amqp.Table
}

func (b *binding) QueueName() string {
	return b.queue
}

func (b *binding) Key() string {
	return b.key
}

func (b *binding) ExchangeName() string {
	return b.exchange
}

func (b *binding) NoWait() bool {
	return b.noWait
}

func (b *binding) Args() amqp.Table {
	return b.args
}

type bindingBuilder struct {
	queue    queue.Queue
	exchange exchange.Exchange
	key      string
	args     amqp.Table
	noWait   bool
}

type QueueBinding interface {
	To(exchange exchange.Exchange) ExchangeBinding
}

type ExchangeBinding interface {
	With(key string) ParameterBinding
}

type ParameterBinding interface {
	Args(args amqp.Table) ParameterBinding
	NoWait(noWait bool) ParameterBinding
	Build() Binding
}

func Bind(queue queue.Queue) QueueBinding {
	return &bindingBuilder{
		queue:  queue,
		noWait: false,
		args:   nil,
	}
}

func (b *bindingBuilder) To(exchange exchange.Exchange) ExchangeBinding {
	b.exchange = exchange
	return b
}

func (b *bindingBuilder) With(key string) ParameterBinding {
	b.key = key
	return b
}

func (b *bindingBuilder) Args(args amqp.Table) ParameterBinding {
	b.args = args
	return b
}

func (b *bindingBuilder) NoWait(noWait bool) ParameterBinding {
	b.noWait = noWait
	return b
}

func (b *bindingBuilder) Build() Binding {
	return &binding{
		queue:    b.queue.Name(),
		key:      b.key,
		exchange: b.exchange.Name(),
		noWait:   b.noWait,
		args:     b.args,
	}
}
