package api

import (
	"github.com/gofiber/fiber/v2"
)

func (h *handler) GetUserRentals(c *fiber.Ctx) error {
	userID := int64(c.Locals("userID").(float64))

	rentals, err := h.backend.GetUserRentals(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "There was a problem on our side")
	}

	return c.Status(fiber.StatusOK).JSON(rentals)
}
