package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alisina-1231/ecommerce-platform/services/payment/internal/database"
	"github.com/alisina-1231/ecommerce-platform/services/payment/internal/handler"
	"github.com/alisina-1231/ecommerce-platform/services/payment/internal/middleware"
	"github.com/alisina-1231/ecommerce-platform/services/payment/internal/repository"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	db, err := database.NewPool(ctx)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	paymentRepository := repository.NewPaymentRepository(db)
	paymentHandler := handler.NewPaymentHandler(paymentRepository)

	router := chi.NewRouter()

	router.Use(middleware.Logging)
	router.Use(middleware.CORS)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"payment"}`))
	})

	router.Route("/payment/api", func(router chi.Router) {
		router.Post("/payments", paymentHandler.CreatePayment)
		router.Get("/payments", paymentHandler.GetUserPayments)
		router.Get("/payments/{paymentID}", paymentHandler.GetPayment)
		router.Put(
			"/payments/{paymentID}/status",
			paymentHandler.UpdatePaymentStatus,
		)
		router.Delete("/payments/{paymentID}", paymentHandler.DeletePayment)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("payment service listening on port %s", port)

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
