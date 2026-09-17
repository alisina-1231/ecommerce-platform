package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alisina-1231/ecommerce-platform/services/product/internal/database"
	"github.com/alisina-1231/ecommerce-platform/services/product/internal/handler"
	"github.com/alisina-1231/ecommerce-platform/services/product/internal/middleware"
	"github.com/alisina-1231/ecommerce-platform/services/product/internal/repository"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	productHandler := handler.NewProductHandler(productRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", productHandler.Health)
	mux.HandleFunc("/product/api/products", productHandler.Products)
	mux.HandleFunc("/product/api/products/", productHandler.Product)

	server := &http.Server{
		Addr: ":" + port,

		Handler: middleware.CORS(
			middleware.Logging(mux),
		),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("product service listening on port %s", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
