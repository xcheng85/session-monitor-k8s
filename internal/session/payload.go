package session

type StreamTaskType string

// must match the definition in "https://dev.azure.com/slb-swt/slbCloud3DViz/_git/viz-3d-service-infrastructure?path=/src/models/redis/stream.ts"
const (
	EnqueueSession StreamTaskType = "EnqueueSession"
	DeleteSession  StreamTaskType = "DeleteSession"
)

type SetNodeProvisionTimeStampActionPayload struct {
	NodeName  string
	Timestamp int64
}

type UpdateSessionTimeStampLikeFieldActionPayload struct {
	SessionId string
	Timestamp int64
}

// must match the definition in "https://dev.azure.com/slb-swt/slbCloud3DViz/_git/viz-3d-service-infrastructure?path=/src/models/redis/stream.ts"
type SetSessionReadyActionPayload struct {
	SessionId              string `json:"sessionId" binding:"required"`
	NodeName               string `json:"nodeName" binding:"required"`
	NodeProvisionTimeStamp int64  `json:"nodeProvisionTimeStamp"`
	PodScheduleTimeStamp   int64  `json:"podScheduleTimeStamp"`
	PodInternalIp          string `json:"podInternalIp"`
}

// must match the definition in "https://dev.azure.com/slb-swt/slbCloud3DViz/_git/viz-3d-service-infrastructure?path=/src/models/redis/stream.ts"
type SetSessionDeletableActionPayload struct {
	SessionId string `json:"sessionId" binding:"required"`
}
