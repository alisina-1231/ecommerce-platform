package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alisina-1231/ecommerce-platform/services/checkout/internal/models"
)

type CartClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCartClient(baseURL string) *CartClient {
	return &CartClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *CartClient) GetCart(
	ctx context.Context,
	userID string,
) (*models.Cart, error) {
	url := fmt.Sprintf(
		"%s/cart/api/cart?user_id=%s",
		c.baseURL,
		userID,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get cart: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"cart service returned status %d",
			resp.StatusCode,
		)
	}

	var cart models.Cart

	if err := json.NewDecoder(resp.Body).Decode(&cart); err != nil {
		return nil, fmt.Errorf("decode cart: %w", err)
	}

	return &cart, nil
}

func (c *CartClient) ClearCart(
	ctx context.Context,
	userID string,
) error {
	url := fmt.Sprintf(
		"%s/cart/api/cart?user_id=%s",
		c.baseURL,
		userID,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		url,
		nil,
	)

	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("clear cart: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent &&
		resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"cart service returned status %d",
			resp.StatusCode,
		)
	}

	return nil
}
