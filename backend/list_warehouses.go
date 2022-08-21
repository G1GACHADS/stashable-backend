package backend

import (
	"context"
)

type ListWarehousesItem struct {
	Attributes    Warehouse `json:"attributes"`
	Relationships struct {
		Address    Address  `json:"address"`
		Categories []string `json:"categories"`
	} `json:"relationships"`
}

type ListWarehousesOutput struct {
	TotalItems int                  `json:"total_items"`
	Items      []ListWarehousesItem `json:"items"`
}

func (b *backend) ListWarehouses(ctx context.Context, limit int) (ListWarehousesOutput, error) {
	var warehouses []ListWarehousesItem

	rows, err := b.clients.DB.Query(ctx, "SELECT * FROM warehouses_list LIMIT $1", limit)
	if err != nil {
		return ListWarehousesOutput{}, err
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		var row ListWarehousesItem
		err := rows.Scan(
			&count,
			&row.Attributes.ID,
			&row.Attributes.AddressID,
			&row.Attributes.Name,
			&row.Attributes.ImageURL,
			&row.Attributes.Description,
			&row.Attributes.BasePrice,
			&row.Attributes.Email,
			&row.Attributes.PhoneNumber,
			&row.Attributes.CreatedAt,
			&row.Relationships.Address.ID,
			&row.Relationships.Address.Province,
			&row.Relationships.Address.City,
			&row.Relationships.Address.StreetName,
			&row.Relationships.Address.ZipCode,
			&row.Relationships.Categories,
		)
		if err != nil {
			return ListWarehousesOutput{}, err
		}
		warehouses = append(warehouses, row)
	}

	return ListWarehousesOutput{
		TotalItems: count,
		Items:      warehouses,
	}, nil
}
