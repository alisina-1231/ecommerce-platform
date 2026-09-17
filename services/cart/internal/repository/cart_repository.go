package repository

import (
	"context"
	"fmt"

	"github.com/alisina-1231/ecommerce-platform/services/cart/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type CartRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewCartRepository(
	client *dynamodb.Client,
	tableName string,
) *CartRepository {
	return &CartRepository{
		client:    client,
		tableName: tableName,
	}
}

// GetCart retrieves a cart by user ID.
func (r *CartRepository) GetCart(
	ctx context.Context,
	userID string,
) (*models.Cart, error) {

	result, err := r.client.GetItem(
		ctx,
		&dynamodb.GetItemInput{
			TableName: aws.String(r.tableName),
			Key: map[string]types.AttributeValue{
				"user_id": &types.AttributeValueMemberS{
					Value: userID,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	if len(result.Item) == 0 {
		return &models.Cart{
			UserID: userID,
			Items:  []models.CartItem{},
		}, nil
	}

	var cart models.Cart

	if err := attributevalue.UnmarshalMap(
		result.Item,
		&cart,
	); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart: %w", err)
	}

	return &cart, nil
}

// SaveCart saves the entire cart.
func (r *CartRepository) SaveCart(
	ctx context.Context,
	cart *models.Cart,
) error {

	item, err := attributevalue.MarshalMap(cart)

	if err != nil {
		return fmt.Errorf("failed to marshal cart: %w", err)
	}

	_, err = r.client.PutItem(
		ctx,
		&dynamodb.PutItemInput{
			TableName: aws.String(r.tableName),
			Item:      item,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to save cart: %w", err)
	}

	return nil
}

// AddItem adds or increases the quantity of a product.
func (r *CartRepository) AddItem(
	ctx context.Context,
	userID string,
	productID int,
	quantity int,
) (*models.Cart, error) {

	cart, err := r.GetCart(ctx, userID)

	if err != nil {
		return nil, err
	}

	found := false

	for i := range cart.Items {

		if cart.Items[i].ProductID == productID {

			cart.Items[i].Quantity += quantity
			found = true

			break
		}
	}

	if !found {

		cart.Items = append(
			cart.Items,
			models.CartItem{
				ProductID: productID,
				Quantity:  quantity,
			},
		)

	}

	if err := r.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// UpdateItem updates the quantity of a product.
func (r *CartRepository) UpdateItem(
	ctx context.Context,
	userID string,
	productID int,
	quantity int,
) (*models.Cart, error) {

	cart, err := r.GetCart(ctx, userID)

	if err != nil {
		return nil, err
	}

	found := false

	for i := range cart.Items {

		if cart.Items[i].ProductID == productID {

			cart.Items[i].Quantity = quantity
			found = true

			break
		}
	}

	if !found {
		return nil, fmt.Errorf("product not found in cart")
	}

	if err := r.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// RemoveItem removes a product from the cart.
func (r *CartRepository) RemoveItem(
	ctx context.Context,
	userID string,
	productID int,
) (*models.Cart, error) {

	cart, err := r.GetCart(ctx, userID)

	if err != nil {
		return nil, err
	}

	found := false

	items := make([]models.CartItem, 0, len(cart.Items))

	for _, item := range cart.Items {

		if item.ProductID == productID {
			found = true
			continue
		}

		items = append(items, item)
	}

	if !found {
		return nil, fmt.Errorf("product not found in cart")
	}

	cart.Items = items

	if err := r.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// ClearCart removes all items from the cart.
func (r *CartRepository) ClearCart(
	ctx context.Context,
	userID string,
) error {

	_, err := r.client.DeleteItem(
		ctx,
		&dynamodb.DeleteItemInput{
			TableName: aws.String(r.tableName),
			Key: map[string]types.AttributeValue{
				"user_id": &types.AttributeValueMemberS{
					Value: userID,
				},
			},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}
