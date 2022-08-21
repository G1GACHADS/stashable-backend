package api

import (
	"errors"
	"net/mail"

	"github.com/G1GACHADS/stashable-backend/internal/backend"
	"github.com/gofiber/fiber/v2"
)

type CreateWarehouseParams struct {
	Name        string  `json:"name"`
	ImageURL    string  `json:"image_url"`
	Description string  `json:"description"`
	BasePrice   float64 `json:"base_price"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phone_number"`

	Province   string `json:"province"`
	City       string `json:"city"`
	StreetName string `json:"street_name"`
	ZipCode    int    `json:"zip_code"`

	CategoryIDs []int64 `json:"categories"`
}

func (p CreateWarehouseParams) Validate() error {
	if err := requiredFields(map[string]any{
		"name":         p.Name,
		"image_url":    p.ImageURL,
		"description":  p.Description,
		"base_price":   p.BasePrice,
		"email":        p.Email,
		"phone_number": p.PhoneNumber,
		"province":     p.Province,
		"city":         p.City,
		"street_name":  p.StreetName,
		"zip_code":     p.ZipCode,
	}); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(p.Email); err != nil {
		return errors.New("invalid email address")
	}

	if len(p.CategoryIDs) == 0 {
		return errors.New("categories is required")
	}

	return nil
}

func (h *handler) CreateWarehouse(c *fiber.Ctx) error {
	c.Accepts("application/json")

	var params CreateWarehouseParams

	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := params.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err := h.backend.CreateWarehouse(c.Context(), backend.CreateWarehouseInput{
		Warehouse: backend.Warehouse{
			Name:        params.Name,
			ImageURL:    params.ImageURL,
			Description: params.Description,
			BasePrice:   params.BasePrice,
			Email:       params.Email,
			PhoneNumber: params.PhoneNumber,
		},
		Address: backend.Address{
			Province:   params.Province,
			City:       params.City,
			StreetName: params.StreetName,
			ZipCode:    params.ZipCode,
		},
		CategoryIDs: params.CategoryIDs,
	})

	switch {
	case errors.Is(err, backend.ErrCategoryDoesNotExists):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Successfully created warehouse",
	})
}
