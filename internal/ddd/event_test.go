package ddd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvent(t *testing.T) {
	event := NewEvent("event-1", "event-1-payload")
	assert.Equal(t, "event-1", event.EventName())
	assert.NotNil(t, event.ID())
	assert.Equal(t, "event-1-payload", event.Payload())
}
