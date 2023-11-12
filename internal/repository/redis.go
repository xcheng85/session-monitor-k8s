package repository

import (
	"context"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"go.uber.org/zap"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(address string, mock bool, logger *zap.Logger) (client *redis.Client, err error) {
	if mock {
		mr, err := miniredis.Run()
		if err != nil {
			return nil, err
		}
		client = redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     address,
			Password: "",
			DB:       0,
		})
	}
	return client, nil
}

func NewRedisRepository(ctx context.Context, config config.IConfig, logger *zap.Logger) (IKVRepository, error) {
	redis_address := config.Get("app.redis_address").(string)
	redis_mock := config.Get("app.redis_mock").(bool)
	client, err := newRedisClient(redis_address, redis_mock, logger)
	if err != nil {
		return nil, err
	}

	redis_repo := &redisRepository{
		client,
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
	status := s.client.Time(ctx)
	unixTimeStamp, err = status.Val().Unix(), status.Err()
	return unixTimeStamp, err
}

func (s *redisRepository) Ping(ctx context.Context) (string, error) {
	status := s.client.Ping(ctx)
	return status.Result()
}

func (s *redisRepository) AddStreamEvent(ctx context.Context, streamKey string, streamId string, payload interface{}) (string, error) {
	status := s.client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		ID:     "*",
		Values: payload,
	})

	return status.Result()
}

func (s *redisRepository) AddToUnsortedSet(ctx context.Context, UnsortedSetKey string, payload ...interface{}) (int64, error) {
	status := s.client.SAdd(ctx, UnsortedSetKey, payload)
	return status.Result()
}

func (s *redisRepository) Set(ctx context.Context, key string, payload interface{}, expiration time.Duration) (string, error) {
	status := s.client.Set(ctx, key, payload, expiration)
	return status.Result()
}
