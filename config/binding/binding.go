package binding

import (
	"github.com/molsbee/rabbitmq/config/exchange"
	"github.com/molsbee/rabbitmq/config/queue"
	"github.com/streadway/amqp"
)

type Binding struct {
	queue    string
	key      string
	exchange string
	noWait   bool
	args     amqp.Table
}

func (b *Binding) QueueName() string {
	return b.queue
}

func (b *Binding) Key() string {
	return b.key
}

func (b *Binding) ExchangeName() string {
	return b.exchange
}

func (b *Binding) NoWait() bool {
	return b.noWait
}

func (b *Binding) Args() amqp.Table {
	return b.args
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
	Build() *Binding
}

type bindingBuilder struct {
	queue    queue.Queue
	exchange exchange.Exchange
	key      string
	args     amqp.Table
	noWait   bool
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

func (b *bindingBuilder) Build() *Binding {
	return &Binding{
		queue:    b.queue.Name(),
		key:      b.key,
		exchange: b.exchange.Name(),
		noWait:   b.noWait,
		args:     b.args,
	}
}
