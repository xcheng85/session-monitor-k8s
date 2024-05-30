package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"github.com/xcheng85/session-monitor-k8s/internal/ddd"
	"github.com/xcheng85/session-monitor-k8s/internal/logger"
	"github.com/xcheng85/session-monitor-k8s/internal/repository"
	"github.com/xcheng85/session-monitor-k8s/internal/session"
	"github.com/xcheng85/session-monitor-k8s/node/internal/domain"
	"go.uber.org/zap"
)

func TestHandleEvent(t *testing.T) {
	scenarios := []struct {
		desc                  string
		inLogger              *zap.Logger
		inConfigMock          func() *config.MockIConfig
		inEventDispatcherMock func() *ddd.MockIEventDispatcher[ddd.IEvent]
		inKVRepositoryMock    func() *repository.MockIKVRepository
		inSessionServiceMock  func() *session.MockISessionService
		inPayload             func() ddd.IEvent
		expectedError         error
	}{
		{
			desc: "Add Node",
			inLogger: logger.NewZapLogger(logger.LogConfig{
				LogLevel: logger.DEBUG,
			}),
			inConfigMock: func() *config.MockIConfig {
				return &config.MockIConfig{}
			},
			inEventDispatcherMock: func() *ddd.MockIEventDispatcher[ddd.IEvent] {
				mockEventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
				mockEventDispatcher.On("Subscribe", mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return &mockEventDispatcher
			},
			inKVRepositoryMock: func() *repository.MockIKVRepository {
				return &repository.MockIKVRepository{}
			},
			inSessionServiceMock: func() *session.MockISessionService {
				return &session.MockISessionService{}
			},
			inPayload: func() ddd.IEvent {
				nodeLables := map[string]string{
					"accelerator": "nvidia",
					"agentpool":   "viz",
				}
				nodeDomain := &domain.Node{
					Name:          "nodeName",
					DriverVersion: "525.0.0",
					Labels:        &nodeLables,
				}
				event := ddd.NewEvent(
					domain.NodeAddEvent,
					&domain.NodeEventPayload{
						Node: nodeDomain,
					})
				return event
			},
			expectedError: nil,
		},
		{
			desc: "Update Node",
			inLogger: logger.NewZapLogger(logger.LogConfig{
				LogLevel: logger.DEBUG,
			}),
			inConfigMock: func() *config.MockIConfig {
				return &config.MockIConfig{}
			},
			inEventDispatcherMock: func() *ddd.MockIEventDispatcher[ddd.IEvent] {
				mockEventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
				mockEventDispatcher.On("Subscribe", mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return &mockEventDispatcher
			},
			inKVRepositoryMock: func() *repository.MockIKVRepository {
				return &repository.MockIKVRepository{}
			},
			inSessionServiceMock: func() *session.MockISessionService {
				return &session.MockISessionService{}
			},
			inPayload: func() ddd.IEvent {
				nodeLables := map[string]string{
					"accelerator": "nvidia",
					"agentpool":   "viz",
				}
				nodeDomain := &domain.Node{
					Name:          "nodeName",
					DriverVersion: "525.0.0",
					Labels:        &nodeLables,
				}
				event := ddd.NewEvent(
					domain.NodeUpdateEvent,
					&domain.NodeEventPayload{
						Node: nodeDomain,
					})
				return event
			},
			expectedError: nil,
		},
		{
			desc: "Delete Node",
			inLogger: logger.NewZapLogger(logger.LogConfig{
				LogLevel: logger.DEBUG,
			}),
			inConfigMock: func() *config.MockIConfig {
				return &config.MockIConfig{}
			},
			inEventDispatcherMock: func() *ddd.MockIEventDispatcher[ddd.IEvent] {
				mockEventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
				mockEventDispatcher.On("Subscribe", mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return &mockEventDispatcher
			},
			inKVRepositoryMock: func() *repository.MockIKVRepository {
				return &repository.MockIKVRepository{}
			},
			inSessionServiceMock: func() *session.MockISessionService {
				return &session.MockISessionService{}
			},
			inPayload: func() ddd.IEvent {
				nodeLables := map[string]string{
					"accelerator": "nvidia",
					"agentpool":   "viz",
				}
				nodeDomain := &domain.Node{
					Name:          "nodeName",
					DriverVersion: "525.0.0",
					Labels:        &nodeLables,
				}
				event := ddd.NewEvent(
					domain.NodeDeleteEvent,
					&domain.NodeEventPayload{
						Node: nodeDomain,
					})
				return event
			},
			expectedError: nil,
		},
		{
			desc: "Update Node Label Cache",
			inLogger: logger.NewZapLogger(logger.LogConfig{
				LogLevel: logger.DEBUG,
			}),
			inConfigMock: func() *config.MockIConfig {
				mockConfig := &config.MockIConfig{}
				mockConfig.On("Get", "app.gpu_agent_pool_set_key").Return("GpuNodePools", nil).Once()
				return mockConfig
			},
			inEventDispatcherMock: func() *ddd.MockIEventDispatcher[ddd.IEvent] {
				mockEventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
				mockEventDispatcher.On("Subscribe", mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return &mockEventDispatcher
			},
			inKVRepositoryMock: func() *repository.MockIKVRepository {
				mockKVRepository := &repository.MockIKVRepository{}
				mockKVRepository.On("AddToUnsortedSet", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)
				return mockKVRepository
			},
			inSessionServiceMock: func() *session.MockISessionService {
				return &session.MockISessionService{}
			},
			inPayload: func() ddd.IEvent {
				nodeLables := map[string]string{
					"accelerator": "nvidia",
					"agentpool":   "viz",
				}
				nodeDomain := &domain.Node{
					Name:          "nodeName",
					DriverVersion: "525.0.0",
					Labels:        &nodeLables,
				}
				event := ddd.NewEvent(
					domain.NodeUpdateLabelsCacheEvent,
					&domain.NodeEventPayload{
						Node: nodeDomain,
					})
				return event
			},
			expectedError: nil,
		},
		{
			desc: "Record Node Provision Timestamp",
			inLogger: logger.NewZapLogger(logger.LogConfig{
				LogLevel: logger.DEBUG,
			}),
			inConfigMock: func() *config.MockIConfig {
				mockConfig := &config.MockIConfig{}
				return mockConfig
			},
			inEventDispatcherMock: func() *ddd.MockIEventDispatcher[ddd.IEvent] {
				mockEventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
				mockEventDispatcher.On("Subscribe", mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return &mockEventDispatcher
			},
			inKVRepositoryMock: func() *repository.MockIKVRepository {
				mockKVRepository := &repository.MockIKVRepository{}
				mockKVRepository.On("GetServerTimestamp", mock.Anything).Return(int64(8888888888), nil)
				return mockKVRepository
			},
			inSessionServiceMock: func() *session.MockISessionService {
				mockSessionService := &session.MockISessionService{}
				mockSessionService.On("SetNodeProvisionTimeStamp", &session.SetNodeProvisionTimeStampActionPayload{
					NodeName:  "nodeName",
					Timestamp: 8888888888,
				}).Return(nil)
				return mockSessionService
			},
			inPayload: func() ddd.IEvent {
				nodeLables := map[string]string{
					"accelerator": "nvidia",
					"agentpool":   "viz",
				}
				nodeDomain := &domain.Node{
					Name:          "nodeName",
					DriverVersion: "525.0.0",
					Labels:        &nodeLables,
				}
				event := ddd.NewEvent(
					domain.NodeRecordNodeProvisionEvent,
					&domain.NodeEventPayload{
						Node: nodeDomain,
					})
				return event
			},
			expectedError: nil,
		},
	}
	for _, s := range scenarios {
		scenario := s
		t.Run(scenario.desc, func(t *testing.T) {
			ctx := context.TODO()
			logger := scenario.inLogger
			config := scenario.inConfigMock()
			eventDispatcher := scenario.inEventDispatcherMock()
			kvRepository := scenario.inKVRepositoryMock()
			sessionService := scenario.inSessionServiceMock()
			payload := scenario.inPayload()

			eventHandler := NewDomainEventHandlers(logger, config, eventDispatcher, kvRepository, sessionService)
			err := eventHandler.HandleEvent(ctx, payload)
			assert.Equal(t, scenario.expectedError, err, scenario.desc)
		})
	}
}
