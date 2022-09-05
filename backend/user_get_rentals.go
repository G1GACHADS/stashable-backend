package backend

import (
	"context"
	"time"
)

type UserGetRentalsItemAttributes struct {
	Rental
	PaymentDue time.Time `json:"payment_due"`
}

type UserGetRentalsItem struct {
	Attributes    UserGetRentalsItemAttributes `json:"attributes"`
	Relationships struct {
		Warehouse  Warehouse `json:"warehouse"`
		Address    Address   `json:"address"`
		Categories []string  `json:"categories"`
	} `json:"relationships"`
}

type UserGetRentalsOutput struct {
	TotalItems int                  `json:"total_items"`
	Items      []UserGetRentalsItem `json:"items"`
}

func (b *backend) UserGetRentals(ctx context.Context, userID int64) (UserGetRentalsOutput, error) {
	query := `
	SELECT
		r.*,
		w.*
	FROM rentals r
	LEFT JOIN warehouses_list w ON w.w_id = r.warehouse_id
	WHERE r.user_id = $1`

	var rentals []UserGetRentalsItem

	rows, err := b.clients.DB.Query(ctx, query, userID)
	if err != nil {
		return UserGetRentalsOutput{}, nil
	}
	defer rows.Close()

	for rows.Next() {
		var rental UserGetRentalsItem
		err = rows.Scan(
			&rental.Attributes.ID,
			&rental.Attributes.UserID,
			&rental.Attributes.WarehouseID,
			&rental.Attributes.CategoryID,
			&rental.Attributes.ImageURLs,
			&rental.Attributes.Name,
			&rental.Attributes.Description,
			&rental.Attributes.Weight,
			&rental.Attributes.Width,
			&rental.Attributes.Height,
			&rental.Attributes.Length,
			&rental.Attributes.Quantity,
			&rental.Attributes.PaidAnnually,
			&rental.Attributes.Type,
			&rental.Attributes.Status,
			&rental.Attributes.CreatedAt,
			&rental.Attributes.RoomID,
			nil, // warehouse count
			&rental.Relationships.Warehouse.ID,
			&rental.Relationships.Warehouse.AddressID,
			&rental.Relationships.Warehouse.Name,
			&rental.Relationships.Warehouse.ImageURL,
			&rental.Relationships.Warehouse.Description,
			&rental.Relationships.Warehouse.BasePrice,
			&rental.Relationships.Warehouse.Email,
			&rental.Relationships.Warehouse.PhoneNumber,
			&rental.Relationships.Warehouse.CreatedAt,
			&rental.Relationships.Warehouse.RoomsCount,
			&rental.Relationships.Address.ID,
			&rental.Relationships.Address.Province,
			&rental.Relationships.Address.City,
			&rental.Relationships.Address.StreetName,
			&rental.Relationships.Address.ZipCode,
			&rental.Relationships.Categories,
		)
		if err != nil {
			return UserGetRentalsOutput{}, err
		}

		if rental.Attributes.PaidAnnually {
			rental.Attributes.PaymentDue =
				rental.Attributes.CreatedAt.Add(365 * 24 * time.Hour)
		} else {
			rental.Attributes.PaymentDue =
				rental.Attributes.CreatedAt.Add(30 * 24 * time.Hour)
		}

		rentals = append(rentals, rental)
	}

	out := UserGetRentalsOutput{
		TotalItems: len(rentals),
		Items:      rentals,
	}

	return out, nil
}
