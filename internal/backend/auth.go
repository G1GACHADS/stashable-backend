package backend

import (
	"context"
	"time"

	"github.com/G1GACHADS/stashable-backend/token/jwt"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticateUserOutput struct {
	Attributes struct {
		User        User   `json:"user"`
		AccessToken string `json:"access_token"`
	} `json:"attributes"`
	Relationships struct {
		Address Address `json:"address"`
	} `json:"relationships"`
}

func (b *backend) AuthenticateUser(ctx context.Context, email, password string) (AuthenticateUserOutput, error) {
	var out AuthenticateUserOutput

	err := b.clients.DB.
		QueryRow(ctx, "SELECT * FROM users LEFT JOIN addresses ON users.address_id = addresses.id WHERE users.email = $1", email).
		Scan(&out.Attributes.User.ID,
			&out.Attributes.User.AddressID,
			&out.Attributes.User.FullName,
			&out.Attributes.User.Email,
			&out.Attributes.User.PhoneNumber,
			&out.Attributes.User.Password,
			&out.Attributes.User.CreatedAt,
			&out.Relationships.Address.ID,
			&out.Relationships.Address.Province,
			&out.Relationships.Address.City,
			&out.Relationships.Address.StreetName,
			&out.Relationships.Address.ZipCode)
	if err != nil {
		return AuthenticateUserOutput{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(out.Attributes.User.Password), []byte(password)); err != nil {
		return AuthenticateUserOutput{}, err
	}

	accessToken, err := jwt.Generate(map[string]any{
		"exp":    time.Now().Add(b.cfg.App.JWTDuration).Unix(),
		"userID": out.Attributes.User.ID,
		"email":  out.Attributes.User.Email,
	}, b.cfg.App.JWTSecretKey)
	if err != nil {
		return AuthenticateUserOutput{}, err
	}

	out.Attributes.AccessToken = accessToken

	return out, nil
}

func (b *backend) RegisterUser(ctx context.Context, user User, address Address) (string, error) {
	tx, err := b.clients.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	var addressID int64
	err = tx.QueryRow(ctx, "INSERT INTO addresses (province, city, street_name, zip_code) VALUES ($1, $2, $3, $4) RETURNING id",
		address.Province, address.City, address.StreetName, address.ZipCode).
		Scan(&addressID)
	if err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	var alreadyExists bool
	err = tx.QueryRow(ctx, `
	SELECT EXISTS (
		SELECT 1 FROM users
		WHERE email = $1
		OR phone_number = $2
	)`, user.Email, user.PhoneNumber).Scan(&alreadyExists)
	if err != nil {
		return "", err
	}

	if alreadyExists {
		return "", ErrUserAlreadyExists
	}

	var userID int64
	err = tx.QueryRow(ctx, "INSERT INTO users (address_id, full_name, email, phone_number, password, created_at) VALUES ($1, $2, $3, $4, $5, now()) RETURNING id",
		addressID, user.FullName, user.Email, user.PhoneNumber, hash).Scan(&userID)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	accessToken, err := jwt.Generate(map[string]any{
		"userID": userID,
		"exp":    time.Now().Add(b.cfg.App.JWTDuration).Unix(),
		"email":  user.Email,
	}, b.cfg.App.JWTSecretKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
