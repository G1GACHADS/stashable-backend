package api

import (
	"errors"

	"github.com/G1GACHADS/backend/internal/backend"
	"github.com/gofiber/fiber/v2"
)

func (h *handler) DeleteWarehouse(c *fiber.Ctx) error {
	warehouseID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide a valid warehouse id",
		})
	}

	err = h.backend.DeleteWarehouse(c.Context(), int64(warehouseID))
	switch {
	case errors.Is(err, backend.ErrWarehouseDoesNotExists):
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
		"message": "Deleted warehouse",
	})
}
