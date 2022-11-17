package api

import (
	"strings"
	"time"

	"github.com/G1GACHADS/stashable-backend/api/middleware"
	"github.com/G1GACHADS/stashable-backend/backend"
	"github.com/G1GACHADS/stashable-backend/config"
	"github.com/G1GACHADS/stashable-backend/core/logger"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	loggerMiddleware "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
)

func NewServer(b backend.Backend, cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "stashable_http_server",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
		ErrorHandler: CustomErrorHandler,
	})

	app.Use(loggerMiddleware.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPatch,
			fiber.MethodDelete,
			fiber.MethodOptions,
			fiber.MethodHead}, ","),
		AllowHeaders: strings.Join([]string{fiber.HeaderAuthorization, fiber.HeaderContentType}, ","),
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

	h := NewHandler(b, &cfg.App)

	// Auth routes
	app.Post("/auth/login", h.AuthenticateUser)
	app.Post("/auth/register", h.RegisterUser)

	// profile
	app.Get("/profile", middleware.Authenticated, h.UserGetProfile)

	// Categories
	app.Post("/categories", middleware.Authenticated, h.CategoryCreate)
	app.Delete("/categories/:id", middleware.Authenticated, h.CategoryDelete)

	// Warehouse routes
	app.Get("/warehouses", h.WarehouseList)
	app.Get("/warehouses/search", h.WarehouseSearch)
	app.Get("/warehouses/:id", h.WarehouseGet)
	app.Post("/warehouses", middleware.Authenticated, h.WarehouseCreate)
	app.Delete("/warehouses/:id", middleware.Authenticated, h.WarehouseDelete)
	app.Get("/warehouses/:id/room/:roomID", middleware.Authenticated, h.RoomGet)

	// Rentals
	app.Get("/rent/history", middleware.Authenticated, h.UserGetRentals)
	app.Get("/rent/history/:id", middleware.Authenticated, h.RentalGet)
	app.Post("/rent/:warehouseID", middleware.Authenticated, h.RentalCreate)

	app.Patch("/rent/:id/pay",
		middleware.Authenticated,
		h.CreateRentalUpdateStatus(backend.RentalStatusPaid.Int()))

	app.Patch("/rent/:id/unpay",
		middleware.Authenticated,
		h.CreateRentalUpdateStatus(backend.RentalStatusUnpaid.Int()))

	app.Patch("/rent/:id/cancel",
		middleware.Authenticated,
		h.CreateRentalUpdateStatus(backend.RentalStatusCancelled.Int()))

	app.Patch("/rent/:id/return",
		middleware.Authenticated,
		h.CreateRentalUpdateStatus(backend.RentalStatusReturned.Int()))

	return app
}

func CustomErrorHandler(c *fiber.Ctx, err error) error {
	statusCode := fiber.StatusInternalServerError

	e, ok := err.(*fiber.Error)
	if ok {
		statusCode = e.Code
	} else {
		logger.M.Error(err)
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"message": e.Message,
	})
}

type handler struct {
	backend backend.Backend
	appCfg  *config.AppConfig
}

func NewHandler(backend backend.Backend, appCfg *config.AppConfig) *handler {
	return &handler{backend: backend, appCfg: appCfg}
}
