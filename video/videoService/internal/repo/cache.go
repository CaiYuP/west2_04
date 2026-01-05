package repo

import (
	"context"
	"time"
)

type Cache interface {
	Put(ctx context.Context, key, value string, expire time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Hset(ctx context.Context, key string, field string, value string)
	HKeys(background context.Context, key string) ([]string, error)
	DeleteAll(background context.Context, key string, fields []string)
}
