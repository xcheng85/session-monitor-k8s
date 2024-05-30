package repository

import (
	"context"
	"github.com/alicebob/miniredis"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type redisClientV9 struct {
	client *redis.Client
	logger *zap.Logger
}

func newRedisClientV9(address string, mock bool, logger *zap.Logger) (*redisClientV9, error) {
	var client *redis.Client
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
	return &redisClientV9{
		client,
		logger,
	}, nil
}

func (s *redisClientV9) GetServerTimestamp(ctx context.Context) (unixTimeStamp int64, err error) {
	status := s.client.Time(ctx)
	unixTimeStamp, err = status.Val().Unix(), status.Err()
	return unixTimeStamp, err
}

func (s *redisClientV9) Ping(ctx context.Context) (message string, err error) {
	message, err = s.client.Ping(ctx).Result()
	if err != nil {
		return message, err
	}
	return message, err
}

func (s *redisClientV9) AddStreamEvent(ctx context.Context, streamKey string, streamId string, payload interface{}) (message string, err error) {
	message, err = s.client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		ID:     streamId,
		Values: payload,
	}).Result()
	return message, err
}

func (s *redisClientV9) AddToUnsortedSet(ctx context.Context, UnsortedSetKey string, objects ...*Object) (int64, error) {
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
