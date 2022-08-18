package backend

import (
	"context"
	"errors"
	"fmt"

	"github.com/G1GACHADS/backend/logger"
	"github.com/bytedance/sonic"
	"github.com/jackc/pgx/v4"
)

type GetWarehouseOutput struct {
	Attributes    Warehouse `json:"attributes"`
	Relationships struct {
		Address    Address  `json:"address"`
		Categories []string `json:"categories"`
	} `json:"relationships"`
}

func (b *backend) GetWarehouse(ctx context.Context, warehouseID int64) (GetWarehouseOutput, error) {
	var warehouse GetWarehouseOutput

	err := b.clients.DB.QueryRow(ctx, "SELECT * FROM warehouses_list WHERE w_id = $1", warehouseID).Scan(
		nil,
		&warehouse.Attributes.ID,
		&warehouse.Attributes.AddressID,
		&warehouse.Attributes.Name,
		&warehouse.Attributes.ImageURL,
		&warehouse.Attributes.Description,
		&warehouse.Attributes.BasePrice,
		&warehouse.Attributes.Email,
		&warehouse.Attributes.PhoneNumber,
		&warehouse.Attributes.CreatedAt,
		&warehouse.Relationships.Address.ID,
		&warehouse.Relationships.Address.Province,
		&warehouse.Relationships.Address.City,
		&warehouse.Relationships.Address.StreetName,
		&warehouse.Relationships.Address.ZipCode,
		&warehouse.Relationships.Categories,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GetWarehouseOutput{}, ErrWarehouseDoesNotExists
		}
		return GetWarehouseOutput{}, err
	}

	go func(warehouse GetWarehouseOutput) {
		cacheKey := fmt.Sprintf("warehouses::%d", warehouseID)
		if exists, _ := b.clients.Cache.Exists(ctx, cacheKey).Result(); exists != 1 {
			// Cache the warehouses for future use
			out, _ := sonic.Marshal(warehouse)
			_, err := b.clients.Cache.Set(ctx, cacheKey, out, 0).Result()
			if err != nil {
				logger.M.Warnf("failed to cache warehouses: %v", err)
			}
		}
	}(warehouse)

	return warehouse, nil
}

func (b *backend) GetWarehouseFromCache(ctx context.Context, warehouseID int64) (GetWarehouseOutput, error) {
	var warehouse GetWarehouseOutput
	cacheKey := fmt.Sprintf("warehouses::%d", warehouseID)
	if exists, _ := b.clients.Cache.Exists(ctx, cacheKey).Result(); exists == 1 {
		out, _ := b.clients.Cache.Get(ctx, cacheKey).Result()
		sonic.Unmarshal([]byte(out), &warehouse)
	}
	return warehouse, nil
}
