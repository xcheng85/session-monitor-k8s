package domain

const PodAggregate = "pods.CustomerAggregate"

type Pod struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	SessionId string `json:"sessionId,omitempty"`
	NodeName  string `json:"nodeName,omitempty"`
	Ip        string `json:"ip,omitempty"`
}
