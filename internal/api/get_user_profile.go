package api

import (
	"github.com/gofiber/fiber/v2"
)

func (h *handler) GetUserProfile(c *fiber.Ctx) error {
	userID := int64(c.Locals("userID").(float64))

	profile, err := h.backend.GetUserProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side",
			"err":     err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"profile":         profile.User,
		"profile_address": profile.Address,
	})
}
