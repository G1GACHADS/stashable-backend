package clients

import (
	"context"
	"fmt"
	"net/url"

	"github.com/G1GACHADS/stashable-backend/internal/config"
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

		connString := url.URL{
			Scheme: "postgres",
			Host:   fmt.Sprintf("%s:%d", cfg.Clients.PostgresHost, cfg.Clients.PostgresPort),
			User:   url.UserPassword(cfg.Clients.PostgresUser, cfg.Clients.PostgresPassword),
			Path:   cfg.Clients.PostgresDB,
		}

		c.DB, err = NewPostgreSQLClient(ctx, connString.String())
		if err != nil {
			return err
		}

		return nil
	})

	group.Go(func() error {
		var err error

		c.Cache, err = NewRedisClient(ctx, &redis.Options{
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
