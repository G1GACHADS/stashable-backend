package backend

import (
	"context"
)

type GetUserProfileOutput struct {
	User User
	Address Address
}

func (b backend) GetUserProfile(ctx context.Context, userID int64) (GetUserProfileOutput, error) {
	var user User
	var userAddress Address

	query := `
	SELECT
		users.id,
		users.full_name,
		users.email,
		users.phone_number,
		users.created_at,
		addresses.id,
		addresses.province,
		addresses.city,
		addresses.street_name,
		addresses.zip_code
	FROM users
	LEFT JOIN addresses ON users.address_id = addresses.id
	WHERE users.id = $1`

	err := b.clients.DB.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PhoneNumber,
		&user.CreatedAt,
		&userAddress.ID,
		&userAddress.Province,
		&userAddress.City,
		&userAddress.StreetName,
		&userAddress.ZipCode,
	)
	if err != nil {
		return GetUserProfileOutput{}, err
	}

	return GetUserProfileOutput{User: user, Address: userAddress}, nil
}
