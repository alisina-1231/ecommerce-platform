package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alisina-1231/ecommerce-platform/services/order/internal/database"
	"github.com/alisina-1231/ecommerce-platform/services/order/internal/handler"
	"github.com/alisina-1231/ecommerce-platform/services/order/internal/middleware"
	"github.com/alisina-1231/ecommerce-platform/services/order/internal/repository"

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

	orderRepository := repository.NewOrderRepository(db)
	orderHandler := handler.NewOrderHandler(orderRepository)

	router := chi.NewRouter()

	router.Use(middleware.Logging)
	router.Use(middleware.CORS)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(
			`{"status":"ok","service":"order"}`,
		))
	})

	router.Route("/order/api", func(router chi.Router) {
		router.Post("/orders", orderHandler.CreateOrder)

		router.Get("/orders", orderHandler.GetUserOrders)

		router.Get(
			"/orders/{orderID}",
			orderHandler.GetOrder,
		)

		router.Put(
			"/orders/{orderID}/status",
			orderHandler.UpdateOrderStatus,
		)

		router.Delete(
			"/orders/{orderID}",
			orderHandler.DeleteOrder,
		)
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8084"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("order service listening on port %s", port)

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
