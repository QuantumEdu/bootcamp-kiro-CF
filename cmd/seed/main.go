package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/config"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/database"
	"github.com/joho/godotenv"
)

// product holds the info needed to generate sale items.
type product struct {
	id    int
	price float64
}

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	log.Println("🌱 POS Seed Tool — opening database...")

	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Database ready (migrations applied automatically)")

	// Count existing data from migration seed
	migrationSales := countRows(db.RW, "ventas")
	migrationProducts := countRows(db.RW, "productos")
	migrationUsers := countRows(db.RW, "usuarios")

	log.Printf("📊 Existing data: %d ventas, %d productos, %d usuarios", migrationSales, migrationProducts, migrationUsers)

	// Generate additional demo sales for the last 30 days
	added, err := generateAdditionalSales(db.RW)
	if err != nil {
		log.Fatalf("failed to generate additional sales: %v", err)
	}

	// Print summary
	totalSales := countRows(db.RW, "ventas")
	totalItems := countRows(db.RW, "venta_items")

	fmt.Println("\n════════════════════════════════════════")
	fmt.Println("  🌱 Seed Complete")
	fmt.Println("════════════════════════════════════════")
	fmt.Printf("  Usuarios:    %d\n", migrationUsers)
	fmt.Printf("  Productos:   %d\n", migrationProducts)
	fmt.Printf("  Ventas:      %d (migration: %d + generated: %d)\n", totalSales, migrationSales, added)
	fmt.Printf("  Items:       %d\n", totalItems)
	fmt.Printf("  Periodo:     últimos 30 días\n")
	fmt.Println("════════════════════════════════════════")
}

// generateAdditionalSales creates 40+ realistic sales spanning the last 30 days.
func generateAdditionalSales(db *sql.DB) (int, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Products available for sale (id, price from 003_seed.sql)
	products := []product{
		{1, 22.00}, {2, 15.00}, {3, 28.00}, {4, 22.00}, {5, 32.00},
		{6, 26.00}, {7, 38.00}, {8, 55.00}, {9, 58.00}, {10, 22.00},
		{11, 28.00}, {12, 18.00}, {13, 22.00}, {14, 12.00}, {15, 15.00},
		{16, 18.00}, {17, 22.00}, {18, 35.00}, {19, 28.00}, {20, 35.00},
		{21, 42.00}, {22, 12.00}, {23, 22.00}, {24, 18.00}, {25, 25.00},
		{26, 20.00}, {27, 35.00}, {28, 95.00}, {29, 45.00}, {30, 38.00},
	}

	userIDs := []int{1, 2}
	clientIDs := []sql.NullInt64{
		{Int64: 1, Valid: true},
		{Int64: 2, Valid: true},
		{Int64: 3, Valid: true},
		{Int64: 4, Valid: true},
		{Int64: 5, Valid: true},
		{Valid: false}, // NULL — walk-in customer
		{Valid: false},
	}
	payMethods := []string{"efectivo", "tarjeta", "transferencia", "efectivo", "efectivo", "tarjeta"}

	now := time.Now()
	salesCount := 45

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("starting transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	insertSale, err := tx.Prepare(`INSERT INTO ventas (usuario_id, cliente_id, total, metodo_pago, created_at) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("preparing sale insert: %w", err)
	}
	defer insertSale.Close()

	insertItem, err := tx.Prepare(`INSERT INTO venta_items (venta_id, producto_id, cantidad, precio_unitario, subtotal) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("preparing item insert: %w", err)
	}
	defer insertItem.Close()

	added := 0
	for i := 0; i < salesCount; i++ {
		// Spread sales across 30 days with some variation in time of day
		daysAgo := rng.Intn(30)
		hoursOffset := rng.Intn(12) + 8 // 8am to 8pm
		saleTime := now.AddDate(0, 0, -daysAgo).Truncate(24 * time.Hour).Add(time.Duration(hoursOffset) * time.Hour)
		saleTimeStr := saleTime.Format("2006-01-02 15:04:05")

		userID := userIDs[rng.Intn(len(userIDs))]
		client := clientIDs[rng.Intn(len(clientIDs))]
		method := payMethods[rng.Intn(len(payMethods))]

		// Each sale has 1-5 items
		numItems := rng.Intn(5) + 1
		var total float64
		type saleItem struct {
			productID int
			qty       int
			price     float64
			subtotal  float64
		}
		items := make([]saleItem, 0, numItems)

		// Pick random products (no duplicates in same sale)
		picked := make(map[int]bool)
		for j := 0; j < numItems; j++ {
			p := products[rng.Intn(len(products))]
			if picked[p.id] {
				continue
			}
			picked[p.id] = true
			qty := rng.Intn(3) + 1
			subtotal := float64(qty) * p.price
			total += subtotal
			items = append(items, saleItem{p.id, qty, p.price, subtotal})
		}

		if len(items) == 0 {
			continue
		}

		// Insert sale
		res, err := insertSale.Exec(userID, client, total, method, saleTimeStr)
		if err != nil {
			return added, fmt.Errorf("inserting sale %d: %w", i, err)
		}

		saleID, err := res.LastInsertId()
		if err != nil {
			return added, fmt.Errorf("getting sale id: %w", err)
		}

		// Insert items
		for _, item := range items {
			if _, err := insertItem.Exec(saleID, item.productID, item.qty, item.price, item.subtotal); err != nil {
				return added, fmt.Errorf("inserting item for sale %d: %w", saleID, err)
			}
		}

		added++
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("committing transaction: %w", err)
	}

	log.Printf("🎲 Generated %d additional sales across 30 days", added)
	return added, nil
}

// countRows returns the number of rows in a table.
func countRows(db *sql.DB, table string) int {
	var count int
	// Table names are hardcoded constants, not user input — safe from injection.
	row := db.QueryRow("SELECT COUNT(*) FROM " + table) //nolint:gosec
	if err := row.Scan(&count); err != nil {
		return 0
	}
	return count
}
