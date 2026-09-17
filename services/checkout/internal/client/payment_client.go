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

type PaymentClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewPaymentClient(baseURL string) *PaymentClient {
	return &PaymentClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *PaymentClient) CreatePayment(
	ctx context.Context,
	request models.CreatePaymentRequest,
) (*models.Payment, error) {
	body, err := json.Marshal(request)

	if err != nil {
		return nil, fmt.Errorf("encode payment: %w", err)
	}

	url := fmt.Sprintf(
		"%s/payment/api/payments",
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
		return nil, fmt.Errorf("create payment: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf(
			"payment service returned status %d",
			resp.StatusCode,
		)
	}

	var payment models.Payment

	if err := json.NewDecoder(resp.Body).Decode(&payment); err != nil {
		return nil, fmt.Errorf("decode payment: %w", err)
	}

	return &payment, nil
}
