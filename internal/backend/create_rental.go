package backend

import (
	"context"
)

type CreateRentalInput struct {
	UserID       int64
	WarehouseID  int64
	ImageURLs    []string
	Name         string
	Description  string
	Weight       float64
	Width        float64
	Height       float64
	Length       float64
	Quantity     int
	PaidAnnually bool
	Type         RentalType
}

func (b *backend) CreateRental(ctx context.Context, input CreateRentalInput) (int64, error) {
	var exists bool
	err := b.clients.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM warehouses WHERE id = $1)", input.WarehouseID).Scan(&exists)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, ErrWarehouseDoesNotExists
	}

	query := `
	INSERT INTO rentals (
		user_id,
		warehouse_id,
		image_urls,
		name,
		description,
		weight,
		length,
		width,
		height,
		quantity,
		paid_annually,
		type,
		status,
		created_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, now())
	RETURNING id`

	args := []any{
		input.UserID,
		input.WarehouseID,
		input.ImageURLs,
		input.Name,
		input.Description,
		input.Weight,
		input.Length,
		input.Width,
		input.Height,
		input.Quantity,
		input.PaidAnnually,
		input.Type,
		RentalStatusUnpaid,
	}

	var rentalID int64
	if err := b.clients.DB.QueryRow(ctx, query, args...).Scan(&rentalID); err != nil {
		return 0, err
	}

	return rentalID, nil
}
