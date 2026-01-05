package dao

import (
	"context"
	"west2-video/common/logs"

	"github.com/go-redis/redis/v8"
	"time"
)

var Rc *RedisCache

type RedisCache struct {
	Rdb *redis.Client
}

func (r *RedisCache) DeleteAll(ctx context.Context, key string, fields []string) {
	err := r.Rdb.Del(ctx, fields...).Err()
	if err != nil {
		logs.LG.Error(err.Error())
	}
}

func (r *RedisCache) HKeys(background context.Context, key string) ([]string, error) {
	return r.Rdb.HKeys(background, key).Result()
}

func (r *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	r.Rdb.Set(ctx, key, value, expire)
	return nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := r.Rdb.Get(ctx, key).Result()
	return result, err
}
func (r *RedisCache) Hset(ctx context.Context, key, field, value string) {
	r.Rdb.HSet(ctx, key, field, value)
}

//func init() {
//	cli := redis.NewClient(
//		config.C.InitRedisOptions(),
//	)
//	Rc = &RedisCache{
//		rdb: cli,
//	}
//}
