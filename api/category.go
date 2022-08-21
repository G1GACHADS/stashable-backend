package api

import (
	"errors"

	"github.com/G1GACHADS/stashable-backend/backend"
	"github.com/gofiber/fiber/v2"
)

type CreateCategoryParams struct {
	Name string `json:"name"`
}

func (p CreateCategoryParams) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

func (h *handler) CreateCategory(c *fiber.Ctx) error {
	c.Accepts("application/json")

	var params CreateCategoryParams

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

	createdCategory, err := h.backend.CreateCategory(c.Context(), params.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side",
			"err":     err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Created category",
		"data":    createdCategory,
	})
}

func (h *handler) DeleteCategory(c *fiber.Ctx) error {
	categoryID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide a valid category id",
		})
	}

	err = h.backend.DeleteCategory(c.Context(), int64(categoryID))
	switch {
	case errors.Is(err, backend.ErrCategoryDoesNotExists):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side",
			"err":     err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Deleted category",
	})
}
