package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}
	defer conn.Close(ctx)

	// Read migration file
	sql, err := os.ReadFile("migrations/postgres/001_init.sql")
	if err != nil {
		log.Fatalf("reading migration file: %v", err)
	}

	// Execute migration
	_, err = conn.Exec(ctx, string(sql))
	if err != nil {
		log.Fatalf("executing migration: %v", err)
	}

	fmt.Println("✅ Migration completed successfully!")
}
