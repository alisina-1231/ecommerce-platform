package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alisina-1231/ecommerce-platform/services/checkout/internal/models"
)

type ProductClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewProductClient(baseURL string) *ProductClient {
	return &ProductClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ProductClient) GetProduct(
	ctx context.Context,
	productID int,
) (*models.Product, error) {
	url := fmt.Sprintf(
		"%s/product/api/products/%d",
		c.baseURL,
		productID,
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
		return nil, fmt.Errorf("get product: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"product service returned status %d",
			resp.StatusCode,
		)
	}

	var product models.Product

	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("decode product: %w", err)
	}

	return &product, nil
}
