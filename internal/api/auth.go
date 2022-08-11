package api

import (
	"github.com/G1GACHADS/backend/internal/backend"
	"github.com/gofiber/fiber/v2"
)

type AuthenticateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *handler) AuthenticateUser(c *fiber.Ctx) error {
	c.Accepts("application/json")

	var params AuthenticateUserParams

	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email/password",
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
