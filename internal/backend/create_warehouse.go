package backend

import "context"

type CreateWarehouseInput struct {
	Warehouse Warehouse
	Address   Address
}

func (b backend) CreateWarehouse(ctx context.Context, input CreateWarehouseInput) error {
	return nil
}
