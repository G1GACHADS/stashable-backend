package api

import (
	"github.com/G1GACHADS/backend/internal/backend"
	"github.com/G1GACHADS/backend/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewServer(backend backend.Backend, cfg *config.Config) *fiber.App {
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/health-check", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "UN",
		})
	})

	return app
}

type handler struct {
	backend backend.Backend
}

func NewHandler(backend backend.Backend) *handler {
	return &handler{backend: backend}
}
