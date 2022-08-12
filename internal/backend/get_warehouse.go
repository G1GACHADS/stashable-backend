package backend

import (
	"context"
	"errors"
	"fmt"

	"github.com/G1GACHADS/backend/internal/logger"
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

func (b backend) GetWarehouse(ctx context.Context, warehouseID int64) (GetWarehouseOutput, error) {
	query := `
	SELECT
		warehouses.id,
		warehouses.name,
		warehouses.image_url,
		warehouses.description,
		warehouses.base_price,
		warehouses.created_at,
		addresses.id,
		addresses.province,
		addresses.city,
		addresses.street_name,
		addresses.zip_code,
		array_agg(categories.name)
	FROM warehouses
	LEFT JOIN addresses ON warehouses.address_id = addresses.id
	LEFT JOIN warehouse_categories as wc ON warehouses.id = wc.warehouse_id
	LEFT JOIN categories ON wc.category_id = categories.id
	WHERE warehouses.id = $1
	GROUP BY
		warehouses.id,
		addresses.id
	`

	var warehouse GetWarehouseOutput

	err := b.clients.DB.QueryRow(ctx, query, warehouseID).Scan(
		&warehouse.Attributes.ID,
		&warehouse.Attributes.Name,
		&warehouse.Attributes.ImageURL,
		&warehouse.Attributes.Description,
		&warehouse.Attributes.BasePrice,
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

func (b backend) GetWarehouseFromCache(ctx context.Context, warehouseID int64) (GetWarehouseOutput, error) {
	var warehouse GetWarehouseOutput
	cacheKey := fmt.Sprintf("warehouses::%d", warehouseID)
	if exists, _ := b.clients.Cache.Exists(ctx, cacheKey).Result(); exists == 1 {
		out, _ := b.clients.Cache.Get(ctx, cacheKey).Result()
		sonic.Unmarshal([]byte(out), &warehouse)
	}
	return warehouse, nil
}
