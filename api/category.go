package api

import (
	"errors"

	"github.com/G1GACHADS/stashable-backend/backend"
	"github.com/gofiber/fiber/v2"
)

type CategoryCreateParams struct {
	Name string `json:"name"`
}

func (p CategoryCreateParams) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

func (h *handler) CategoryCreate(c *fiber.Ctx) error {
	c.Accepts("application/json")

	var params CategoryCreateParams

	if err := c.BodyParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := params.Validate(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	createdCategory, err := h.backend.CategoryCreate(c.Context(), params.Name)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "There was a problem on our side")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Created category",
		"data":    createdCategory,
	})
}

func (h *handler) CategoryDelete(c *fiber.Ctx) error {
	categoryID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide a valid category id")
	}

	err = h.backend.CategoryDelete(c.Context(), int64(categoryID))
	switch {
	case errors.Is(err, backend.ErrCategoryDoesNotExists):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case err != nil:
		return fiber.NewError(fiber.StatusInternalServerError, "There was a problem on our side")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Deleted category",
	})
}
