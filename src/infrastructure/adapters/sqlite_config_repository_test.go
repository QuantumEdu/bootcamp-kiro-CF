package adapters_test

import (
	"context"
	"testing"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/adapters"
)

func TestSQLiteConfigRepository_Get_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteConfigRepository(db.RW)

	val, err := repo.Get(context.Background(), "nonexistent_key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty string for missing key, got %q", val)
	}
}

func TestSQLiteConfigRepository_SetAndGet(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteConfigRepository(db.RW)
	ctx := context.Background()

	// Set a value
	err := repo.Set(ctx, "openrouter_api_key", "encrypted_value_123")
	if err != nil {
		t.Fatalf("unexpected error on Set: %v", err)
	}

	// Get it back
	val, err := repo.Get(ctx, "openrouter_api_key")
	if err != nil {
		t.Fatalf("unexpected error on Get: %v", err)
	}
	if val != "encrypted_value_123" {
		t.Errorf("expected %q, got %q", "encrypted_value_123", val)
	}
}

func TestSQLiteConfigRepository_Set_Upsert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteConfigRepository(db.RW)
	ctx := context.Background()

	// Set initial value
	if err := repo.Set(ctx, "api_key", "value_v1"); err != nil {
		t.Fatalf("first Set failed: %v", err)
	}

	// Upsert with new value
	if err := repo.Set(ctx, "api_key", "value_v2"); err != nil {
		t.Fatalf("second Set (upsert) failed: %v", err)
	}

	// Verify updated value
	val, err := repo.Get(ctx, "api_key")
	if err != nil {
		t.Fatalf("Get after upsert failed: %v", err)
	}
	if val != "value_v2" {
		t.Errorf("expected %q after upsert, got %q", "value_v2", val)
	}
}

func TestSQLiteConfigRepository_MultipleKeys(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteConfigRepository(db.RW)
	ctx := context.Background()

	// Set multiple keys
	if err := repo.Set(ctx, "key_a", "alpha"); err != nil {
		t.Fatalf("Set key_a failed: %v", err)
	}
	if err := repo.Set(ctx, "key_b", "beta"); err != nil {
		t.Fatalf("Set key_b failed: %v", err)
	}

	// Verify each key returns its own value
	valA, err := repo.Get(ctx, "key_a")
	if err != nil {
		t.Fatalf("Get key_a failed: %v", err)
	}
	if valA != "alpha" {
		t.Errorf("expected %q for key_a, got %q", "alpha", valA)
	}

	valB, err := repo.Get(ctx, "key_b")
	if err != nil {
		t.Fatalf("Get key_b failed: %v", err)
	}
	if valB != "beta" {
		t.Errorf("expected %q for key_b, got %q", "beta", valB)
	}
}
