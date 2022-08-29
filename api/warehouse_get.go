package api

import (
	"errors"

	"github.com/G1GACHADS/stashable-backend/backend"
	"github.com/gofiber/fiber/v2"
)

func (h *handler) WarehouseGet(c *fiber.Ctx) error {
	warehouseID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide a valid warehouse id")
	}

	cachedWarehouse, _ := h.backend.WarehouseGetFromCache(c.Context(), int64(warehouseID))
	if cachedWarehouse.Attributes.ID != 0 {
		return c.Status(fiber.StatusOK).JSON(cachedWarehouse)
	}

	warehouse, err := h.backend.WarehouseGet(c.Context(), int64(warehouseID))
	switch {
	case errors.Is(err, backend.ErrWarehouseDoesNotExists):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case err != nil:
		return fiber.NewError(fiber.StatusInternalServerError, "There was a problem on our side")
	}

	return c.Status(fiber.StatusOK).JSON(warehouse)
}
