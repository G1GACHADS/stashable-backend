package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) WarehouseSearch(c *fiber.Ctx) error {
	searchQuery := c.Query("q")
	priceAsc := c.Query("order_by") == "asc"
	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide a valid limit")
	}

	searchResult, err := h.backend.WarehouseSearch(c.Context(), searchQuery, limit, priceAsc)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "There was a problem on our side")
	}

	return c.Status(fiber.StatusOK).JSON(searchResult)
}
