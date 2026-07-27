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
		log.Fatal("DATABASE_URL required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer conn.Close(ctx)

	// Reset all sequences to max(id)+1
	tables := []string{"usuarios", "categorias", "productos", "clientes", "ventas", "venta_items", "inventario_movimientos", "configuracion"}
	for _, t := range tables {
		sql := fmt.Sprintf("SELECT setval(pg_get_serial_sequence('%s', 'id'), COALESCE(MAX(id), 1)) FROM %s", t, t)
		var val int64
		conn.QueryRow(ctx, sql).Scan(&val)
		fmt.Printf("✅ %s sequence → %d\n", t, val)
	}
	fmt.Println("Done!")
}
