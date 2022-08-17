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
	ListWarehouses(ctx context.Context, limit int) (ListWarehousesOutput, error)
	GetWarehouse(ctx context.Context, warehouseID int64) (GetWarehouseOutput, error)
	GetWarehouseFromCache(ctx context.Context, warehouseID int64) (GetWarehouseOutput, error)
	SearchWarehouses(ctx context.Context, searchQuery string, limit int, priceAscending bool) (SearchWarehousesOutput, error)
	CreateWarehouse(ctx context.Context, input CreateWarehouseInput) error
	DeleteWarehouse(ctx context.Context, warehouseID int64) error

	// Rentals
	GetUserRentals(ctx context.Context, userID int64) (GetUserRentalsOutput, error)
	GetUserRentalsFromCache(ctx context.Context, userID int64) (GetUserRentalsOutput, error)
	CreateRental(ctx context.Context, input CreateRentalInput) (int64, error)
	UpdateRentalStatus(ctx context.Context, rentalID, userID int64, status RentalStatus) error

	// Addresses
	UpdateAddress(ctx context.Context, addressID int64) error
}

type backend struct {
	clients *clients.Clients
	cfg     *config.Config
}

func New(clients *clients.Clients, cfg *config.Config) Backend {
	return &backend{clients: clients, cfg: cfg}
}
