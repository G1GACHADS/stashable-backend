package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/G1GACHADS/backend/internal/backend"
	"github.com/G1GACHADS/backend/internal/clients"
	"github.com/G1GACHADS/backend/internal/config"
	"github.com/G1GACHADS/backend/internal/logger"
	"github.com/bxcodec/faker/v3"
	_ "github.com/joho/godotenv/autoload"
)

type Limiter struct {
	limit chan struct{}
	wg    sync.WaitGroup
}

func NewLimiter(n int) *Limiter {
	return &Limiter{limit: make(chan struct{}, n)}
}

func (lim *Limiter) Go(ctx context.Context, fn func()) bool {
	if ctx.Err() != nil {
		return false
	}

	select {
	case lim.limit <- struct{}{}:
	case <-ctx.Done():
		return false
	}

	lim.wg.Add(1)
	go func() {
		defer func() {
			<-lim.limit
			lim.wg.Done()
		}()

		fn()
	}()

	return true
}

func (lim *Limiter) Wait() {
	lim.wg.Wait()
}

const (
	n = 100
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

	b := backend.New(clients, config)

	// b.CreateCategory(ctx, "Chemical")
	// b.CreateCategory(ctx, "Electric Components")
	// b.CreateCategory(ctx, "Fragile / Glass")
	// b.CreateCategory(ctx, "Heavy Materials")

	minBasePrice := float64(200000)
	maxBasePrice := float64(25000000)

	limiter := NewLimiter(8)
	for i := 0; i < n; i++ {
		started := limiter.Go(ctx, func() {
			err := b.CreateWarehouse(ctx, backend.CreateWarehouseInput{
				Warehouse: backend.Warehouse{
					Name:        "PT. " + faker.Word(),
					ImageURL:    "https://source.unsplash.com/random/800x800",
					Description: faker.Paragraph(),
					BasePrice:   minBasePrice + rand.Float64()*(maxBasePrice-minBasePrice),
					Email:       faker.Email(),
					PhoneNumber: faker.Phonenumber(),
				},
				Address: backend.Address{
					Province:   faker.Word(),
					City:       faker.Word(),
					StreetName: faker.Sentence(),
					ZipCode:    rand.Intn(180000-170000) + 170000,
				},
				CategoryIDs: []int64{1, 2, 3, 4},
			})
			if err != nil {
				logger.M.Warnf("Worker-#%d: failed inserting\nreason:%v", i, err)
			}
		})

		if !started {
			logger.M.Fatal(ctx.Err())
		}
	}

	limiter.Wait()
	logger.M.Info("Database populate process finished!")
}
