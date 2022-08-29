package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/G1GACHADS/stashable-backend/core/logger"
	"github.com/bytedance/sonic"
)

type UserGetProfileOutput struct {
	Attributes    User `json:"attributes"`
	Relationships struct {
		Address `json:"address"`
	} `json:"relationships"`
}

func (b *backend) UserGetProfile(ctx context.Context, userID int64) (UserGetProfileOutput, error) {
	var out UserGetProfileOutput

	query := `
	SELECT
		users.id,
		users.address_id,
		users.full_name,
		users.email,
		users.phone_number,
		users.created_at,
		addresses.*
	FROM users
	LEFT JOIN addresses ON users.address_id = addresses.id
	WHERE users.id = $1`

	err := b.clients.DB.QueryRow(ctx, query, userID).Scan(
		&out.Attributes.ID,
		&out.Attributes.AddressID,
		&out.Attributes.FullName,
		&out.Attributes.Email,
		&out.Attributes.PhoneNumber,
		&out.Attributes.CreatedAt,
		&out.Relationships.Address.ID,
		&out.Relationships.Address.Province,
		&out.Relationships.Address.City,
		&out.Relationships.Address.StreetName,
		&out.Relationships.Address.ZipCode,
	)
	if err != nil {
		return UserGetProfileOutput{}, err
	}

	// Cache the profile for future use
	go func(profile UserGetProfileOutput) {
		cacheKey := fmt.Sprintf("profile::%d", out.Attributes.ID)
		if exists, _ := b.clients.Cache.Exists(ctx, cacheKey).Result(); exists != 1 {
			out, _ := sonic.Marshal(profile)
			_, err := b.clients.Cache.Set(ctx, fmt.Sprintf("profile::%d", profile.Attributes.ID), out, time.Hour).Result()
			if err != nil {
				logger.M.Warnf("failed to cache profile: %v", err)
			}
			logger.M.Debugf("cached profile for user %d", profile.Attributes.ID)
		}
	}(out)

	return out, nil
}

func (b *backend) UserGetProfileFromCache(ctx context.Context, userID int64) (UserGetProfileOutput, error) {
	var profile UserGetProfileOutput
	cacheKey := fmt.Sprintf("profile::%d", userID)
	if exists, _ := b.clients.Cache.Exists(ctx, cacheKey).Result(); exists == 1 {
		out, _ := b.clients.Cache.Get(ctx, cacheKey).Result()
		sonic.Unmarshal([]byte(out), &profile)
	}
	return profile, nil
}
