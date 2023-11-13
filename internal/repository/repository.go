package repository

import (
	"context"
	"time"
)

type Object struct {
	Key        string
	Payload    interface{}
	Expiration time.Duration
}

// interface of kv store and sql-like are quite different
//
//go:generate mockery --name IKVRepository
type IKVRepository interface {
	GetServerTimestamp(ctx context.Context) (int64, error)
	Ping(ctx context.Context) (string, error)
	AddStreamEvent(ctx context.Context, streamKey string, streamId string, payload interface{}) (string, error)
	AddToUnsortedSet(ctx context.Context, UnsortedSetKey string, objects ...*Object) (int64, error)
}
