package rabbitmq

import (
	"github.com/streadway/amqp"
)

// AmqpChannel - Interface definition of streadway amqp channel this is used to allow for swapping out
// implementations during testing
type AmqpChannel interface {
	Ack(tag uint64, multiple bool) error
	Cancel(consumer string, noWait bool) error
	Close() error
	Confirm(noWait bool) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	ExchangeDeclarePassive(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	ExchangeDelete(name string, ifUnused, noWait bool) error
	ExchangeUnbind(destination, key, source string, noWait bool, args amqp.Table) error
	Flow(active bool) error
	Get(queue string, autoAck bool) (msg amqp.Delivery, ok bool, err error)
	Nack(tag uint64, multiple bool, requeue bool) error
	NotifyCancel(c chan string) chan string
	NotifyClose(c chan *amqp.Error) chan *amqp.Error
	NotifyConfirm(ack, nack chan uint64) (chan uint64, chan uint64)
	NotifyFlow(c chan bool) chan bool
	NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation
	NotifyReturn(c chan amqp.Return) chan amqp.Return
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Qos(prefetchCount, prefetchSize int, global bool) error
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueDeclarePassive(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error)
	QueueInspect(name string) (amqp.Queue, error)
	QueuePurge(name string, noWait bool) (int, error)
	QueueUnbind(name, key, exchange string, args amqp.Table) error
	Recover(requeue bool) error
	Reject(tag uint64, requeue bool) error
	Tx() error
	TxCommit() error
	TxRollback() error
}