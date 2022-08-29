package backend

import (
	"context"
	"fmt"

	"github.com/G1GACHADS/stashable-backend/core/logger"
	"github.com/jackc/pgx/v4"
)

func (b *backend) RentalUpdateStatus(ctx context.Context, rentalID, userID int64, status RentalStatus) error {
	var dbUserID int64
	err := b.clients.DB.QueryRow(ctx, "SELECT user_id FROM rentals WHERE id = $1", rentalID).Scan(&dbUserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrRentalDoesNotExists
		}
		return err
	}

	if dbUserID != userID {
		return ErrRentalDoesNotBelongToUser
	}

	_, err = b.clients.DB.Exec(ctx, "UPDATE rentals SET status = $1 WHERE id = $2", status, rentalID)
	if err != nil {
		return err
	}

	// Clear the cache
	go func(ctx context.Context) {
		cacheKey := fmt.Sprintf("user::rentals::%d", userID)
		_, err := b.clients.Cache.Del(ctx, cacheKey).Result()
		if err != nil {
			logger.M.Warnf("Couldn't refresh warehouses cache: %v", err)
		}
	}(ctx)

	return nil
}
