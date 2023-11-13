package ddd

import (
	"context"
	"sync"
)

//go:generate mockery --name IEventHandler
type IEventHandler[T IEvent] interface {
	HandleEvent(ctx context.Context, event T) error
}

//go:generate mockery --name IEventPublisher
type IEventPublisher[T IEvent] interface {
	Publish(ctx context.Context, events ...T) error
}

//go:generate mockery --name IEventSubscriber
type IEventSubscriber[T IEvent] interface {
	Subscribe(handler IEventHandler[T], events ...string)
}

// dispatcher needs to subscribe and publish both.
// he can control the event flow with business logic
//
//go:generate mockery --name IEventDispatcher
type IEventDispatcher[T IEvent] interface {
	IEventPublisher[T]
	IEventSubscriber[T]
}

// decorator pattern
type filterableEventHandlers[T IEvent] struct {
	handler IEventHandler[T]
	filters map[string]struct{}
}

// container of all the event handler
type EventDispatcher[T IEvent] struct {
	eventHandlers []filterableEventHandlers[T]
	mutex         sync.Mutex
}

// composite new interface from two existing interface
var _ IEventDispatcher[IEvent] = (*EventDispatcher[IEvent])(nil)

func NewEventDispatcher[T IEvent]() IEventDispatcher[T] {
	return &EventDispatcher[T]{
		eventHandlers: []filterableEventHandlers[T]{},
	}
}

func (d *EventDispatcher[T]) Subscribe(handler IEventHandler[T], events ...string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	var filters map[string]struct{}
	if len(events) > 0 {
		filters = make(map[string]struct{})
		for _, event := range events {
			filters[event] = struct{}{}
		}
	}

	d.eventHandlers = append(d.eventHandlers, filterableEventHandlers[T]{
		handler,
		filters,
	})
}

func (d *EventDispatcher[T]) Publish(ctx context.Context, events ...T) error {
	for _, event := range events {
		for _, eventHandler := range d.eventHandlers {
			if eventHandler.filters != nil {

				if _, exists := eventHandler.filters[event.EventName()]; !exists {
					continue
				}
			}
			err := eventHandler.handler.HandleEvent(ctx, event)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
