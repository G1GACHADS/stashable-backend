package clients

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	ctx   context.Context
	redis *redis.Client
}

func NewRedisClient(ctx context.Context, opts *redis.Options) (*RedisClient, error) {
	redis := redis.NewClient(opts)
	if _, err := redis.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &RedisClient{ctx: ctx, redis: redis}, nil
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.redis.Get(r.ctx, key).Result()
}

func (r *RedisClient) Set(key string, value any) error {
	return r.redis.Set(r.ctx, key, value, 0).Err()
}

func (r *RedisClient) SetWithTTL(key string, value any, ttl time.Duration) error {
	return r.redis.Set(r.ctx, key, value, ttl).Err()
}

func (r *RedisClient) Delete(key string) error {
	return r.redis.Del(r.ctx, key).Err()
}
