package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) ListWarehouses(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide a valid limit")
	}

	warehouses, err := h.backend.ListWarehouses(c.Context(), limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "There was a problem on our side")
	}

	return c.Status(fiber.StatusOK).JSON(warehouses)
}
