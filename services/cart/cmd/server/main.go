package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alisina-1231/ecommerce-platform/services/cart/internal/database"
	"github.com/alisina-1231/ecommerce-platform/services/cart/internal/handler"
	"github.com/alisina-1231/ecommerce-platform/services/cart/internal/middleware"
	"github.com/alisina-1231/ecommerce-platform/services/cart/internal/repository"
)

func main() {

	ctx := context.Background()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8082"
	}

	tableName := os.Getenv("DYNAMODB_TABLE_NAME")

	if tableName == "" {
		tableName = "ecommerce-platform-dev-cart"
	}

	// --------------------------------------------------------
	// Initialize DynamoDB
	// --------------------------------------------------------

	dynamoClient, err := database.NewDynamoDBClient(ctx)

	if err != nil {
		log.Fatalf(
			"failed to initialize DynamoDB: %v",
			err,
		)
	}

	// --------------------------------------------------------
	// Check DynamoDB connection
	// --------------------------------------------------------

	if err := database.CheckConnection(
		ctx,
		dynamoClient,
		tableName,
	); err != nil {

		log.Fatalf(
			"DynamoDB connection failed: %v",
			err,
		)
	}

	log.Printf(
		"connected to DynamoDB table: %s",
		tableName,
	)

	// --------------------------------------------------------
	// Repository
	// --------------------------------------------------------

	cartRepository := repository.NewCartRepository(
		dynamoClient,
		tableName,
	)

	// --------------------------------------------------------
	// Handler
	// --------------------------------------------------------

	cartHandler := handler.NewCartHandler(
		cartRepository,
	)

	// --------------------------------------------------------
	// Routes
	// --------------------------------------------------------

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/health",
		cartHandler.Health,
	)

	mux.HandleFunc(
		"/cart/api/cart",
		func(w http.ResponseWriter, r *http.Request) {

			switch r.Method {

			case http.MethodGet:
				cartHandler.GetCart(w, r)

			case http.MethodDelete:
				cartHandler.ClearCart(w, r)

			default:
				http.Error(
					w,
					"method not allowed",
					http.StatusMethodNotAllowed,
				)

			}
		},
	)

	mux.HandleFunc(
		"/cart/api/cart/items",
		func(w http.ResponseWriter, r *http.Request) {

			switch r.Method {

			case http.MethodPost:
				cartHandler.AddItem(w, r)

			case http.MethodPut:
				cartHandler.UpdateItem(w, r)

			case http.MethodDelete:
				cartHandler.RemoveItem(w, r)

			default:
				http.Error(
					w,
					"method not allowed",
					http.StatusMethodNotAllowed,
				)

			}
		},
	)

	// --------------------------------------------------------
	// HTTP Server
	// --------------------------------------------------------

	server := &http.Server{
		Addr: ":" + port,

		Handler: middleware.CORS(
			middleware.Logging(mux),
		),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf(
		"cart service listening on port %s",
		port,
	)

	if err := server.ListenAndServe(); err != nil {

		log.Fatalf(
			"server stopped: %v",
			err,
		)
	}
}
