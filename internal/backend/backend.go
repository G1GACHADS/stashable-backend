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

	// Warehouses
	ListWarehouses(ctx context.Context) ([]ListWarehousesOutput, error)
	ListWarehousesFromCache(ctx context.Context) ([]ListWarehousesOutput, error)
}

type backend struct {
	clients *clients.Clients
	cfg     *config.Config
}

func New(clients *clients.Clients, cfg *config.Config) Backend {
	return backend{clients: clients, cfg: cfg}
}
