package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"voucher-api/internal/database"
	"voucher-api/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database connection
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize handlers
	h := handlers.NewHandler(db)

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Post("/brand", h.CreateBrand)
	r.Post("/voucher", h.CreateVoucher)
	r.Get("/voucher", h.GetVoucher)
	r.Get("/voucher/brand", h.GetVouchersByBrand)
	r.Post("/transaction/redemption", h.CreateRedemption)
	r.Get("/transaction/redemption", h.GetRedemption)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
