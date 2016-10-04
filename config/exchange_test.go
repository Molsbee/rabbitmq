package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuilder(t *testing.T) {
	// arrange
	expected := Exchange{
		Name:       "test",
		Type:       "direct",
		Durable:    false,
		AutoDelete: false,
		Internal:   false,
		Args:       nil,
	}

	// act
	exchangeBuilder := ExchangeBuilder(expected.Name).Type(expected.Type)
	exchangeBuilder.Durable(expected.Durable).AutoDelete(expected.AutoDelete)
	exchangeBuilder.Internal(expected.Internal).NoWait(expected.NoWait)
	actual := exchangeBuilder.Args(expected.Args).Build()

	// assert
	assert.Equal(t, expected, *actual)
}

func TestBuilder_DefaultValues(t *testing.T) {
	// act
	exchange := ExchangeBuilder("test").Build()

	// assert
	assert.Equal(t, "direct", exchange.Type)
	assert.Equal(t, true, exchange.Durable)
	assert.Equal(t, false, exchange.AutoDelete)
	assert.Equal(t, false, exchange.Internal)
	assert.Equal(t, false, exchange.NoWait)
	assert.Nil(t, exchange.Args)
}
