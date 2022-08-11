package clients

import (
	"context"
	"fmt"

	"github.com/G1GACHADS/backend/internal/config"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/errgroup"
)

type Clients struct {
	Postgres *PostgreSQLConnection
	Redis    *RedisClient
}

func New(ctx context.Context, cfg *config.Config) (*Clients, error) {
	group, ctx := errgroup.WithContext(ctx)

	c := &Clients{}

	group.Go(func() error {
		var err error

		connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.Clients.PostgresUser,
			cfg.Clients.PostgresPassword,
			cfg.Clients.PostgresHost,
			cfg.Clients.PostgresPort,
			cfg.Clients.PostgresDB)

		c.Postgres, err = NewPostgreSQLClient(ctx, connString)
		if err != nil {
			return err
		}

		return nil
	})

	group.Go(func() error {
		var err error

		c.Redis, err = NewRedisClient(ctx, &redis.Options{
			Network:      "tcp",
			Addr:         cfg.Clients.RedisAddress,
			DB:           cfg.Clients.RedisDB,
			Password:     cfg.Clients.RedisPassword,
			ReadTimeout:  cfg.Clients.RedisReadTimeout,
			WriteTimeout: cfg.Clients.RedisWriteTimeout,
		})
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
