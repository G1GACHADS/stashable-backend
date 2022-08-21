package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) SearchWarehouses(c *fiber.Ctx) error {
	searchQuery := c.Query("q")
	priceAsc := c.Query("order_by") == "asc"
	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide a valid limit",
		})
	}

	searchResult, err := h.backend.SearchWarehouses(c.Context(), searchQuery, limit, priceAsc)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side",
			"err":     err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(searchResult)
}
