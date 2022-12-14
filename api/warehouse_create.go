package api

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/G1GACHADS/stashable-backend/backend"
	"github.com/gofiber/fiber/v2"
)

type WarehouseCreateParamsRoom struct {
	ImageURL string  `json:"image_url"`
	Name     string  `json:"name"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Length   float64 `json:"length"`
	Price    float64 `json:"price"`
}

type WarehouseCreateParams struct {
	Name        string `json:"name"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`

	Province   string `json:"province"`
	City       string `json:"city"`
	StreetName string `json:"street_name"`
	ZipCode    int    `json:"zip_code"`

	Rooms       []WarehouseCreateParamsRoom `json:"rooms"`
	CategoryIDs []int64                     `json:"categories"`
}

func (p WarehouseCreateParams) Validate() error {
	if err := requiredFields(map[string]any{
		"name":         p.Name,
		"image_url":    p.ImageURL,
		"description":  p.Description,
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

func (h *handler) WarehouseCreate(c *fiber.Ctx) error {
	c.Accepts("application/json")

	var params WarehouseCreateParams

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

	err := h.backend.WarehouseCreate(c.Context(), backend.WarehouseCreateInput{
		Warehouse: backend.Warehouse{
			Name:        params.Name,
			ImageURL:    params.ImageURL,
			Description: params.Description,
			BasePrice:   rooms[0].Price,
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
