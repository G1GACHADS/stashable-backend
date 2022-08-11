package api

import (
	"errors"
	"net/mail"

	"github.com/G1GACHADS/backend/internal/backend"
	"github.com/gofiber/fiber/v2"
)

type AuthenticateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p AuthenticateUserParams) Validate() error {
	if p.Email == "" {
		return errors.New("email is required")
	}

	if p.Password == "" {
		return errors.New("password is required")
	}

	if _, err := mail.ParseAddress(p.Email); err != nil {
		return errors.New("invalid email address")
	}

	return nil
}

func (h *handler) AuthenticateUser(c *fiber.Ctx) error {
	c.Accepts("application/json")

	var params AuthenticateUserParams

	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email/password",
		})
	}

	if err := params.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	accessToken, err := h.backend.AuthenticateUser(params.Email, params.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email/password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessToken,
	})
}

type RegisterUserParams struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`

	Province   string `json:"province"`
	City       string `json:"city"`
	StreetName string `json:"street_name"`
	ZipCode    int    `json:"zip_code"`
}

func (p RegisterUserParams) Validate() error {
	if p.FullName == "" {
		return errors.New("full_name is required")
	}

	if p.Email == "" {
		return errors.New("email is required")
	}

	if _, err := mail.ParseAddress(p.Email); err != nil {
		return errors.New("invalid email address")
	}

	if p.PhoneNumber == "" {
		return errors.New("phone_number is required")
	}

	if p.Password == "" {
		return errors.New("password is required")
	}

	if p.Province == "" {
		return errors.New("province is required")
	}

	if p.City == "" {
		return errors.New("city is required")
	}

	if p.StreetName == "" {
		return errors.New("street_name is required")
	}

	if p.ZipCode == 0 {
		return errors.New("zip_code is required")
	}

	return nil
}

func (h *handler) RegisterUser(c *fiber.Ctx) error {
	c.Accepts("application/json")

	var params RegisterUserParams

	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Please fill all required fields",
		})
	}

	accessToken, err := h.backend.RegisterUser(
		backend.User{
			FullName:    params.FullName,
			Email:       params.Email,
			PhoneNumber: params.PhoneNumber,
			Password:    params.Password,
		},
		backend.Address{
			Province:   params.Province,
			City:       params.City,
			StreetName: params.StreetName,
			ZipCode:    params.ZipCode,
		})

	switch {
	case err == backend.ErrUserAlreadyExists:
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "User already exists",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side",
			"err":     err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"access_token": accessToken,
	})
}
