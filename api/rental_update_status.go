package api

import (
	"errors"

	"github.com/G1GACHADS/stashable-backend/backend"
	"github.com/G1GACHADS/stashable-backend/core/logger"
	"github.com/gofiber/fiber/v2"
)

// CreateRentalUpdateStatus to create support for multiple handlers for different rental status updates
func (h *handler) CreateRentalUpdateStatus(status int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rentalID, err := c.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Please provide a valid rental id")
		}
		userID := int64(c.Locals("userID").(float64))

		err = h.backend.RentalUpdateStatus(c.Context(), int64(rentalID), userID, status)
		switch {
		case errors.Is(err, backend.ErrRentalDoesNotExists):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case errors.Is(err, backend.ErrRentalDoesNotBelongToUser):
			return fiber.NewError(fiber.StatusForbidden, err.Error())
		case err != nil:
			logger.M.Errorf("Couldn't update rental status: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "There was a problem on our side")
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Updated rental status",
			"status":  status,
		})
	}
}
