package models

type CartItem struct {
	ProductID int `json:"product_id" dynamodbav:"product_id"`
	Quantity  int `json:"quantity" dynamodbav:"quantity"`
}

type Cart struct {
	UserID string     `json:"user_id" dynamodbav:"user_id"`
	Items  []CartItem `json:"items" dynamodbav:"items"`
}
