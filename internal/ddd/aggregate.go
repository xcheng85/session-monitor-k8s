package ddd

// I of SOLID: interface segration
type IEventManager interface {
	AddEvent(string, EventPayload)
	Events() []IAggregateEvent
	ClearEvents()
}

// IAggregatable: event sourcing for any domain object
type IAggregatable interface {
	IDer
	IEventManager
}

// IAggregatable extends entity
type aggregatable struct {
	Entity
	events []IAggregateEvent
}

var _ IAggregatable = (*aggregatable)(nil)

// entity could be pod, node
// events are attached to the entity correspondingly
func NewAggregatable(id, name string) IAggregatable {
	return &aggregatable{
		Entity: NewEntity(id, name),
		events: []IAggregateEvent{},
	}
}

func (a *aggregatable) Events() []IAggregateEvent {
	return a.events
}
func (a *aggregatable) ClearEvents() {
	a.events = []IAggregateEvent{}
}
func (a *aggregatable) AddEvent(name string, payload EventPayload) {
	a.events = append(
		a.events,
		aggregateEvent{
			event: NewEvent(name, payload).(event),
		},
	)
}

type IAggregateEvent interface {
	IEvent
	AggregateName() string
	AggregateID() string
	AggregateVersion() int
}

type aggregateEvent struct {
	event
}

var _ IAggregateEvent = (*aggregateEvent)(nil)

// predefined key for event aggregation
const (
	AggregateNameKey    = "aggregate-name"
	AggregateIDKey      = "aggregate-id"
	AggregateVersionKey = "aggregate-version"
)

func (e aggregateEvent) AggregateName() string { return e.metadata.Get(AggregateNameKey).(string) }
func (e aggregateEvent) AggregateID() string   { return e.metadata.Get(AggregateIDKey).(string) }
func (e aggregateEvent) AggregateVersion() int { return e.metadata.Get(AggregateVersionKey).(int) }
