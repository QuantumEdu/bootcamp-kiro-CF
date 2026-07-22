package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/config"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (ignore error in production where env vars are set directly)
	_ = godotenv.Load()

	// Load configuration
	cfg := config.Load()

	// Initialize database (runs migrations automatically)
	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","port":%q}`, cfg.Port)
	})

	// Start server
	addr := ":" + cfg.Port
	log.Printf("POS AI-First starting on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
		os.Exit(1)
	}
}
