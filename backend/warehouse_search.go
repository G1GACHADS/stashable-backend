package backend

import (
	"context"
)

type WarehouseSearchItem struct {
	Attributes    Warehouse `json:"attributes"`
	Relationships struct {
		Address    Address  `json:"address"`
		Categories []string `json:"categories"`
	} `json:"relationships"`
}

type WarehouseSearchOutput struct {
	TotalItems int                   `json:"total_items"`
	Items      []WarehouseSearchItem `json:"items"`
}

func (b *backend) WarehouseSearch(ctx context.Context, searchQuery string, limit int, priceAscending bool) (WarehouseSearchOutput, error) {
	sql := `SELECT * FROM warehouses_list
	WHERE
		name ILIKE '%' || $1 || '%' OR
		province ILIKE '%' || $1 || '%' OR
		city ILIKE '%' || $1 || '%' OR
		street_name ILIKE '%' || $1 || '%'
	LIMIT $2
	`

	rows, err := b.clients.DB.Query(ctx, sql, searchQuery, limit)
	if err != nil {
		return WarehouseSearchOutput{}, err
	}
	defer rows.Close()

	var warehouses []WarehouseSearchItem

	for rows.Next() {
		var row WarehouseSearchItem
		err := rows.Scan(
			nil,
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
			return WarehouseSearchOutput{}, err
		}
		warehouses = append(warehouses, row)
	}

	return WarehouseSearchOutput{
		TotalItems: len(warehouses),
		Items:      warehouses,
	}, nil
}
