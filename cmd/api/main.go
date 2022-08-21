package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/G1GACHADS/stashable-backend/api"
	"github.com/G1GACHADS/stashable-backend/backend"
	"github.com/G1GACHADS/stashable-backend/clients"
	"github.com/G1GACHADS/stashable-backend/config"
	"github.com/G1GACHADS/stashable-backend/core/logger"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger.Init(true)

	ctx, cancel := context.WithCancel(context.Background())

	// Wait for kill signals to gracefully shutdown the server
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-c

		cancel()
	}()

	config := config.New()
	clients, err := clients.New(ctx, config)
	if err != nil {
		logger.M.Fatal(err.Error())
	}

	backend := backend.New(clients, config)
	srv := api.NewServer(backend, config)

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return srv.Listen(config.App.Address)
	})

	group.Go(func() error {
		<-groupCtx.Done()
		logger.M.Warn("Shutting down server")
		return srv.Shutdown()
	})

	if err := group.Wait(); err != nil {
		logger.M.Warnf("Exit reason: %v\n", err)
	}
}
