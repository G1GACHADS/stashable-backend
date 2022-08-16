package api

import (
	"errors"

	"github.com/G1GACHADS/backend/internal/backend"
	"github.com/gofiber/fiber/v2"
)

// CreateUpdateRentalStatusHandler to create multiple handlers for different rental status updates
func (h *handler) CreateUpdateRentalStatusHandler(status backend.RentalStatus) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rentalID, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Please provide a valid rental id",
			})
		}
		userID := int64(c.Locals("userID").(float64))

		err = h.backend.UpdateRentalStatus(c.Context(), int64(rentalID), userID, status)
		switch {
		case errors.Is(err, backend.ErrRentalDoesNotExists):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": err.Error(),
			})
		case errors.Is(err, backend.ErrRentalDoesNotBelongToUser):
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": err.Error(),
			})
		case err != nil:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "There was a problem on our side",
				"err":     err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Updated rental status",
			"status":  status,
		})
	}
}
