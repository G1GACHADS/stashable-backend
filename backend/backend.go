package backend

import (
	"context"

	"github.com/G1GACHADS/stashable-backend/clients"
	"github.com/G1GACHADS/stashable-backend/config"
)

type Backend interface {
	// Auth
	AuthenticateUser(ctx context.Context, email, password string) (AuthenticateUserOutput, error)
	RegisterUser(ctx context.Context, user User, address Address) (RegisterUserOutput, error)

	// Profile
	UserGetProfile(ctx context.Context, userID int64) (UserGetProfileOutput, error)
	UserGetProfileFromCache(ctx context.Context, userID int64) (UserGetProfileOutput, error)

	// Categories
	CategoryCreate(ctx context.Context, name string) (Category, error)
	CategoryDelete(ctx context.Context, categoryID int64) error

	// Warehouses
	WarehouseList(ctx context.Context, limit int) (WarehouseListOutput, error)
	WarehouseGet(ctx context.Context, warehouseID int64) (WarehouseGetOutput, error)
	WarehouseGetFromCache(ctx context.Context, warehouseID int64) (WarehouseGetOutput, error)
	WarehouseSearch(ctx context.Context, searchQuery string, limit int, priceAscending bool) (WarehouseSearchOutput, error)
	WarehouseCreate(ctx context.Context, input WarehouseCreateInput) error
	WarehouseDelete(ctx context.Context, warehouseID int64) error

	// Rentals
	UserGetRentals(ctx context.Context, userID int64) (UserGetRentalsOutput, error)
	RentalGet(ctx context.Context, userID, rentalID int64) (RentalGetOutput, error)
	RentalCreate(ctx context.Context, input RentalCreateInput) (int64, error)
	RentalUpdateStatus(ctx context.Context, rentalID, userID int64, status RentalStatus) error
}

type backend struct {
	clients *clients.Clients
	cfg     *config.Config
}

func New(clients *clients.Clients, cfg *config.Config) Backend {
	return &backend{clients: clients, cfg: cfg}
}
