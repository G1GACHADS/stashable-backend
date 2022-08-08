package backend

import "github.com/G1GACHADS/backend/internal/clients"

type Backend interface{}

type backend struct {
	clients *clients.Clients
}

func New(clients *clients.Clients) Backend {
	return backend{clients: clients}
}
