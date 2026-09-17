package models

import "time"

type Payment struct {
	ID            int64     `json:"id"`
	OrderID       int64     `json:"order_id"`
	UserID        string    `json:"user_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	TransactionID *string   `json:"transaction_id,omitempty"`
	FailureReason *string   `json:"failure_reason,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreatePaymentRequest struct {
	OrderID       int64   `json:"order_id"`
	UserID        string  `json:"user_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	PaymentMethod string  `json:"payment_method"`
}

type UpdatePaymentStatusRequest struct {
	Status        string  `json:"status"`
	TransactionID *string `json:"transaction_id,omitempty"`
	FailureReason *string `json:"failure_reason,omitempty"`
}
