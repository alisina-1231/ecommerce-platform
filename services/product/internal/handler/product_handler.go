package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alisina-1231/ecommerce-platform/services/product/internal/models"
	"github.com/alisina-1231/ecommerce-platform/services/product/internal/repository"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(
	repo *repository.ProductRepository,
) *ProductHandler {
	return &ProductHandler{
		repo: repo,
	}
}

func (h *ProductHandler) Health(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	response := map[string]string{
		"status":  "healthy",
		"service": "product",
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf(
			"failed to write health response: %v",
			err,
		)
	}
}

func (h *ProductHandler) Products(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {
	case http.MethodGet:
		h.listProducts(w, r)

	case http.MethodPost:
		h.createProduct(w, r)

	default:
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func (h *ProductHandler) Product(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := productIDFromPath(r.URL.Path)

	if err != nil {
		http.Error(
			w,
			"invalid product id",
			http.StatusBadRequest,
		)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getProduct(w, r, id)

	case http.MethodPut:
		h.updateProduct(w, r, id)

	case http.MethodDelete:
		h.deleteProduct(w, r, id)

	default:
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func (h *ProductHandler) listProducts(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, cancel := requestContext(r)
	defer cancel()

	products, err := h.repo.List(ctx)
	if err != nil {
		log.Printf(
			"failed to list products: %v",
			err,
		)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		products,
	)
}

func (h *ProductHandler) createProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	var product models.Product

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&product); err != nil {
		http.Error(
			w,
			"invalid JSON body",
			http.StatusBadRequest,
		)
		return
	}

	if err := validateProduct(&product); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	if err := h.repo.Create(ctx, &product); err != nil {
		log.Printf(
			"failed to create product: %v",
			err,
		)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		product,
	)
}

func (h *ProductHandler) getProduct(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	ctx, cancel := requestContext(r)
	defer cancel()

	product, err := h.repo.GetByID(ctx, id)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrProductNotFound,
		) {
			http.Error(
				w,
				"product not found",
				http.StatusNotFound,
			)
			return
		}

		log.Printf(
			"failed to get product: %v",
			err,
		)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		product,
	)
}

func (h *ProductHandler) updateProduct(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	var product models.Product

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(
			w,
			"invalid JSON body",
			http.StatusBadRequest,
		)
		return
	}

	product.ID = id

	if err := validateProduct(&product); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	if err := h.repo.Update(ctx, &product); err != nil {
		if errors.Is(
			err,
			repository.ErrProductNotFound,
		) {
			http.Error(
				w,
				"product not found",
				http.StatusNotFound,
			)
			return
		}

		log.Printf(
			"failed to update product: %v",
			err,
		)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		product,
	)
}

func (h *ProductHandler) deleteProduct(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	ctx, cancel := requestContext(r)
	defer cancel()

	if err := h.repo.Delete(ctx, id); err != nil {
		if errors.Is(
			err,
			repository.ErrProductNotFound,
		) {
			http.Error(
				w,
				"product not found",
				http.StatusNotFound,
			)
			return
		}

		log.Printf(
			"failed to delete product: %v",
			err,
		)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func productIDFromPath(
	path string,
) (int64, error) {
	const prefix = "/product/api/products/"

	idString := strings.TrimPrefix(path, prefix)

	if idString == "" || idString == path {
		return 0, errors.New("missing product id")
	}

	return strconv.ParseInt(
		idString,
		10,
		64,
	)
}

func validateProduct(
	product *models.Product,
) error {
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("name is required")
	}

	if product.Price < 0 {
		return errors.New("price cannot be negative")
	}

	if len(product.Currency) != 3 {
		return errors.New(
			"currency must be a 3-letter code",
		)
	}

	product.Currency = strings.ToUpper(
		strings.TrimSpace(product.Currency),
	)

	return nil
}

func requestContext(
	r *http.Request,
) (context.Context, context.CancelFunc) {
	return context.WithTimeout(
		r.Context(),
		5*time.Second,
	)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data interface{},
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf(
			"failed to encode JSON response: %v",
			err,
		)
	}
}
