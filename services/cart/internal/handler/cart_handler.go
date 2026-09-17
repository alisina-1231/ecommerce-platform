package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/alisina-1231/ecommerce-platform/services/cart/internal/repository"
)

type CartHandler struct {
	repo *repository.CartRepository
}

func NewCartHandler(
	repo *repository.CartRepository,
) *CartHandler {
	return &CartHandler{
		repo: repo,
	}
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}

func writeError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	writeJSON(
		w,
		status,
		map[string]string{
			"error": message,
		},
	)
}

// GET /health
func (h *CartHandler) Health(
	w http.ResponseWriter,
	r *http.Request,
) {
	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status":  "healthy",
			"service": "cart",
		},
	)
}

// GET /cart/api/cart?user_id=user-1
func (h *CartHandler) GetCart(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}

	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"user_id is required",
		)
		return
	}

	cart, err := h.repo.GetCart(
		r.Context(),
		userID,
	)

	if err != nil {
		log.Printf("get cart error: %v", err)

		writeError(
			w,
			http.StatusInternalServerError,
			"failed to get cart",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		cart,
	)
}

// POST /cart/api/cart/items
func (h *CartHandler) AddItem(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}

	userID := r.URL.Query().Get("user_id")
	productID := r.URL.Query().Get("product_id")
	quantity := r.URL.Query().Get("quantity")

	if userID == "" || productID == "" || quantity == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"user_id, product_id and quantity are required",
		)
		return
	}

	productIDValue, err := strconv.Atoi(productID)

	if err != nil || productIDValue <= 0 {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid product_id",
		)
		return
	}

	quantityValue, err := strconv.Atoi(quantity)

	if err != nil || quantityValue <= 0 {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid quantity",
		)
		return
	}

	cart, err := h.repo.AddItem(
		r.Context(),
		userID,
		productIDValue,
		quantityValue,
	)

	if err != nil {
		log.Printf("add item error: %v", err)

		writeError(
			w,
			http.StatusInternalServerError,
			"failed to add item",
		)

		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		cart,
	)
}

// PUT /cart/api/cart/items
func (h *CartHandler) UpdateItem(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPut {
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}

	userID := r.URL.Query().Get("user_id")
	productID := r.URL.Query().Get("product_id")
	quantity := r.URL.Query().Get("quantity")

	if userID == "" || productID == "" || quantity == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"user_id, product_id and quantity are required",
		)
		return
	}

	productIDValue, err := strconv.Atoi(productID)

	if err != nil || productIDValue <= 0 {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid product_id",
		)
		return
	}

	quantityValue, err := strconv.Atoi(quantity)

	if err != nil || quantityValue <= 0 {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid quantity",
		)
		return
	}

	cart, err := h.repo.UpdateItem(
		r.Context(),
		userID,
		productIDValue,
		quantityValue,
	)

	if err != nil {
		log.Printf("update item error: %v", err)

		writeError(
			w,
			http.StatusNotFound,
			"product not found in cart",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		cart,
	)
}

// DELETE /cart/api/cart/items
func (h *CartHandler) RemoveItem(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodDelete {
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}

	userID := r.URL.Query().Get("user_id")
	productID := r.URL.Query().Get("product_id")

	if userID == "" || productID == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"user_id and product_id are required",
		)
		return
	}

	productIDValue, err := strconv.Atoi(productID)

	if err != nil || productIDValue <= 0 {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid product_id",
		)
		return
	}

	cart, err := h.repo.RemoveItem(
		r.Context(),
		userID,
		productIDValue,
	)

	if err != nil {
		log.Printf("remove item error: %v", err)

		writeError(
			w,
			http.StatusNotFound,
			"product not found in cart",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		cart,
	)
}

// DELETE /cart/api/cart
func (h *CartHandler) ClearCart(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodDelete {
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}

	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"user_id is required",
		)
		return
	}

	err := h.repo.ClearCart(
		r.Context(),
		userID,
	)

	if err != nil {
		log.Printf("clear cart error: %v", err)

		writeError(
			w,
			http.StatusInternalServerError,
			"failed to clear cart",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"message": "cart cleared successfully",
		},
	)
}
