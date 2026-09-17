package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/alisina-1231/ecommerce-platform/services/order/internal/models"
	"github.com/alisina-1231/ecommerce-platform/services/order/internal/repository"

	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	repository *repository.OrderRepository
}

func NewOrderHandler(repository *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{
		repository: repository,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var request models.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if request.UserID == "" ||
		request.TotalAmount <= 0 ||
		request.Currency == "" ||
		request.ShippingAddress == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "user_id, total_amount, currency and shipping_address are required",
		})
		return
	}

	order, err := h.repository.Create(r.Context(), request)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to create order",
		})
		return
	}

	writeJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(
		chi.URLParam(r, "orderID"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid order_id",
		})
		return
	}

	order, err := h.repository.GetByID(r.Context(), id)

	if errors.Is(err, repository.ErrOrderNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "order not found",
		})
		return
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to get order",
		})
		return
	}

	writeJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(
		r.URL.Query().Get("user_id"),
	)

	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "user_id is required",
		})
		return
	}

	orders, err := h.repository.GetByUserID(
		r.Context(),
		userID,
	)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to get orders",
		})
		return
	}

	writeJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) UpdateOrderStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(
		chi.URLParam(r, "orderID"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid order_id",
		})
		return
	}

	var request models.UpdateOrderStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	status := strings.ToUpper(
		strings.TrimSpace(request.Status),
	)

	allowedStatuses := map[string]bool{
		"PENDING":    true,
		"CONFIRMED":  true,
		"PROCESSING": true,
		"SHIPPED":    true,
		"DELIVERED":  true,
		"CANCELLED":  true,
	}

	if !allowedStatuses[status] {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid order status",
		})
		return
	}

	order, err := h.repository.UpdateStatus(
		r.Context(),
		id,
		status,
	)

	if errors.Is(err, repository.ErrOrderNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "order not found",
		})
		return
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to update order status",
		})
		return
	}

	writeJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) DeleteOrder(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(
		chi.URLParam(r, "orderID"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid order_id",
		})
		return
	}

	err = h.repository.Delete(r.Context(), id)

	if errors.Is(err, repository.ErrOrderNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "order not found",
		})
		return
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to delete order",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	value any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}
