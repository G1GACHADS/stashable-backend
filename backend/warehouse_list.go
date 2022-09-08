package backend

import (
	"context"
)

type WarehouseListItem struct {
	Attributes    Warehouse `json:"attributes"`
	Relationships struct {
		Address    Address  `json:"address"`
		Categories []string `json:"categories"`
	} `json:"relationships"`
}

type WarehouseListOutput struct {
	TotalItems int                 `json:"total_items"`
	Items      []WarehouseListItem `json:"items"`
}

func (b *backend) WarehouseList(ctx context.Context, limit int) (WarehouseListOutput, error) {
	var warehouses []WarehouseListItem

	rows, err := b.clients.DB.Query(ctx, "SELECT * FROM warehouses_list LIMIT $1", limit)
	if err != nil {
		return WarehouseListOutput{}, err
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		var row WarehouseListItem
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
			&row.Attributes.RoomsCount,
			&row.Relationships.Address.ID,
			&row.Relationships.Address.Province,
			&row.Relationships.Address.City,
			&row.Relationships.Address.StreetName,
			&row.Relationships.Address.ZipCode,
			&row.Relationships.Categories,
		)
		if err != nil {
			return WarehouseListOutput{}, err
		}
		warehouses = append(warehouses, row)
	}

	return WarehouseListOutput{
		TotalItems: count,
		Items:      warehouses,
	}, nil
}
