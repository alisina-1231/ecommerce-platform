package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alisina-1231/ecommerce-platform/services/checkout/internal/client"
	"github.com/alisina-1231/ecommerce-platform/services/checkout/internal/models"
)

type CheckoutHandler struct {
	cartClient    *client.CartClient
	productClient *client.ProductClient
	orderClient   *client.OrderClient
	paymentClient *client.PaymentClient
}

func NewCheckoutHandler(
	cartClient *client.CartClient,
	productClient *client.ProductClient,
	orderClient *client.OrderClient,
	paymentClient *client.PaymentClient,
) *CheckoutHandler {
	return &CheckoutHandler{
		cartClient:    cartClient,
		productClient: productClient,
		orderClient:   orderClient,
		paymentClient: paymentClient,
	}
}

func (h *CheckoutHandler) Checkout(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request models.CheckoutRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if request.UserID == "" ||
		request.ShippingAddress == "" ||
		request.PaymentMethod == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "user_id, shipping_address and payment_method are required",
		})
		return
	}

	ctx := r.Context()

	// 1. Get cart.
	cart, err := h.cartClient.GetCart(
		ctx,
		request.UserID,
	)

	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{
			"error": fmt.Sprintf("failed to get cart: %v", err),
		})
		return
	}

	if len(cart.Items) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "cart is empty",
		})
		return
	}

	// 2. Retrieve products and calculate total.
	var total float64
	currency := "USD"

	for _, item := range cart.Items {
		if item.Quantity <= 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid cart quantity",
			})
			return
		}

		product, err := h.productClient.GetProduct(
			ctx,
			item.ProductID,
		)

		if err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{
				"error": fmt.Sprintf(
					"failed to get product %d: %v",
					item.ProductID,
					err,
				),
			})
			return
		}

		if product.Stock < item.Quantity {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf(
					"insufficient stock for product %d",
					item.ProductID,
				),
			})
			return
		}

		total += product.Price * float64(item.Quantity)

		if product.Currency != "" {
			currency = product.Currency
		}
	}

	// 3. Create order.
	order, err := h.orderClient.CreateOrder(
		ctx,
		models.CreateOrderRequest{
			UserID:          request.UserID,
			TotalAmount:     total,
			Currency:        currency,
			ShippingAddress: request.ShippingAddress,
		},
	)

	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{
			"error": fmt.Sprintf(
				"failed to create order: %v",
				err,
			),
		})
		return
	}

	// 4. Create payment.
	payment, err := h.paymentClient.CreatePayment(
		ctx,
		models.CreatePaymentRequest{
			OrderID:       order.ID,
			UserID:        request.UserID,
			Amount:        total,
			Currency:      currency,
			PaymentMethod: request.PaymentMethod,
		},
	)

	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{
			"error": fmt.Sprintf(
				"failed to create payment: %v",
				err,
			),
		})
		return
	}

	// 5. Clear cart.
	if err := h.cartClient.ClearCart(
		ctx,
		request.UserID,
	); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{
			"error": fmt.Sprintf(
				"checkout completed but failed to clear cart: %v",
				err,
			),
		})
		return
	}

	// 6. Return checkout result.
	writeJSON(w, http.StatusCreated, models.CheckoutResponse{
		OrderID:     order.ID,
		PaymentID:   payment.ID,
		Status:      payment.Status,
		TotalAmount: total,
		Currency:    currency,
	})
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	value any,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}
