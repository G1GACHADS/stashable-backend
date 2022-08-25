package api

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/G1GACHADS/stashable-backend/backend"
	"github.com/gofiber/fiber/v2"
)

type CreateWarehouseParamsRoom struct {
	ImageURL string  `json:"image_url"`
	Name     string  `json:"name"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Length   float64 `json:"length"`
	Price    float64 `json:"price"`
}

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

	Rooms       []CreateWarehouseParamsRoom `json:"rooms"`
	CategoryIDs []int64                     `json:"categories"`
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

	if len(p.Rooms) == 0 {
		return errors.New("rooms is required")
	}

	if len(p.CategoryIDs) == 0 {
		return errors.New("categories is required")
	}

	for idx, room := range p.Rooms {
		if err := requiredFields(map[string]any{
			fmt.Sprintf("rooms[%d].name", idx):      room.Name,
			fmt.Sprintf("rooms[%d].image_url", idx): room.ImageURL,
			fmt.Sprintf("rooms[%d].width", idx):     room.Width,
			fmt.Sprintf("rooms[%d].height", idx):    room.Height,
			fmt.Sprintf("rooms[%d].length", idx):    room.Length,
			fmt.Sprintf("rooms[%d].price", idx):     room.Price,
		}); err != nil {
			return err
		}
	}

	if _, err := mail.ParseAddress(p.Email); err != nil {
		return errors.New("invalid email address")
	}

	return nil
}

func (h *handler) CreateWarehouse(c *fiber.Ctx) error {
	c.Accepts("application/json")

	var params CreateWarehouseParams

	if err := c.BodyParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := params.Validate(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	rooms := make([]backend.Room, len(params.Rooms))
	for idx, room := range params.Rooms {
		rooms[idx] = backend.Room{
			ImageURL: room.ImageURL,
			Name:     room.Name,
			Width:    room.Width,
			Height:   room.Height,
			Length:   room.Length,
			Price:    room.Price,
		}
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
		Rooms:       rooms,
		CategoryIDs: params.CategoryIDs,
	})

	switch {
	case errors.Is(err, backend.ErrCategoryDoesNotExists):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case err != nil:
		return fiber.NewError(fiber.StatusInternalServerError, "There was a problem on our side")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Successfully created warehouse",
	})
}
