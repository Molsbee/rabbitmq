package config

import "github.com/streadway/amqp"

type Exchange struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

type ExchangeBuilder interface {
	Type(string) ExchangeBuilder
	Durable(bool) ExchangeBuilder
	AutoDelete(bool) ExchangeBuilder
	Internal(bool) ExchangeBuilder
	NoWait(bool) ExchangeBuilder
	Args(amqp.Table) ExchangeBuilder
	Build() *Exchange
}

type exchangeBuilder struct {
	name         string
	exchangeType string
	durable      bool
	autoDelete   bool
	internal     bool
	noWait       bool
	args         amqp.Table
}

// ExchangeBuilder will initially create an ExchangeBuilder that has the name provided along with default values preset.
// Each method call will override the default settings and Build will return a new initialized Exchange.
// Defaults
// [Type: direct, Durable: true, AutoDelete: false, Internal: false, NoWait: false, Args: nil]
func ExchangeBuilder(name string) ExchangeBuilder {
	return &exchangeBuilder{
		name:         name,
		exchangeType: "direct",
		durable:      true,
		autoDelete:   false,
		internal:     false,
		noWait:       false,
		args:         nil,
	}
}

// Type allows you to define what type of exchange you are working with ["direct", "fanout", "topic", "headers"]
func (e *exchangeBuilder) Type(exchangeType string) ExchangeBuilder {
	e.exchangeType = exchangeType
	return e
}

// Durable will persist your exchange declarations and survive rabbitmq restarts.
func (e *exchangeBuilder) Durable(durable bool) ExchangeBuilder {
	e.durable = durable
	return e
}

// AutoDelete will act in opposition of Durable if set it will delete the exchange when their are no bindings present
func (e *exchangeBuilder) AutoDelete(autoDelete bool) ExchangeBuilder {
	e.autoDelete = autoDelete
	return e
}

// Internal exchanges won't accept publishings and aren't exposed to users of the broker
func (e *exchangeBuilder) Internal(internal bool) ExchangeBuilder {
	e.internal = internal
	return e
}

// NoWait will declared the exchange without waiting on server response
func (e *exchangeBuilder) NoWait(noWait bool) ExchangeBuilder {
	e.noWait = noWait
	return e
}

// Args are optional collection of arguments that are used by specific exchange types.
func (e *exchangeBuilder) Args(args amqp.Table) ExchangeBuilder {
	e.args = args
	return e
}

// Build constructs the data from the ExchangeBuilder into a new Exchange that can be used to configure environment
func (e *exchangeBuilder) Build() *Exchange {
	return &Exchange{
		Name:       e.name,
		Type:       e.exchangeType,
		Durable:    e.durable,
		AutoDelete: e.autoDelete,
		Internal:   e.internal,
		NoWait:     e.noWait,
		Args:       e.args,
	}
}
