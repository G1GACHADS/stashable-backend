package backend

import (
	"context"
	"errors"
	"fmt"

	"github.com/G1GACHADS/stashable-backend/logger"
	"github.com/jackc/pgx/v4"
)

func (b *backend) DeleteWarehouse(ctx context.Context, warehouseID int64) error {
	var addressID int64
	err := b.clients.DB.QueryRow(ctx, "SELECT address_id FROM warehouses WHERE id = $1", warehouseID).Scan(&addressID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrWarehouseDoesNotExists
		}
		return err
	}

	// We should just only delete the address since
	// it will cascade delete the warehouse
	_, err = b.clients.DB.Exec(ctx, "DELETE FROM addresses WHERE id = $1", addressID)
	if err != nil {
		return err
	}

	// Clear the cache
	go func(ctx context.Context) {
		_, err := b.clients.Cache.Del(ctx, "warehouses").Result()
		if err != nil {
			logger.M.Warnf("Couldn't refresh warehouses cache: %v", err)
		}

		_, err = b.clients.Cache.Del(ctx, fmt.Sprintf("warehouses::%d", warehouseID)).Result()
		if err != nil {
			logger.M.Warnf("Couldn't refresh warehouses cache: %v", err)
		}
	}(ctx)

	return nil
}
