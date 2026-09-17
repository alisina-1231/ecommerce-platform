package models

import "time"

type Order struct {
	ID              int64     `json:"id"`
	UserID          string    `json:"user_id"`
	TotalAmount     float64   `json:"total_amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	ShippingAddress string    `json:"shipping_address"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateOrderRequest struct {
	UserID          string  `json:"user_id"`
	TotalAmount     float64 `json:"total_amount"`
	Currency        string  `json:"currency"`
	ShippingAddress string  `json:"shipping_address"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}
