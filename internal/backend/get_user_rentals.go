package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/G1GACHADS/backend/internal/logger"
	"github.com/bytedance/sonic"
)

type GetUserRentalsItemAttributes struct {
	Rental
	PaymentDue time.Time `json:"payment_due"`
}

type GetUserRentalsItem struct {
	Attributes    GetUserRentalsItemAttributes `json:"attributes"`
	Relationships struct {
		User      User      `json:"user"`
		Warehouse Warehouse `json:"warehouse"`
	} `json:"relationships"`
}

type GetUserRentalsOutput struct {
	TotalItems int                  `json:"total_items"`
	Items      []GetUserRentalsItem `json:"items"`
}

func (b *backend) GetUserRentals(ctx context.Context, userID int64) (GetUserRentalsOutput, error) {
	query := `
	SELECT
		rentals.id,
		rentals.user_id,
		rentals.warehouse_id,
		rentals.image_urls,
		rentals.name,
		rentals.description,
		rentals.weight,
		rentals.width,
		rentals.height,
		rentals.length,
		rentals.quantity,
		rentals.paid_annually,
		rentals.type,
		rentals.status,
		rentals.created_at,
		users.id,
		users.address_id,
		users.full_name,
		users.email,
		users.phone_number,
		users.created_at,
		warehouses.id,
		warehouses.address_id,
		warehouses.name,
		warehouses.image_url,
		warehouses.description,
		warehouses.description,
		warehouses.base_price,
		warehouses.created_at
	FROM rentals
	LEFT JOIN users ON rentals.user_id = users.id
	LEFT JOIN warehouses on rentals.warehouse_id = warehouses.id
	WHERE rentals.user_id = $1
	`

	var rentals []GetUserRentalsItem

	rows, err := b.clients.DB.Query(ctx, query, userID)
	if err != nil {
		return GetUserRentalsOutput{}, nil
	}

	for rows.Next() {
		var rental GetUserRentalsItem
		rows.Scan(
			&rental.Attributes.ID,
			&rental.Attributes.UserID,
			&rental.Attributes.WarehouseID,
			&rental.Attributes.ImageURLs,
			&rental.Attributes.Name,
			&rental.Attributes.Description,
			&rental.Attributes.Weight,
			&rental.Attributes.Width,
			&rental.Attributes.Height,
			&rental.Attributes.Length,
			&rental.Attributes.Quantity,
			&rental.Attributes.PaidAnnually,
			&rental.Attributes.Type,
			&rental.Attributes.Status,
			&rental.Attributes.CreatedAt,
			&rental.Relationships.User.ID,
			&rental.Relationships.User.AddressID,
			&rental.Relationships.User.FullName,
			&rental.Relationships.User.Email,
			&rental.Relationships.User.PhoneNumber,
			&rental.Relationships.User.CreatedAt,
			&rental.Relationships.Warehouse.ID,
			&rental.Relationships.Warehouse.AddressID,
			&rental.Relationships.Warehouse.Name,
			&rental.Relationships.Warehouse.ImageURL,
			&rental.Relationships.Warehouse.Description,
			&rental.Relationships.Warehouse.Description,
			&rental.Relationships.Warehouse.BasePrice,
			&rental.Relationships.Warehouse.CreatedAt,
		)
		if err != nil {
			return GetUserRentalsOutput{}, err
		}

		if rental.Attributes.PaidAnnually {
			rental.Attributes.PaymentDue =
				rental.Attributes.CreatedAt.Add(365 * 24 * time.Hour)
		} else {
			rental.Attributes.PaymentDue =
				rental.Attributes.CreatedAt.Add(30 * 24 * time.Hour)
		}

		rentals = append(rentals, rental)
	}

	out := GetUserRentalsOutput{
		TotalItems: len(rentals),
		Items:      rentals,
	}

	go func(rentals GetUserRentalsOutput) {
		cacheKey := fmt.Sprintf("user::rentals::%d", userID)
		if exists, _ := b.clients.Cache.Exists(ctx, cacheKey).Result(); exists != 1 {
			// Cache the warehouses for future use
			out, _ := sonic.Marshal(rentals)
			_, err := b.clients.Cache.Set(ctx, cacheKey, out, time.Hour).Result()
			if err != nil {
				logger.M.Warnf("failed to cache warehouses: %v", err)
			}
		}
	}(out)

	return out, nil
}

func (b *backend) GetUserRentalsFromCache(ctx context.Context, userID int64) (GetUserRentalsOutput, error) {
	var rentals GetUserRentalsOutput
	cacheKey := fmt.Sprintf("user::rentals::%d", userID)
	if exists, _ := b.clients.Cache.Exists(ctx, cacheKey).Result(); exists == 1 {
		out, _ := b.clients.Cache.Get(ctx, cacheKey).Result()
		sonic.Unmarshal([]byte(out), &rentals)
	}
	return rentals, nil
}
