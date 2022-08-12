package backend

import (
	"context"

	"github.com/G1GACHADS/backend/internal/logger"
	"github.com/bytedance/sonic"
)

type ListWarehousesOutput struct {
	Warehouse  Warehouse `json:"warehouse"`
	Address    Address   `json:"warehouse_address"`
	Categories []string  `json:"categories"`
}

func (b backend) ListWarehouses(ctx context.Context) ([]ListWarehousesOutput, error) {
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
	GROUP BY
		warehouses.id,
		addresses.id
	`

	var warehouses []ListWarehousesOutput

	rows, err := b.clients.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var row ListWarehousesOutput
		err := rows.Scan(
			&row.Warehouse.ID,
			&row.Warehouse.Name,
			&row.Warehouse.ImageURL,
			&row.Warehouse.Description,
			&row.Warehouse.BasePrice,
			&row.Warehouse.CreatedAt,
			&row.Address.ID,
			&row.Address.Province,
			&row.Address.City,
			&row.Address.StreetName,
			&row.Address.ZipCode,
			&row.Categories,
		)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, row)
	}

	go func(warehouses []ListWarehousesOutput) {
		if exists, _ := b.clients.Cache.Exists(ctx, "warehouses").Result(); exists != 1 {
			// Cache the warehouses for future use
			out, _ := sonic.Marshal(warehouses)
			_, err := b.clients.Cache.Set(ctx, "warehouses", out, 0).Result()
			if err != nil {
				logger.M.Warnf("failed to cache warehouses: %v", err)
			}
		}
	}(warehouses)

	return warehouses, nil
}

func (b backend) ListWarehousesFromCache(ctx context.Context) ([]ListWarehousesOutput, error) {
	var warehouses []ListWarehousesOutput
	if exists, _ := b.clients.Cache.Exists(ctx, "warehouses").Result(); exists == 1 {
		out, _ := b.clients.Cache.Get(ctx, "warehouses").Result()
		sonic.Unmarshal([]byte(out), &warehouses)
	}
	return warehouses, nil
}
