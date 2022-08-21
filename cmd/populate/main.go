package main

import (
	"context"
	"math/rand"

	"github.com/G1GACHADS/stashable-backend/backend"
	"github.com/G1GACHADS/stashable-backend/clients"
	"github.com/G1GACHADS/stashable-backend/config"
	"github.com/G1GACHADS/stashable-backend/core/logger"
	"github.com/bxcodec/faker/v3"
	_ "github.com/joho/godotenv/autoload"
)

const (
	n = 100
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

	categoryChemical, _ := b.CreateCategory(ctx, "Chemical")
	categoryElectricComponents, _ := b.CreateCategory(ctx, "Electric Components")
	categoryFragileGlass, _ := b.CreateCategory(ctx, "Fragile / Glass")
	categoryHeavyMaterials, _ := b.CreateCategory(ctx, "Heavy Materials")
	categories := []int64{
		categoryChemical.ID,
		categoryElectricComponents.ID,
		categoryFragileGlass.ID,
		categoryHeavyMaterials.ID}

	minBasePrice := float64(200000)
	maxBasePrice := float64(25000000)

	limiter := NewLimiter(12)
	for i := 0; i < n; i++ {
		started := limiter.Go(ctx, i, func(id int) {
			// Generate category ids with random length and unique random values
			categoryIDs := make([]int64, rand.Intn(len(categories))+1)
			for i := 0; i < len(categoryIDs); i++ {
				categoryIDs[i] = categories[rand.Intn(len(categories))]
			}
			categoryIDs = RemoveDuplicates(categoryIDs)

			err := b.CreateWarehouse(ctx, backend.CreateWarehouseInput{
				Warehouse: backend.Warehouse{
					Name:        "PT. " + faker.Word() + faker.Word(),
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
				CategoryIDs: categoryIDs,
			})
			if err != nil {
				logger.M.Warnf("Worker-#%d: failed inserting\nreason:%v", id, err)
				logger.M.Warnf("Category IDs: %v", categoryIDs)
			}
		})

		if !started {
			logger.M.Fatal(ctx.Err())
		}
	}

	limiter.Wait()
	logger.M.Info("Database populate process finished!")
}
