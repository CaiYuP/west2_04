package repo

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Put(ctx context.Context, key, value string, expire time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Hset(ctx context.Context, key string, field string, value string)
	HKeys(background context.Context, key string) ([]string, error)
	DeleteAll(background context.Context, key string, fields []string)

	// 排行榜（ZSET）相关方法
	ZAdd(ctx context.Context, key string, score float64, member string) error
	ZAddBatch(ctx context.Context, key string, members []*redis.Z) error
	ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error)
	ZCard(ctx context.Context, key string) (int64, error)
	ZRem(ctx context.Context, key string, members ...string) error

	// 访问量计数器相关方法
	IncrVisitCount(ctx context.Context, videoID uint64) (int64, error)
	GetVisitCount(ctx context.Context, videoID uint64) (int64, error)
	GetAllVisitCountKeys(ctx context.Context) ([]string, error)
	GetVisitCountsBatch(ctx context.Context, videoIDs []uint64) (map[uint64]int64, error)
	DelVisitCount(ctx context.Context, videoID uint64) error
	DelVisitCountBatch(ctx context.Context, videoIDs []uint64) error
}
