package domain

const (
	NodeNilEvent                 = "NodeNilEvent"
	NodeInformerErrorEvent       = "NodeInformerErrorEvent"
	NodeAddEvent                 = "NodeAddEvent"
	NodeDeleteEvent              = "NodeDeleteEvent"
	NodeUpdateEvent              = "NodeUpdateEvent"
	NodeRecordNodeProvisionEvent = "NodeRecordNodeProvisionEvent"
	NodeUpdateLabelsCacheEvent   = "NodeUpdateLabelsCacheEvent"
)

type NodeInformerErrorPayload struct {
	Err error
}

type NodeEventPayload struct {
	Node *Node
}
