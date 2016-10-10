package exchange

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExchangeBuilder(t *testing.T) {
	// arrange
	expected := exchange{
		name:         "test",
		exchangeType: "direct",
		durable:      false,
		autoDelete:   false,
		internal:     false,
		args:         nil,
	}

	// act
	exchangeBuilder := Builder(expected.Name()).Type(expected.Type())
	exchangeBuilder.Durable(expected.Durable()).AutoDelete(expected.AutoDelete())
	exchangeBuilder.Internal(expected.Internal()).NoWait(expected.NoWait())
	actual := exchangeBuilder.Args(expected.Args()).Build()

	// assert
	assert.Equal(t, &expected, actual)
}

func TestExchangeBuilder_DefaultValues(t *testing.T) {
	// act
	exchange := Builder("test").Build()

	// assert
	assert.Equal(t, "direct", exchange.Type())
	assert.Equal(t, true, exchange.Durable())
	assert.Equal(t, false, exchange.AutoDelete())
	assert.Equal(t, false, exchange.Internal())
	assert.Equal(t, false, exchange.NoWait())
	assert.Nil(t, exchange.Args())
}
