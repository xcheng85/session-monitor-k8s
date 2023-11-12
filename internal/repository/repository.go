package repository

import (
	"context"
	"time"
)

// interface of kv store and sql-like are quite different
type IKVRepository interface {
	GetServerTimestamp(ctx context.Context) (int64, error)
	Ping(ctx context.Context) (string, error)
	AddStreamEvent(ctx context.Context, streamKey string, streamId string, payload interface{}) (string, error)
	AddToUnsortedSet(ctx context.Context, UnsortedSetKey string, payload ...interface{}) (int64, error)
	Set(ctx context.Context, key string, payload interface{}, expiration time.Duration) (string, error)
}
