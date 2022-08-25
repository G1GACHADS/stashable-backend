package api

import (
	"errors"

	"github.com/G1GACHADS/stashable-backend/backend"
	"github.com/gofiber/fiber/v2"
)

func (h *handler) DeleteWarehouse(c *fiber.Ctx) error {
	warehouseID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide a valid warehouse id")
	}

	err = h.backend.DeleteWarehouse(c.Context(), int64(warehouseID))
	switch {
	case errors.Is(err, backend.ErrWarehouseDoesNotExists):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case err != nil:
		return fiber.NewError(fiber.StatusInternalServerError, "There was a problem on our side")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Deleted warehouse",
	})
}
