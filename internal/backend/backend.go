package backend

import (
	"github.com/G1GACHADS/backend/internal/clients"
	"github.com/G1GACHADS/backend/internal/config"
)

type Backend interface {
	AuthenticateUser(email, password string) (string, error)
	RegisterUser(user User, address Address) (string, error)
}

type backend struct {
	clients *clients.Clients
	cfg     *config.Config
}

func New(clients *clients.Clients, cfg *config.Config) Backend {
	return backend{clients: clients, cfg: cfg}
}
