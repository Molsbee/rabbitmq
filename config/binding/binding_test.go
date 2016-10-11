package binding

import (
	"github.com/molsbee/rabbitmq/config/exchange"
	"github.com/molsbee/rabbitmq/config/queue"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBindingBuilder(t *testing.T) {
	// arrange
	q := queue.Builder("test").Build()
	x := exchange.Builder("test.exchange").Build()

	// act
	binding := Bind(q).To(x).With("test").Build()

	// assert
	assert.Equal(t, q.Name(), binding.QueueName())
	assert.Equal(t, x.Name(), binding.ExchangeName())
	assert.Equal(t, "test", binding.Key())
	assert.False(t, binding.NoWait())
	assert.Nil(t, binding.Args())
}
