package ddd

import (
	"github.com/google/uuid"
	"time"
)

type EventPayload interface{}
type IEvent interface {
	IDer
	EventName() string
	Payload() EventPayload
	Metadata() Metadata // for event sourcing and aggregation purpose
	OccurredAt() time.Time
}
type event struct {
	Entity
	payload    EventPayload
	metadata   Metadata
	occurredAt time.Time
}

var _ IEvent = (*event)(nil)

func NewEvent(name string, payload EventPayload) IEvent {
	return &event{
		Entity:     NewEntity(uuid.New().String(), name),
		payload:    payload,
		metadata:   make(Metadata),
		occurredAt: time.Now(),
	}
}

func (e event) EventName() string     { return e.EntityName() }
func (e event) Payload() EventPayload { return e.payload }
func (e event) Metadata() Metadata    { return e.metadata }
func (e event) OccurredAt() time.Time { return e.occurredAt }
