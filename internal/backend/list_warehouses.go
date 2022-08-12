package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/G1GACHADS/backend/internal/logger"
	"github.com/bytedance/sonic"
)

type ListWarehousesItem struct {
	Attributes    Warehouse `json:"attributes"`
	Relationships struct {
		Address    Address  `json:"address"`
		Categories []string `json:"categories"`
	} `json:"relationships"`
}

type ListWarehousesOutput struct {
	TotalItems int                  `json:"total_items"`
	Items      []ListWarehousesItem `json:"items"`
}

func (b backend) ListWarehouses(ctx context.Context, limit int) (ListWarehousesOutput, error) {
	query := `
	SELECT
		(SELECT count(id) FROM warehouses),
		warehouses.id,
		warehouses.address_id,
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
	GROUP BY
		warehouses.id,
		addresses.id
	LIMIT $1
	`

	var warehouses []ListWarehousesItem

	rows, err := b.clients.DB.Query(ctx, query, limit)
	if err != nil {
		return ListWarehousesOutput{}, err
	}

	var count int

	for rows.Next() {
		var row ListWarehousesItem
		err := rows.Scan(
			&count,
			&row.Attributes.ID,
			&row.Attributes.AddressID,
			&row.Attributes.Name,
			&row.Attributes.ImageURL,
			&row.Attributes.Description,
			&row.Attributes.BasePrice,
			&row.Attributes.CreatedAt,
			&row.Relationships.Address.ID,
			&row.Relationships.Address.Province,
			&row.Relationships.Address.City,
			&row.Relationships.Address.StreetName,
			&row.Relationships.Address.ZipCode,
			&row.Relationships.Categories,
		)
		if err != nil {
			return ListWarehousesOutput{}, err
		}
		warehouses = append(warehouses, row)
	}

	out := ListWarehousesOutput{
		TotalItems: count,
		Items:      warehouses,
	}

	go func(warehouses ListWarehousesOutput) {
		cacheKey := fmt.Sprintf("warehouses.limit(%d)", limit)
		if exists, _ := b.clients.Cache.Exists(ctx, cacheKey).Result(); exists != 1 {
			// Cache the warehouses for future use
			out, _ := sonic.Marshal(warehouses)
			_, err := b.clients.Cache.Set(ctx, cacheKey, out, 5*time.Minute).Result()
			if err != nil {
				logger.M.Warnf("failed to cache warehouses: %v", err)
			}
		}
	}(out)

	return out, nil
}

func (b backend) ListWarehousesFromCache(ctx context.Context, limit int) (ListWarehousesOutput, error) {
	var warehouses ListWarehousesOutput
	cacheKey := fmt.Sprintf("warehouses.limit(%d)", limit)
	if exists, _ := b.clients.Cache.Exists(ctx, cacheKey).Result(); exists == 1 {
		out, _ := b.clients.Cache.Get(ctx, cacheKey).Result()
		sonic.Unmarshal([]byte(out), &warehouses)
	}
	return warehouses, nil
}
