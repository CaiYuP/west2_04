package dao

import (
	"context"
	"fmt"
	"strconv"
	"west2-video/common/logs"

	"github.com/go-redis/redis/v8"
	"time"
)

var Rc *RedisCache

const (
	VideoRankingKey      = "video:ranking:visit"
	VideoVisitKeyPrefix  = "video:visit:"
	VideoDetailKeyPrefix = "video:detail:"
)

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
	if err != nil {
		if err == redis.Nil {
			return "", nil
		} else {
			return "", err
		}
	}
	return result, err
}
func (r *RedisCache) Hset(ctx context.Context, key, field, value string) {
	r.Rdb.HSet(ctx, key, field, value)
}

func (r *RedisCache) ZAdd(ctx context.Context, key string, score float64, member string) error {
	return r.Rdb.ZAdd(ctx, key, &redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

func (r *RedisCache) ZAddBatch(ctx context.Context, key string, members []*redis.Z) error {
	return r.Rdb.ZAdd(ctx, key, members...).Err()
}

func (r *RedisCache) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	result, err := r.Rdb.ZRevRange(ctx, key, start, stop).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return result, err
}

func (r *RedisCache) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	result, err := r.Rdb.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return result, err
}

//总数
func (r *RedisCache) ZCard(ctx context.Context, key string) (int64, error) {
	result, err := r.Rdb.ZCard(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		} else {
			return 0, err
		}
	}
	return result, err
}

// 删除排行榜成员
func (r *RedisCache) ZRem(ctx context.Context, key string, members ...string) error {
	args := make([]interface{}, len(members))
	for i, member := range members {
		args[i] = member
	}
	return r.Rdb.ZRem(ctx, key, args...).Err()
}

// 增加视频访问量（原子操作）
func (r *RedisCache) IncrVisitCount(ctx context.Context, videoID uint64) (int64, error) {
	key := fmt.Sprintf(VideoVisitKeyPrefix+"%d", videoID)
	return r.Rdb.Incr(ctx, key).Result()
}

func (r *RedisCache) GetVisitCount(ctx context.Context, videoID uint64) (int64, error) {
	key := fmt.Sprintf(VideoVisitKeyPrefix+"%d", videoID)
	result, err := r.Rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	count, err := strconv.ParseInt(result, 10, 64)
	return count, err
}

//获取所有访问量计数器键
func (r *RedisCache) GetAllVisitCountKeys(ctx context.Context) ([]string, error) {
	return r.Rdb.Keys(ctx, "video:visit:*").Result()
}

// 批量获取访问量
func (r *RedisCache) GetVisitCountsBatch(ctx context.Context, videoIDs []uint64) (map[uint64]int64, error) {
	pipe := r.Rdb.Pipeline()
	cmds := make(map[uint64]*redis.StringCmd)

	for _, videoID := range videoIDs {
		key := fmt.Sprintf(VideoVisitKeyPrefix+"%d", videoID)
		cmds[videoID] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	result := make(map[uint64]int64)
	for videoID, cmd := range cmds {
		val, err := cmd.Result()
		if err == redis.Nil {
			result[videoID] = 0
		} else if err != nil {
			return nil, err
		} else {
			count, _ := strconv.ParseInt(val, 10, 64)
			result[videoID] = count
		}
	}

	return result, nil
}

func (r *RedisCache) DelVisitCount(ctx context.Context, videoID uint64) error {
	key := fmt.Sprintf(VideoVisitKeyPrefix+"%d", videoID)
	return r.Rdb.Del(ctx, key).Err()
}

func (r *RedisCache) DelVisitCountBatch(ctx context.Context, videoIDs []uint64) error {
	if len(videoIDs) == 0 {
		return nil
	}

	keys := make([]string, len(videoIDs))
	for i, videoID := range videoIDs {
		keys[i] = fmt.Sprintf(VideoVisitKeyPrefix+"%d", videoID)
	}
	return r.Rdb.Del(ctx, keys...).Err()
}
