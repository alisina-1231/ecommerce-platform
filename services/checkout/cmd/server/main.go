package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alisina-1231/ecommerce-platform/services/checkout/internal/client"
	"github.com/alisina-1231/ecommerce-platform/services/checkout/internal/handler"
	"github.com/alisina-1231/ecommerce-platform/services/checkout/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func main() {
	cartURL := getEnv(
		"CART_SERVICE_URL",
		"http://localhost:8082",
	)

	productURL := getEnv(
		"PRODUCT_SERVICE_URL",
		"http://localhost:8081",
	)

	orderURL := getEnv(
		"ORDER_SERVICE_URL",
		"http://localhost:8084",
	)

	paymentURL := getEnv(
		"PAYMENT_SERVICE_URL",
		"http://localhost:8085",
	)

	cartClient := client.NewCartClient(cartURL)
	productClient := client.NewProductClient(productURL)
	orderClient := client.NewOrderClient(orderURL)
	paymentClient := client.NewPaymentClient(paymentURL)

	checkoutHandler := handler.NewCheckoutHandler(
		cartClient,
		productClient,
		orderClient,
		paymentClient,
	)

	router := chi.NewRouter()

	router.Use(middleware.Logging)
	router.Use(middleware.CORS)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(
			`{"status":"ok","service":"checkout"}`,
		))
	})

	router.Post(
		"/checkout/api/checkout",
		checkoutHandler.Checkout,
	)

	port := getEnv("PORT", "8083")

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf(
		"checkout service listening on port %s",
		port,
	)

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatalf(
			"server failed: %v",
			err,
		)
	}
}
