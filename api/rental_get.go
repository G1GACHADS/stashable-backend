package api

import (
	"github.com/gofiber/fiber/v2"
)

func (h *handler) RentalGet(c *fiber.Ctx) error {
	userID := int64(c.Locals("userID").(float64))
	rentalID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide a valid rental id")
	}

	rental, err := h.backend.RentalGet(c.Context(), userID, int64(rentalID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "There was a problem on our side")
	}

	return c.Status(fiber.StatusOK).JSON(rental)
}
