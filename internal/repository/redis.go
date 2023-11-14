package repository

import (
	"context"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"go.uber.org/zap"
)

type redisRepository struct {
	client *redis.Client
	logger *zap.Logger
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
		ID:     streamId,
		Values: payload,
	})
	return status.Result()
}

func (s *redisRepository) AddToUnsortedSet(ctx context.Context, UnsortedSetKey string, objects ...*Object) (int64, error) {
	cmds, err := s.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		allKeys := []string{}
		for _, obj := range objects {
			key, payload, Expiration := obj.Key, obj.Payload, obj.Expiration
			pipe.Set(ctx, key, payload, Expiration)
			allKeys = append(allKeys, key)
		}
		pipe.SAdd(ctx, UnsortedSetKey, allKeys)
		return nil
	})
	if err != nil {
		return 0, err
	}
	if len(cmds) > 0 {
		lastCmd := cmds[len(cmds)-1]
		numKeyAdded := lastCmd.(*redis.IntCmd).Val()
		s.logger.Sugar().Infof("numKeyAdded: %d", numKeyAdded)
		return numKeyAdded, nil
	}
	return 0, nil
}
