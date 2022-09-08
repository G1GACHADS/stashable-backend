package api

import (
	"github.com/G1GACHADS/stashable-backend/backend"
	"github.com/gofiber/fiber/v2"
)

func (h *handler) RoomGet(c *fiber.Ctx) error {
	type roomGetRelationshipsResponse struct {
		Warehouse  backend.Warehouse `json:"warehouse"`
		Address    backend.Address   `json:"address"`
		Categories []string          `json:"categories"`
	}

	type roomGetResponse struct {
		Attributes    backend.Room                 `json:"attributes"`
		Relationships roomGetRelationshipsResponse `json:"relationships"`
	}

	roomID, err := c.ParamsInt("roomID")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide a valid rental id")
	}

	room, err := h.backend.RoomGet(c.Context(), int64(roomID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "There was a problem on our side")
	}

	warehouse, err := h.backend.WarehouseGet(c.Context(), room.Attributes.WarehouseID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "There was a problem on our side")
	}

	return c.Status(fiber.StatusOK).JSON(roomGetResponse{
		Attributes: room.Attributes,
		Relationships: roomGetRelationshipsResponse{
			Warehouse:  warehouse.Attributes,
			Address:    warehouse.Relationships.Address,
			Categories: warehouse.Relationships.Categories,
		},
	})
}
