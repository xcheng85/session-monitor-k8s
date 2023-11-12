package domain

import (
	"github.com/xcheng85/session-monitor-k8s/internal/ddd"
)

const PodAggregate = "pods.CustomerAggregate"

type Pod struct {
	ddd.IAggregatable
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	SessionId string `json:"sessionId,omitempty"`
	NodeName  string `json:"nodeName,omitempty"`
	Ip        string `json:"ip,omitempty"`
}

func NewPod(id string) *Pod {
	return &Pod{
		IAggregatable: ddd.NewAggregatable(id, PodAggregate),
	}
}

func RegisterPod(name, namespace string) (*Pod, error) {
	pod := NewPod(name)
	pod.Name = name
	pod.Namespace = namespace
	pod.AddEvent(PodAddEvent, &PodEventPayload{
		Pod: pod,
	})
	return pod, nil
}
