package clients

import (
	"context"

	"github.com/G1GACHADS/stashable-backend/config"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/sync/errgroup"
)

type Clients struct {
	DB    *pgxpool.Pool
	Cache *redis.Client
}

func New(ctx context.Context, cfg *config.Config) (*Clients, error) {
	var group errgroup.Group

	c := &Clients{}

	group.Go(func() error {
		var err error

		c.DB, err = NewPostgreSQLClient(ctx, cfg.Clients.DatabaseURL)
		if err != nil {
			return err
		}

		return nil
	})

	group.Go(func() error {
		opts, err := redis.ParseURL(cfg.Clients.RedisAddress)
		if err != nil {
			return err
		}

		opts.ReadTimeout = cfg.Clients.RedisReadTimeout
		opts.WriteTimeout = cfg.Clients.RedisWriteTimeout

		c.Cache, err = NewRedisClient(ctx, opts)
		if err != nil {
			return err
		}

		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return c, nil
}
