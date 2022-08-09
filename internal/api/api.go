package api

import (
	"time"

	"github.com/G1GACHADS/backend/internal/backend"
	"github.com/G1GACHADS/backend/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"
)

func NewServer(backend backend.Backend, cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "storage_system_http_server",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: fiber.HeaderAuthorization,
	}))

	app.Use(helmet.New(helmet.Config{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSPreloadEnabled:    true,
		HSTSMaxAge:            63072000,
		HSTSExcludeSubdomains: false,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}))

	app.Get("/health-check", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "UP",
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