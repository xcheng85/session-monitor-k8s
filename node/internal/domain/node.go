package domain

import (
	"github.com/xcheng85/session-monitor-k8s/internal/ddd"
)

const NodeAggregate = "nodes.CustomerAggregate"

type Node struct {
	ddd.IAggregatable
	Name string `json:"name,omitempty"`
}
