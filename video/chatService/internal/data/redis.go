package data

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
	"west2-video/chatService/configs"
)

// RedisClient Redis客户端
type RedisClient struct {
	*redis.Client
}

// NewRedisClient 创建新的Redis客户端
func NewRedisClient(cfg configs.RedisConfig) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic("failed to connect to redis: " + err.Error())
	}

	return &RedisClient{client}
}
