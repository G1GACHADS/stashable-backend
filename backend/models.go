package backend

import (
	"errors"
	"time"
)

type User struct {
	ID          int64     `json:"id"`
	AddressID   int64     `json:"address_id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type Warehouse struct {
	ID          int64  `json:"id"`
	AddressID   int64  `json:"address_id"`
	Name        string `json:"name"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	// BasePrice is the first room's price
	BasePrice   float64   `json:"base_price"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	RoomsCount  int       `json:"rooms_count"`
	CreatedAt   time.Time `json:"created_at"`
}

var (
	ErrWarehouseDoesNotExists = errors.New("warehouse does not exists")
)

type Room struct {
	ID          int64   `json:"id"`
	WarehouseID int64   `json:"warehouse_id"`
	ImageURL    string  `json:"image_url"`
	Name        string  `json:"name"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	Length      float64 `json:"length"`
	Price       float64 `json:"price"`
}

type RentalShippingType string

const (
	RentalPickUpTruckShipping RentalShippingType = "pick-up-truck"
	RentalPickUpBoxShipping   RentalShippingType = "pick-up-box"
	RentalVanShipping         RentalShippingType = "van"
	RentalTruckShipping       RentalShippingType = "truck"
	RentalSelfServiceShipping RentalShippingType = "self-service"
)

func (rst RentalShippingType) Int() int {
	switch rst {
	case RentalPickUpTruckShipping:
		return 1
	case RentalPickUpBoxShipping:
		return 2
	case RentalVanShipping:
		return 3
	case RentalTruckShipping:
		return 4
	case RentalSelfServiceShipping:
		return 5
	default:
		return 0

	}
}

type RentalStatus string

const (
	RentalStatusUnpaid    RentalStatus = "unpaid"
	RentalStatusPaid      RentalStatus = "paid"
	RentalStatusCancelled RentalStatus = "cancelled"
	RentalStatusReturned  RentalStatus = "returned"
)

func (rs RentalStatus) Int() int {
	switch rs {
	case RentalStatusUnpaid:
		return 1
	case RentalStatusPaid:
		return 2
	case RentalStatusCancelled:
		return 3
	case RentalStatusReturned:
		return 4
	default:
		return 0
	}
}

type Rental struct {
	ID             int64              `json:"id"`
	UserID         int64              `json:"user_id"`
	WarehouseID    int64              `json:"warehouse_id"`
	RoomID         int64              `json:"room_id"`
	CategoryID     int64              `json:"category_id"`
	ShippingTypeID int                `json:"shipping_type_id"`
	StatusID       int                `json:"status_id"`
	ShippingType   RentalShippingType `json:"shipping_type"`
	Status         RentalStatus       `json:"status"`
	ImageURLs      []string           `json:"image_urls"`
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	Weight         float64            `json:"weight"`
	Width          float64            `json:"width"`
	Height         float64            `json:"height"`
	Length         float64            `json:"length"`
	Quantity       int                `json:"quantity"`
	PaidAnnually   bool               `json:"paid_annually"`
	CreatedAt      time.Time          `json:"created_at"`
}

var (
	ErrRentalDoesNotExists       = errors.New("rental does not exists")
	ErrRentalDoesNotBelongToUser = errors.New("rental does not belong to user")
)

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

var (
	ErrCategoryDoesNotExists = errors.New("category does not exists")
)

type Address struct {
	ID         int64  `json:"id"`
	Province   string `json:"province"`
	City       string `json:"city"`
	StreetName string `json:"street_name"`
	ZipCode    int    `json:"zip_code"`
}
