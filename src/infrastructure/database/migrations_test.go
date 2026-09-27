package database

import (
	"database/sql"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestMigrationsRestartPreservesSeed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "demo.db")
	first, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	tables := []string{"clientes", "productos", "ventas", "venta_items"}
	counts := make(map[string]int)
	for _, table := range tables {
		var count int
		if err := first.RW.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			t.Fatal(err)
		}
		counts[table] = count
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := New(path)
	if err != nil {
		t.Fatalf("second startup: %v", err)
	}
	defer second.Close()
	for _, table := range tables {
		var got int
		if err := second.RW.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&got); err != nil {
			t.Fatal(err)
		}
		if got != counts[table] {
			t.Errorf("%s count = %d, want %d", table, got, counts[table])
		}
	}
}

func TestMigrationFailureRollsBackAndRetries(t *testing.T) {
	rw, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "rollback.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer rw.Close()
	db := &DB{RW: rw}
	source := fstest.MapFS{"migrations/001.sql": &fstest.MapFile{Data: []byte(`CREATE TABLE demo (id INTEGER); INSERT INTO demo VALUES (1); INSERT INTO missing VALUES (2);`)}}
	if err := db.runMigrations(source); err == nil {
		t.Fatal("expected failed migration")
	}
	var count int
	if err := rw.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE name='demo'`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("failed migration left partial schema")
	}
	if err := rw.QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("failed migration was recorded")
	}
	source["migrations/001.sql"].Data = []byte(`CREATE TABLE demo (id INTEGER); INSERT INTO demo VALUES (1);`)
	if err := db.runMigrations(source); err != nil {
		t.Fatalf("retry: %v", err)
	}
	if err := db.runMigrations(source); err != nil {
		t.Fatalf("repeat: %v", err)
	}
	if err := rw.QueryRow(`SELECT COUNT(*) FROM demo`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("retry/repeat count = %d, want 1", count)
	}
}

func TestLegacyDatabaseIsNotModified(t *testing.T) {
	rw, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "legacy.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer rw.Close()
	if _, err := rw.Exec(`CREATE TABLE existing (id INTEGER); INSERT INTO existing VALUES (7)`); err != nil {
		t.Fatal(err)
	}
	db := &DB{RW: rw}
	if err := db.RunMigrations(); err == nil {
		t.Fatal("expected unsupported legacy error")
	}
	var count int
	if err := rw.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE name='schema_migrations'`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("legacy database was assigned a ledger")
	}
	var id int
	if err := rw.QueryRow(`SELECT id FROM existing`).Scan(&id); err != nil {
		t.Fatal(err)
	}
	if id != 7 {
		t.Fatalf("existing data changed: %d", id)
	}
}
