package domain

const (
	NodeNilEvent                 = "NodeNilEvent"
	NodeInformerErrorEvent       = "NodeInformerErrorEvent"
	NodeAddEvent                 = "NodeAddEvent"
	NodeDeleteEvent              = "NodeDeleteEvent"
	NodeUpdateEvent              = "NodeUpdateEvent"
	NodeRecordNodeProvisionEvent = "NodeRecordNodeProvisionEvent"
)

type NodeInformerErrorPayload struct {
	Err error
}

type NodeEventPayload struct {
	Name string `json:"name,omitempty"`
}
