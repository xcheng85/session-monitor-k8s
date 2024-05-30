package repository

import (
	"context"
	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"go.uber.org/zap"
)

type redisRepository struct {
	clients []IKVRepository
	logger  *zap.Logger
}

func NewRedisRepository(ctx context.Context, config config.IConfig, logger *zap.Logger) (IKVRepository, error) {
	redis_address := config.Get("app.redis_address").(string)
	redisv9_address := config.Get("app.redisv9_address").(string)
	redis_mock := config.Get("app.redis_mock").(bool)
	redis_disable := config.Get("app.redis_disable").(bool)
	redisv9_disable := config.Get("app.redisv9_disable").(bool)

	clients := []IKVRepository{}
	var clientv8 *redisClientV8 = nil
	var clientv9 *redisClientV9 = nil
	var err error
	// reverse the order, so client v8 has high priority in the testing phase.
	if !redisv9_disable {
		clientv9, err = newRedisClientV9(redisv9_address, redis_mock, logger)
		if err != nil {
			return nil, err
		}
		clients = append(clients, clientv9)
		logger.Sugar().Infof("Redis Client v9 is appended!")
	}
	if !redis_disable {
		clientv8, err = newRedisClientV8(redis_address, redis_mock, logger)
		if err != nil {
			return nil, err
		}
		clients = append(clients, clientv8)
		logger.Sugar().Infof("Redis Client v8 is appended!")
	}
	redis_repo := &redisRepository{
		clients,
		logger,
	}
	resp, err := redis_repo.Ping(ctx)
	if err != nil {
		logger.Sugar().Errorf("an error '%s' occurred when opening a mock redis connection", err)
		return nil, err
	} else {
		logger.Sugar().Debugf("resp from redis server: '%s'", resp)
		return redis_repo, nil
	}
}

func (s *redisRepository) GetServerTimestamp(ctx context.Context) (unixTimeStamp int64, err error) {
	for _, client := range s.clients {
		unixTimeStamp, err = client.GetServerTimestamp(ctx)
	}
	return unixTimeStamp, err
}

func (s *redisRepository) Ping(ctx context.Context) (message string, err error) {
	for _, client := range s.clients {
		message, err = client.Ping(ctx)
	}
	return message, err
}

func (s *redisRepository) AddStreamEvent(ctx context.Context, streamKey string, streamId string, payload interface{}) (message string, err error) {
	for _, client := range s.clients {
		message, err = client.AddStreamEvent(ctx, streamKey, streamId, payload)
	}
	return message, err
}

func (s *redisRepository) AddToUnsortedSet(ctx context.Context, UnsortedSetKey string, objects ...*Object) (keyUpdated int64, err error) {
	for _, client := range s.clients {
		keyUpdated, err = client.AddToUnsortedSet(ctx, UnsortedSetKey, objects...)
	}
	return keyUpdated, err
}
