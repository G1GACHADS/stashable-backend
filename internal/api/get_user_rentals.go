package api

import "github.com/gofiber/fiber/v2"

func (h *handler) GetUserRentals(c *fiber.Ctx) error {
	userID := int64(c.Locals("userID").(float64))

	cachedRentals, _ := h.backend.GetUserRentalsFromCache(c.Context(), userID)
	if cachedRentals.TotalItems != 0 {
		return c.Status(fiber.StatusOK).JSON(cachedRentals)
	}

	rentals, err := h.backend.GetUserRentals(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side",
			"err":     err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(rentals)
}
