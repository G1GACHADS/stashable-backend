package clients

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(ctx context.Context, opts *redis.Options) (*redis.Client, error) {
	redis := redis.NewClient(opts)
	if _, err := redis.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return redis, nil
}
