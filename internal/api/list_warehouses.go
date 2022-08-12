package api

import (
	"github.com/G1GACHADS/backend/internal/backend"
	"github.com/gofiber/fiber/v2"
)

func (h *handler) ListWarehouses(c *fiber.Ctx) error {
	warehouses, _ := h.backend.ListWarehousesFromCache(c.Context())
	if len(warehouses) != 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"total_items": len(warehouses),
			"items":       warehouses,
		})
	}

	warehouses, err := h.backend.ListWarehouses(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side",
			"err":     err.Error(),
		})
	}

	if len(warehouses) == 0 {
		// render an empty array instead of null
		warehouses = []backend.ListWarehousesOutput{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total_items": len(warehouses),
		"items":       warehouses,
	})
}
