package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/alisina-1231/ecommerce-platform/services/payment/internal/models"
	"github.com/alisina-1231/ecommerce-platform/services/payment/internal/repository"

	"github.com/go-chi/chi/v5"
)

type PaymentHandler struct {
	repository *repository.PaymentRepository
}

func NewPaymentHandler(
	repository *repository.PaymentRepository,
) *PaymentHandler {
	return &PaymentHandler{
		repository: repository,
	}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var request models.CreatePaymentRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if request.OrderID <= 0 ||
		request.UserID == "" ||
		request.Amount <= 0 ||
		request.Currency == "" ||
		request.PaymentMethod == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "order_id, user_id, amount, currency and payment_method are required",
		})
		return
	}

	payment, err := h.repository.Create(r.Context(), request)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to create payment",
		})
		return
	}

	writeJSON(w, http.StatusCreated, payment)
}

func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "paymentID"), 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid payment_id",
		})
		return
	}

	payment, err := h.repository.GetByID(r.Context(), id)
	if errors.Is(err, repository.ErrPaymentNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "payment not found",
		})
		return
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to get payment",
		})
		return
	}

	writeJSON(w, http.StatusOK, payment)
}

func (h *PaymentHandler) GetUserPayments(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))

	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "user_id is required",
		})
		return
	}

	payments, err := h.repository.GetByUserID(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to get payments",
		})
		return
	}

	writeJSON(w, http.StatusOK, payments)
}

func (h *PaymentHandler) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "paymentID"), 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid payment_id",
		})
		return
	}

	var request models.UpdatePaymentStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	allowedStatuses := map[string]bool{
		"PENDING":    true,
		"PROCESSING": true,
		"COMPLETED":  true,
		"FAILED":     true,
		"CANCELLED":  true,
	}

	request.Status = strings.ToUpper(strings.TrimSpace(request.Status))

	if !allowedStatuses[request.Status] {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid payment status",
		})
		return
	}

	payment, err := h.repository.UpdateStatus(
		r.Context(),
		id,
		request,
	)

	if errors.Is(err, repository.ErrPaymentNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "payment not found",
		})
		return
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to update payment status",
		})
		return
	}

	writeJSON(w, http.StatusOK, payment)
}

func (h *PaymentHandler) DeletePayment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "paymentID"), 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid payment_id",
		})
		return
	}

	err = h.repository.Delete(r.Context(), id)
	if errors.Is(err, repository.ErrPaymentNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "payment not found",
		})
		return
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to delete payment",
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
