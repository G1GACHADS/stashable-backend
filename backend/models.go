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
	RentalSelfServiceShipping RentalShippingType = "self-service"
	RentalDeliveryShipping    RentalShippingType = "delivery"
)

type RentalStatus string

const (
	RentalStatusUnpaid    RentalStatus = "unpaid"
	RentalStatusPaid      RentalStatus = "paid"
	RentalStatusCancelled RentalStatus = "cancelled"
	RentalStatusReturned  RentalStatus = "returned"
)

type Rental struct {
	ID           int64              `json:"id"`
	UserID       int64              `json:"user_id"`
	WarehouseID  int64              `json:"warehouse_id"`
	RoomID       int64              `json:"room_id"`
	CategoryID   int64              `json:"category_id"`
	ImageURLs    []string           `json:"image_urls"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Weight       float64            `json:"weight"`
	Width        float64            `json:"width"`
	Height       float64            `json:"height"`
	Length       float64            `json:"length"`
	Quantity     int                `json:"quantity"`
	PaidAnnually bool               `json:"paid_annually"`
	Type         RentalShippingType `json:"type"`
	Status       RentalStatus       `json:"status"`
	CreatedAt    time.Time          `json:"created_at"`
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
