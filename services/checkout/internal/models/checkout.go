package models

type CheckoutRequest struct {
	UserID          string `json:"user_id"`
	ShippingAddress string `json:"shipping_address"`
	PaymentMethod   string `json:"payment_method"`
}

type CheckoutResponse struct {
	OrderID     int64   `json:"order_id"`
	PaymentID   int64   `json:"payment_id"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_amount"`
	Currency    string  `json:"currency"`
}

type Cart struct {
	UserID string     `json:"user_id"`
	Items  []CartItem `json:"items"`
}

type CartItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Product struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Stock    int     `json:"stock"`
}

type CreateOrderRequest struct {
	UserID          string  `json:"user_id"`
	TotalAmount     float64 `json:"total_amount"`
	Currency        string  `json:"currency"`
	ShippingAddress string  `json:"shipping_address"`
}

type Order struct {
	ID          int64   `json:"id"`
	UserID      string  `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`
}

type CreatePaymentRequest struct {
	OrderID       int64   `json:"order_id"`
	UserID        string  `json:"user_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	PaymentMethod string  `json:"payment_method"`
}

type Payment struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}
