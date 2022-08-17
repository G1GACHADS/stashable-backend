package main

import (
	"context"
	"math/rand"
	"sync"

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

func generateUniqueRandomSlice(n int64) []int64 {
	s := make([]int64, n)
	for i := int64(1); i <= n; i++ {
		s[i-1] = i
	}

	return s
}

const (
	n = 1000
)

func main() {
	logger.Init(true)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := config.New()
	clients, err := clients.New(ctx, config)
	if err != nil {
		logger.M.Fatal(err.Error())
	}

	b := backend.New(clients, config)

	b.CreateCategory(ctx, "Chemical")
	b.CreateCategory(ctx, "Electric Components")
	b.CreateCategory(ctx, "Fragile / Glass")
	b.CreateCategory(ctx, "Heavy Materials")

	minBasePrice := float64(200000)
	maxBasePrice := float64(25000000)

	limiter := NewLimiter(12)
	for i := 0; i < n; i++ {
		started := limiter.Go(ctx, func() {
			generatedCategoryIDs := generateUniqueRandomSlice(rand.Int63n(3) + 1)
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
					ZipCode:    rand.Intn(18000-17000) + 17000,
				},
				CategoryIDs: generatedCategoryIDs,
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
