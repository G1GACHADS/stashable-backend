package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/G1GACHADS/stashable-backend/logger"
	"github.com/bytedance/sonic"
)

type GetUserRentalsItemAttributes struct {
	Rental
	PaymentDue time.Time `json:"payment_due"`
}

type GetUserRentalsItem struct {
	Attributes    GetUserRentalsItemAttributes `json:"attributes"`
	Relationships struct {
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
		r.*,
		w.*
	FROM rentals r
	LEFT JOIN warehouses AS w ON r.warehouse_id = w.id
	WHERE r.user_id = $1`

	var rentals []GetUserRentalsItem

	rows, err := b.clients.DB.Query(ctx, query, userID)
	if err != nil {
		return GetUserRentalsOutput{}, nil
	}
	defer rows.Close()

	for rows.Next() {
		var rental GetUserRentalsItem
		err = rows.Scan(
			&rental.Attributes.ID,
			&rental.Attributes.UserID,
			&rental.Attributes.WarehouseID,
			&rental.Attributes.CategoryID,
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
			&rental.Relationships.Warehouse.ID,
			&rental.Relationships.Warehouse.AddressID,
			&rental.Relationships.Warehouse.Name,
			&rental.Relationships.Warehouse.ImageURL,
			&rental.Relationships.Warehouse.Description,
			&rental.Relationships.Warehouse.BasePrice,
			&rental.Relationships.Warehouse.Email,
			&rental.Relationships.Warehouse.PhoneNumber,
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
