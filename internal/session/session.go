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
	SetNodeProvisionTimeStamp(*SetNodeProvisionTimeStampActionPayload) error
	SetSessionReady(*SetSessionReadyActionPayload) error
	SetSessionDeletable(*SetSessionDeletableActionPayload) error
	SetPodScheduleTimeStamp(*UpdateSessionTimeStampLikeFieldActionPayload) error
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

	// nodeName, sessionId := payload.NodeName, payload.SessionId
	// nodeProvisionStoreKey := fmt.Sprintf("NodeProvisionTimeStamp.'%s'", nodeName)
	// svc.logger.Sugar().Info(nodeProvisionStoreKey)
	// nodeProvisionedTimeStamp := svc.config.Get(nodeProvisionStoreKey).(int64)

	// podScheduledStoreKey := fmt.Sprintf("PodScheduleTimeStamp.'%s'", sessionId)
	// svc.logger.Sugar().Info(podScheduledStoreKey)
	// podScheduledTimeStamp := svc.config.Get(podScheduledStoreKey).(int64)

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
	svc.logger.Sugar().Infow("SetNodeProvisionTimeStamp", "nodeName", nodeName, "timestamp", timestamp)
	nodeProvisionTimestampStoreKey := fmt.Sprintf("NodeProvisionTimeStamp.'%s'", nodeName)
	svc.logger.Sugar().Info(nodeProvisionTimestampStoreKey)
	svc.config.Set(nodeProvisionTimestampStoreKey, timestamp)
	return nil
}

func (svc *sessionService) SetPodScheduleTimeStamp(payload *UpdateSessionTimeStampLikeFieldActionPayload) error {
	sessionId, timestamp := payload.SessionId, payload.Timestamp
	svc.logger.Sugar().Infow("SetNodeProvisionTimeStamp", "sessionId", sessionId, "timestamp", timestamp)
	podScheduledTimestampStoreKey := fmt.Sprintf("PodScheduleTimeStamp.'%s'", sessionId)
	svc.logger.Sugar().Info(podScheduledTimestampStoreKey)
	svc.config.Set(podScheduledTimestampStoreKey, timestamp)
	return nil
}
