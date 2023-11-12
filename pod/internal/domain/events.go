package domain

const (
	PodNilEvent               = "PodNilEvent"
	PodInformerErrorEvent     = "PodInformerErrorEvent"
	PodAddEvent               = "PodAddEvent"
	PodDeleteEvent            = "PodDeleteEvent"
	PodReadyEvent             = "PodReadyEvent"
	PodRecordPodScheduleEvent = "PodRecordPodScheduleEvent"
)

type PodInformerErrorPayload struct {
	Err error
}

type PodEventPayload struct {
	Pod *Pod
}
