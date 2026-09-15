package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Order struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Status    string  `json:"status"`
	Total     float64 `json:"total"`
	Currency  string  `json:"currency"`
	CreatedAt string  `json:"created_at"`
}

type CreateOrderRequest struct {
	UserID   string  `json:"user_id"`
	Total    float64 `json:"total"`
	Currency string  `json:"currency"`
}

var orders = []Order{}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status":  "healthy",
		"service": "order",
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to write health response: %v", err)
	}
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getOrders(w)
	case http.MethodPost:
		createOrder(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getOrders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(orders); err != nil {
		log.Printf("failed to write orders response: %v", err)
	}
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var request CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	if request.Total <= 0 {
		http.Error(w, "total must be greater than zero", http.StatusBadRequest)
		return
	}

	if request.Currency == "" {
		request.Currency = "USD"
	}

	order := Order{
		ID:        "order-" + time.Now().Format("20060102150405"),
		UserID:    request.UserID,
		Status:    "pending",
		Total:     request.Total,
		Currency:  request.Currency,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	orders = append(orders, order)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("failed to write order response: %v", err)
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/orders/api/orders", ordersHandler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMiddleware(loggingMiddleware(mux)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("order service listening on port %s", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"method=%s path=%s duration=%s",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}
