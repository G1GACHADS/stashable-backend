package backend

import (
	"context"
	"fmt"

	"github.com/G1GACHADS/stashable-backend/core/logger"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type WarehouseCreateInput struct {
	Warehouse   Warehouse
	Address     Address
	Rooms       []Room
	CategoryIDs []int64
}

func (b *backend) WarehouseCreate(ctx context.Context, input WarehouseCreateInput) error {
	tx, err := b.clients.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := b.checkIfCategoryIDsExists(ctx, tx, input.CategoryIDs); err != nil {
		return err
	}

	warehouseID, err := b.insertAddressAndWarehouse(ctx, tx, input)
	if err != nil {
		return err
	}

	err = b.bulkInsertWarehouseRooms(ctx, tx, warehouseID, input.Rooms)
	if err != nil {
		return err
	}

	err = b.bulkInsertWarehouseCategories(ctx, tx, warehouseID, input.CategoryIDs)
	if err != nil {
		return err
	}

	// Clear the cache
	go func(ctx context.Context) {
		_, err := b.clients.Cache.Del(ctx, "warehouses").Result()
		if err != nil {
			logger.M.Warnf("Couldn't refresh warehouses cache: %v", err)
		}
	}(ctx)

	return tx.Commit(ctx)
}

func (b *backend) checkIfCategoryIDsExists(ctx context.Context, tx pgx.Tx, ids []int64) error {
	var categoryCount int

	pgIds := &pgtype.Int8Array{}
	pgIds.Set(ids)

	// Check if all the CategoryIDs exists in the db
	err := tx.QueryRow(ctx,
		`SELECT count(id) FROM categories WHERE id = ANY($1)`, pgIds).
		Scan(&categoryCount)
	if err != nil {
		return err
	}

	if categoryCount != len(ids) {
		return ErrCategoryDoesNotExists
	}

	return nil
}

func (b *backend) insertAddressAndWarehouse(ctx context.Context, tx pgx.Tx, input WarehouseCreateInput) (int64, error) {
	var addressID int64
	err := tx.QueryRow(ctx, "INSERT INTO addresses (province, city, street_name, zip_code) VALUES ($1, $2, $3, $4) RETURNING id",
		input.Address.Province,
		input.Address.City,
		input.Address.StreetName,
		input.Address.ZipCode).Scan(&addressID)
	if err != nil {
		return 0, err
	}

	var warehouseID int64
	err = tx.QueryRow(ctx, `
	INSERT INTO warehouses
		(address_id, name, image_url, description, base_price, email, phone_number, created_at)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, now())
	RETURNING id
	`,
		addressID,
		input.Warehouse.Name,
		input.Warehouse.ImageURL,
		input.Warehouse.Description,
		input.Warehouse.BasePrice,
		input.Warehouse.Email,
		input.Warehouse.PhoneNumber).Scan(&warehouseID)
	if err != nil {
		return 0, err
	}

	return warehouseID, nil
}

func (b *backend) bulkInsertWarehouseRooms(ctx context.Context, tx pgx.Tx, warehouseID int64, rooms []Room) error {
	var inputRows [][]any
	for _, room := range rooms {
		inputRows = append(inputRows, []any{
			warehouseID,
			room.ImageURL,
			room.Name,
			room.Width,
			room.Height,
			room.Length,
			room.Price,
		})
	}

	copyCount, err := tx.CopyFrom(ctx,
		pgx.Identifier{"rooms"},
		[]string{
			"warehouse_id",
			"image_url",
			"name",
			"width",
			"height",
			"length",
			"price",
		},
		pgx.CopyFromRows(inputRows))
	if err != nil {
		return err
	}

	if int(copyCount) != len(rooms) {
		return fmt.Errorf("expected %d rows to be copied, but %d were copied", len(rooms), copyCount)
	}

	return nil
}

func (b *backend) bulkInsertWarehouseCategories(ctx context.Context, tx pgx.Tx, warehouseID int64, categoriesID []int64) error {
	var inputRows [][]any
	for _, categoryID := range categoriesID {
		inputRows = append(inputRows, []any{warehouseID, categoryID})
	}

	copyCount, err := tx.CopyFrom(ctx,
		pgx.Identifier{"warehouse_categories"},
		[]string{"warehouse_id", "category_id"},
		pgx.CopyFromRows(inputRows))
	if err != nil {
		return err
	}

	if int(copyCount) != len(categoriesID) {
		return fmt.Errorf("expected %d rows to be copied, but %d were copied", len(categoriesID), copyCount)
	}

	return nil
}
