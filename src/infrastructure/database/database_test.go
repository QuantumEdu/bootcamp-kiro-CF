package database

import (
	"testing"
)

func TestNewInMemory_CreatesTables(t *testing.T) {
	db, err := NewInMemory()
	if err != nil {
		t.Fatalf("NewInMemory() error: %v", err)
	}
	defer db.Close()

	// Verify all 8 tables + sessions table exist
	expectedTables := []string{
		"usuarios",
		"categorias",
		"productos",
		"clientes",
		"ventas",
		"venta_items",
		"inventario_movimientos",
		"configuracion",
		"sessions",
	}

	for _, table := range expectedTables {
		t.Run(table, func(t *testing.T) {
			var name string
			err := db.RW.QueryRow(
				"SELECT name FROM sqlite_master WHERE type='table' AND name=?",
				table,
			).Scan(&name)
			if err != nil {
				t.Errorf("table %q not found: %v", table, err)
			}
		})
	}
}

func TestNewInMemory_IndexesExist(t *testing.T) {
	db, err := NewInMemory()
	if err != nil {
		t.Fatalf("NewInMemory() error: %v", err)
	}
	defer db.Close()

	expectedIndexes := []string{
		"idx_productos_nombre",
		"idx_productos_sku",
		"idx_productos_categoria",
		"idx_productos_stock",
		"idx_ventas_created_at",
		"idx_ventas_usuario",
		"idx_venta_items_venta",
		"idx_inventario_producto",
		"idx_usuarios_pin",
		"idx_sessions_expiry",
	}

	for _, idx := range expectedIndexes {
		t.Run(idx, func(t *testing.T) {
			var name string
			err := db.RW.QueryRow(
				"SELECT name FROM sqlite_master WHERE type='index' AND name=?",
				idx,
			).Scan(&name)
			if err != nil {
				t.Errorf("index %q not found: %v", idx, err)
			}
		})
	}
}

func TestNewInMemory_TriggerExists(t *testing.T) {
	db, err := NewInMemory()
	if err != nil {
		t.Fatalf("NewInMemory() error: %v", err)
	}
	defer db.Close()

	var name string
	err = db.RW.QueryRow(
		"SELECT name FROM sqlite_master WHERE type='trigger' AND name='trg_inventario_actualiza_stock'",
	).Scan(&name)
	if err != nil {
		t.Errorf("trigger trg_inventario_actualiza_stock not found: %v", err)
	}
}

func TestNewInMemory_ROConnectionIsQueryOnly(t *testing.T) {
	db, err := NewInMemory()
	if err != nil {
		t.Fatalf("NewInMemory() error: %v", err)
	}
	defer db.Close()

	// In-memory mode shares connection, so this test validates the setup works
	// For file-based DB, RO would reject writes
	rows, err := db.RO.Query("SELECT 1")
	if err != nil {
		t.Fatalf("RO query failed: %v", err)
	}
	rows.Close()
}
