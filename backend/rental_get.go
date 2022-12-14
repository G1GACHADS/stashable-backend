package backend

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
)

type RentalGetAttributes struct {
	Rental
	PaymentDue time.Time `json:"payment_due"`
}

type RentalGetOutput struct {
	Attributes    RentalGetAttributes `json:"attributes"`
	Relationships struct {
		Warehouse  Warehouse `json:"warehouse"`
		Address    Address   `json:"address"`
		Categories []string  `json:"categories"`
		Room       Room      `json:"room"`
	} `json:"relationships"`
}

func (b *backend) RentalGet(ctx context.Context, userID, rentalID int64) (RentalGetOutput, error) {
	query := `
	SELECT
		r.*,
		st.name,
		rs.name,
		w.*,
		rm.*
	FROM rentals r
	LEFT JOIN shipping_types st ON st.id = r.shipping_type_id
	LEFT JOIN rental_statuses rs ON rs.id = r.rental_status_id
	LEFT JOIN warehouses_list w ON w.w_id = r.warehouse_id
	LEFT JOIN rooms rm ON rm.id = r.room_id
	WHERE r.id = $1`

	var rental RentalGetOutput

	err := b.clients.DB.QueryRow(ctx, query, rentalID).Scan(
		&rental.Attributes.ID,
		&rental.Attributes.UserID,
		&rental.Attributes.WarehouseID,
		&rental.Attributes.CategoryID,
		&rental.Attributes.ShippingTypeID,
		&rental.Attributes.StatusID,
		&rental.Attributes.ImageURLs,
		&rental.Attributes.Name,
		&rental.Attributes.Description,
		&rental.Attributes.Weight,
		&rental.Attributes.Width,
		&rental.Attributes.Height,
		&rental.Attributes.Length,
		&rental.Attributes.Quantity,
		&rental.Attributes.PaidAnnually,
		&rental.Attributes.CreatedAt,
		&rental.Attributes.RoomID,
		&rental.Attributes.ShippingType,
		&rental.Attributes.Status,
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
		&rental.Relationships.Room.ID,
		&rental.Relationships.Room.WarehouseID,
		&rental.Relationships.Room.ImageURL,
		&rental.Relationships.Room.Name,
		&rental.Relationships.Room.Width,
		&rental.Relationships.Room.Height,
		&rental.Relationships.Room.Length,
		&rental.Relationships.Room.Price,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RentalGetOutput{}, ErrRentalDoesNotExists
		}
		return RentalGetOutput{}, err
	}

	if rental.Attributes.UserID != userID {
		return RentalGetOutput{}, ErrRentalDoesNotExists
	}

	if rental.Attributes.PaidAnnually {
		rental.Attributes.PaymentDue =
			rental.Attributes.CreatedAt.Add(365 * 24 * time.Hour)
	} else {
		rental.Attributes.PaymentDue =
			rental.Attributes.CreatedAt.Add(30 * 24 * time.Hour)
	}

	return rental, nil
}
