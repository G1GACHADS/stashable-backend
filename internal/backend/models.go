package backend

import "time"

type User struct {
	ID          int64  `json:"id"`
	AddressID   int64  `json:"address_id"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string
	CreatedAt   time.Time `json:"created_at"`
}

type Warehouse struct {
	ID          int64     `json:"id"`
	AddressID   int64     `json:"address_id"`
	Name        string    `json:"name"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"base_price"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}

type Item struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	WarehouseID int64     `json:"warehouse_id"`
	CategoryID  int64     `json:"category_id"`
	Name        string    `json:"name"`
	Weight      float64   `json:"weight"`
	Width       float64   `json:"width"`
	Height      float64   `json:"height"`
	Length      float64   `json:"length"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Address struct {
	ID         int64  `json:"id"`
	Province   string `json:"province"`
	City       string `json:"city"`
	StreetName string `json:"street_name"`
	ZipCode    int    `json:"zip_code"`
}
