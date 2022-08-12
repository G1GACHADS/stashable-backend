package backend

import "context"

type SearchWarehouseFilter uint

const (
	SearchWarehouseFilterByName = iota << 1
	SearchWarehouseFilterByCategory
	SearchWarehouseFilterByPrice
)

func (b backend) SearchWarehouse(ctx context.Context, filter SearchWarehouseFilter, query string) ([]Warehouse, error) {
	return nil, nil
}
