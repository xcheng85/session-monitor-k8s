package session

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"github.com/xcheng85/session-monitor-k8s/internal/repository"
	"go.uber.org/zap"
)

//go:generate mockery --name ISessionService
type ISessionService interface {
	SetSessionReady(*SetSessionReadyActionPayload) error
	SetSessionDeletable(*SetSessionDeletableActionPayload) error
	SetNodeProvisionTimeStamp(*SetNodeProvisionTimeStampActionPayload) error
	SetPodScheduleTimeStamp(*UpdateSessionTimeStampLikeFieldActionPayload) error
	GetNodeProvisionTimeStamp(NodeName string) (int64, error)
	GetPodScheduleTimeStamp(sessionId string) (int64, error)
}

type sessionService struct {
	ctx    context.Context
	logger *zap.Logger
	config config.IConfig
	kvRepo repository.IKVRepository
}

var _ ISessionService = (*sessionService)(nil)

func NewSessionService(ctx context.Context, logger *zap.Logger, config config.IConfig, kvRepo repository.IKVRepository) ISessionService {
	return &sessionService{
		ctx,
		logger,
		config,
		kvRepo,
	}
}

func (svc *sessionService) SetSessionReady(payload *SetSessionReadyActionPayload) error {
	// reuse viper as config store
	streamKey := svc.config.Get("app.enqueue_session_stream_key").(string)
	out, _ := json.Marshal(payload)
	svc.logger.Sugar().Info(string(out))
	currentServerUnixTimestamp, err := svc.kvRepo.GetServerTimestamp(svc.ctx)
	if err != nil {
		return err
	}
	payloadToKvStore := []interface{}{"TaskType", string(EnqueueSession), "TaskInfo", string(out), "TaskCreateTimeStamp", currentServerUnixTimestamp}
	streamId, err := svc.kvRepo.AddStreamEvent(svc.ctx, streamKey, "*", payloadToKvStore)
	svc.logger.Sugar().Infof("SetSessionReady create streamTask: %s", streamId)
	return err
}

func (svc *sessionService) SetSessionDeletable(payload *SetSessionDeletableActionPayload) error {
	// reuse viper as config store
	streamKey := svc.config.Get("app.delete_session_stream_key").(string)
	out, _ := json.Marshal(payload)
	svc.logger.Sugar().Info(string(out))
	currentServerUnixTimestamp, err := svc.kvRepo.GetServerTimestamp(svc.ctx)
	if err != nil {
		return err
	}
	payloadToKvStore := []interface{}{"TaskType", string(DeleteSession), "TaskInfo", string(out), "TaskCreateTimeStamp", currentServerUnixTimestamp}
	streamId, err := svc.kvRepo.AddStreamEvent(svc.ctx, streamKey, "*", payloadToKvStore)
	svc.logger.Sugar().Infof("SetSessionDeletable create streamTask: %s", streamId)
	return err
}

func (svc *sessionService) SetNodeProvisionTimeStamp(payload *SetNodeProvisionTimeStampActionPayload) error {
	nodeName, timestamp := payload.NodeName, payload.Timestamp
	nodeProvisionTimestampStoreKey := fmt.Sprintf("NodeProvisionTimeStamp.'%s'", nodeName)
	svc.logger.Sugar().Infow("SetNodeProvisionTimeStamp", "nodeName", nodeName, "timestamp", timestamp, "key", nodeProvisionTimestampStoreKey)
	svc.config.Set(nodeProvisionTimestampStoreKey, timestamp)
	return nil
}

func (svc *sessionService) SetPodScheduleTimeStamp(payload *UpdateSessionTimeStampLikeFieldActionPayload) error {
	sessionId, timestamp := payload.SessionId, payload.Timestamp
	podScheduledTimestampStoreKey := fmt.Sprintf("PodScheduleTimeStamp.'%s'", sessionId)
	svc.logger.Sugar().Infow("SetNodeProvisionTimeStamp", "sessionId", sessionId, "timestamp", timestamp, "key", podScheduledTimestampStoreKey)
	svc.config.Set(podScheduledTimestampStoreKey, timestamp)
	return nil
}

func (svc *sessionService) GetNodeProvisionTimeStamp(nodeName string) (int64, error) {
	nodeProvisionTimestampStoreKey := fmt.Sprintf("NodeProvisionTimeStamp.'%s'", nodeName)
	svc.logger.Sugar().Infow("GetNodeProvisionTimeStamp", "key", nodeProvisionTimestampStoreKey)
	val, ok := svc.config.Get(nodeProvisionTimestampStoreKey).(int64)
	if ok {
		return val, nil
	} else {
		return 0, NewInvalidStoreKeyErr(nodeProvisionTimestampStoreKey)
	}
}

func (svc *sessionService) GetPodScheduleTimeStamp(sessionId string) (int64, error) {
	podScheduledTimestampStoreKey := fmt.Sprintf("PodScheduleTimeStamp.'%s'", sessionId)
	svc.logger.Sugar().Infow("GetPodScheduleTimeStamp", "key", podScheduledTimestampStoreKey)
	val, ok := svc.config.Get(podScheduledTimestampStoreKey).(int64)
	if ok {
		return val, nil
	} else {
		return 0, NewInvalidStoreKeyErr(podScheduledTimestampStoreKey)
	}
}
