package ddd

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/xcheng85/session-monitor-k8s/internal/test"
)

func TestEventDispatcher(t *testing.T) {
	ctx := context.TODO()

	eventDispatcher := NewEventDispatcher[IEvent]()
	initialEventHanders := test.GetUnexportedField(reflect.ValueOf(eventDispatcher).Elem().FieldByName("eventHandlers")).([]filterableEventHandlers[IEvent])
	assert.Equal(t, 0, len(initialEventHanders))

	mockEventHandler := &MockIEventHandler[IEvent]{}
	mockEventHandler.On("HandleEvent", mock.Anything, mock.Anything).Return(nil)

	eventDispatcher.Subscribe(mockEventHandler, "event-1", "event-2")

	newEventHandlers := test.GetUnexportedField(reflect.ValueOf(eventDispatcher).Elem().FieldByName("eventHandlers")).([]filterableEventHandlers[IEvent])
	assert.Equal(t, 1, len(newEventHandlers))

	event1, event2, eventBogus := NewEvent("event-1", nil), NewEvent("event-2", nil), NewEvent("bogus-event", nil)

	eventDispatcher.Publish(ctx, event1, event2, eventBogus)
	mockEventHandler.AssertNumberOfCalls(t, "HandleEvent", 2)
	mockEventHandler.AssertCalled(t, "HandleEvent", ctx, event1)
	mockEventHandler.AssertCalled(t, "HandleEvent", ctx, event2)
	mockEventHandler.AssertNotCalled(t, "HandleEvent", ctx, eventBogus)
}
