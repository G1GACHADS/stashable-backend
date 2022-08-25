package api

import (
	"github.com/gofiber/fiber/v2"
)

func (h *handler) GetUserProfile(c *fiber.Ctx) error {
	userID := int64(c.Locals("userID").(float64))

	cachedProfile, _ := h.backend.GetUserProfileFromCache(c.Context(), userID)
	if cachedProfile.Attributes.ID != 0 {
		return c.Status(fiber.StatusOK).JSON(cachedProfile)
	}

	profile, err := h.backend.GetUserProfile(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "There was a problem on our side")
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}
