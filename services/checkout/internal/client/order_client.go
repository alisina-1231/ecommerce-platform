package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alisina-1231/ecommerce-platform/services/checkout/internal/models"
)

type OrderClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewOrderClient(baseURL string) *OrderClient {
	return &OrderClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *OrderClient) CreateOrder(
	ctx context.Context,
	request models.CreateOrderRequest,
) (*models.Order, error) {
	body, err := json.Marshal(request)

	if err != nil {
		return nil, fmt.Errorf("encode order: %w", err)
	}

	url := fmt.Sprintf(
		"%s/order/api/orders",
		c.baseURL,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf(
			"order service returned status %d",
			resp.StatusCode,
		)
	}

	var order models.Order

	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("decode order: %w", err)
	}

	return &order, nil
}
