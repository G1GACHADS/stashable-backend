package backend

import (
	"context"
	"errors"
)

type CreateRentalInput struct {
	UserID       int64
	WarehouseID  int64
	CategoryID   int64
	RoomID       int64
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

var ErrWarehouseOrCategoryOrRoomDoesNotExists = errors.New("warehouse or category or room does not exists")

func (b *backend) CreateRental(ctx context.Context, input CreateRentalInput) (int64, error) {
	existsQuery := `
	SELECT EXISTS (SELECT 1 FROM warehouses WHERE id = $1) AND
		   EXISTS (SELECT 1 FROM categories WHERE id = $2) AND
		   EXISTS (SELECT 1 FROM rooms WHERE id = $3)`

	var exists bool
	err := b.clients.DB.QueryRow(ctx, existsQuery, input.WarehouseID, input.CategoryID, input.RoomID).
		Scan(&exists)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, ErrWarehouseOrCategoryOrRoomDoesNotExists
	}

	query := `
	INSERT INTO rentals (
		user_id,
		warehouse_id,
		category_id,
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
		created_at,
		room_id
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, now(), $15)
	RETURNING id`

	args := []any{
		input.UserID,
		input.WarehouseID,
		input.CategoryID,
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
		input.RoomID,
	}

	var rentalID int64
	if err := b.clients.DB.QueryRow(ctx, query, args...).Scan(&rentalID); err != nil {
		return 0, err
	}

	return rentalID, nil
}
