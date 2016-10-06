package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestQueueBuilder(t *testing.T) {
	// arrange
	expected := queue{
		name:       "test",
		durable:    false,
		autoDelete: true,
		exclusive:  true,
		noWait:     true,
		args:       nil,
	}

	// act
	queuBuilder := Builder(expected.Name()).Durable(expected.Durable())
	queuBuilder.AutoDelete(expected.AutoDelete()).Exclusive(expected.Exclusive())
	actual := queuBuilder.NoWait(expected.NoWait()).Args(expected.Args()).Build()

	// assert
	assert.Equal(t, &expected, actual)
}

func TestQueueBuilder_DefaultValues(t *testing.T) {
	// arrange
	name := "test"

	// act
	queue := Builder(name).Build()

	// assert
	assert.Equal(t, name, queue.Name())
	assert.True(t, queue.Durable())
	assert.False(t, queue.AutoDelete())
	assert.False(t, queue.Exclusive())
	assert.False(t, queue.NoWait())
	assert.Nil(t, queue.Args())
}
