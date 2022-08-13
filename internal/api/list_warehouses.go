package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) ListWarehouses(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide a valid limit",
		})
	}

	warehouses, err := h.backend.ListWarehouses(c.Context(), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side",
			"err":     err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(warehouses)
}
