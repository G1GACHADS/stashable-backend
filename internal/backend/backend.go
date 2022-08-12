package backend

import (
	"context"

	"github.com/G1GACHADS/backend/internal/clients"
	"github.com/G1GACHADS/backend/internal/config"
)

type Backend interface {
	// Auth
	AuthenticateUser(ctx context.Context, email, password string) (string, error)
	RegisterUser(ctx context.Context, user User, address Address) (string, error)

	// Profile
	GetUserProfile(ctx context.Context, userID int64) (GetUserProfileOutput, error)
	GetUserProfileFromCache(ctx context.Context, userID int64) (GetUserProfileOutput, error)

	// Categories
	CreateCategory(ctx context.Context, name string) (Category, error)
	DeleteCategory(ctx context.Context, categoryID int64) error

	// Warehouses
	ListWarehouses(ctx context.Context) ([]ListWarehousesOutput, error)
	ListWarehousesFromCache(ctx context.Context) ([]ListWarehousesOutput, error)
	GetWarehouse(ctx context.Context, warehouseID int64) (GetWarehouseOutput, error)
	GetWarehouseFromCache(ctx context.Context, warehouseID int64) (GetWarehouseOutput, error)
	SearchWarehouse(ctx context.Context, filter SearchWarehouseFilter, query string) ([]Warehouse, error)
	CreateWarehouse(ctx context.Context, input CreateWarehouseInput) error
	DeleteWarehouse(ctx context.Context, warehouseID int64) error

	// Addresses
	UpdateAddress(ctx context.Context, addressID int64) error
}

type backend struct {
	clients *clients.Clients
	cfg     *config.Config
}

func New(clients *clients.Clients, cfg *config.Config) Backend {
	return backend{clients: clients, cfg: cfg}
}
