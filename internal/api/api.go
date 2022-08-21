package api

import (
	"strings"
	"time"

	"github.com/G1GACHADS/backend/internal/api/middleware"
	"github.com/G1GACHADS/backend/internal/backend"
	"github.com/G1GACHADS/backend/internal/config"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/helmet/v2"
)

func NewServer(b backend.Backend, cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "stashable_http_server",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
	})

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{fiber.MethodGet, fiber.MethodPost, fiber.MethodHead}, ","),
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

	app.Static("/public", "./public")

	app.Get("/metrics", monitor.New(monitor.Config{Title: "Stashable metrics"}))

	app.Get("/health-check", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "UP",
		})
	})

	h := NewHandler(b, &cfg.App)

	// Auth routes
	app.Post("/auth/login", h.AuthenticateUser)
	app.Post("/auth/register", h.RegisterUser)

	// profile
	app.Get("/profile", middleware.Authenticated, h.GetUserProfile)

	// Categories
	app.Post("/categories", middleware.Authenticated, h.CreateCategory)
	app.Delete("/categories/:id", middleware.Authenticated, h.DeleteCategory)

	// Warehouse routes
	app.Get("/warehouses", h.ListWarehouses)
	app.Get("/warehouses/search", h.SearchWarehouses)
	app.Get("/warehouses/:id", h.GetWarehouse)
	app.Post("/warehouses", middleware.Authenticated, h.CreateWarehouse)
	app.Delete("/warehouses/:id", middleware.Authenticated, h.DeleteWarehouse)

	// Rentals
	app.Get("/rent/history", middleware.Authenticated, h.GetUserRentals)
	app.Post("/rent/:warehouseID", middleware.Authenticated, h.CreateRental)

	app.Patch("/rent/:id/pay",
		middleware.Authenticated,
		h.CreateUpdateRentalStatusHandler(backend.RentalStatusPaid))

	app.Patch("/rent/:id/cancel",
		middleware.Authenticated,
		h.CreateUpdateRentalStatusHandler(backend.RentalStatusCancelled))

	app.Patch("/rent/:id/return",
		middleware.Authenticated,
		h.CreateUpdateRentalStatusHandler(backend.RentalStatusReturned))

	return app
}

type handler struct {
	backend backend.Backend
	appCfg  *config.AppConfig
}

func NewHandler(backend backend.Backend, appCfg *config.AppConfig) *handler {
	return &handler{backend: backend, appCfg: appCfg}
}
