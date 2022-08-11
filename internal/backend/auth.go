package backend

import (
	"github.com/G1GACHADS/backend/internal/token/jwt"
	"golang.org/x/crypto/bcrypt"
)

func (b backend) AuthenticateUser(email, password string) (string, error) {
	var user User

	err := b.clients.DB.QueryRow("SELECT id, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	accessToken, err := jwt.Generate(map[string]any{
		"userID": user.ID,
		"exp":    b.cfg.App.JWTDuration,
		"email":  user.Email,
	}, b.cfg.App.JWTSecretKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (b backend) RegisterUser(user User, address Address) (string, error) {
	tx, err := b.clients.DB.BeginTx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var addressID int64
	err = tx.QueryRow("INSERT INTO addresses (province, city, street_name, zip_code) VALUES ($1, $2, $3, $4) RETURNING id",
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
	err = tx.QueryRow(`
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
	err = tx.QueryRow("INSERT INTO users (address_id, full_name, email, phone_number, password, created_at) VALUES ($1, $2, $3, $4, $5, now()) RETURNING id",
		addressID, user.FullName, user.Email, user.PhoneNumber, hash).Scan(&userID)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	accessToken, err := jwt.Generate(map[string]any{
		"userID": userID,
		"exp":    b.cfg.App.JWTDuration,
		"email":  user.Email,
	}, b.cfg.App.JWTSecretKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}