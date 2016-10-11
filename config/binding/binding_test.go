package binding

import (
	"testing"
	"github.com/molsbee/rabbitmq/config/queue"
	"github.com/molsbee/rabbitmq/config/exchange"
	"github.com/stretchr/testify/assert"
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
