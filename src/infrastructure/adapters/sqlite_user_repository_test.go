package adapters_test

import (
	"context"
	"testing"
	"time"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/adapters"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/database"
)

// setupTestDB creates an in-memory database with migrations applied.
// The 003_seed.sql migration already inserts 2 users:
//   - ID 1: "Admin" (admin, pin_hash=sha256("1234"))
//   - ID 2: "Maria Cajera" (cajero, pin_hash=sha256("123"))
func setupTestDB(t *testing.T) *database.DB {
	t.Helper()

	db, err := database.NewInMemory()
	if err != nil {
		t.Fatalf("creating in-memory DB: %v", err)
	}

	return db
}

func TestSQLiteUserRepository_FindAll_ReturnsSeededUsers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteUserRepository(db.RW)

	users, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}

	if len(users) < 2 {
		t.Errorf("FindAll() got %d users, want at least 2", len(users))
	}

	// Verify first user (from seed migration)
	if users[0].Nombre != "Admin" {
		t.Errorf("users[0].Nombre = %q, want %q", users[0].Nombre, "Admin")
	}
	if users[0].Rol != "admin" {
		t.Errorf("users[0].Rol = %q, want %q", users[0].Rol, "admin")
	}
	if !users[0].Activo {
		t.Error("users[0].Activo = false, want true")
	}

	// Verify second user
	if users[1].Nombre != "Maria Cajera" {
		t.Errorf("users[1].Nombre = %q, want %q", users[1].Nombre, "Maria Cajera")
	}
	if users[1].Rol != "cajero" {
		t.Errorf("users[1].Rol = %q, want %q", users[1].Rol, "cajero")
	}
}

func TestSQLiteUserRepository_FindByID_ExistingUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteUserRepository(db.RW)

	user, err := repo.FindByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("FindByID(1) error = %v", err)
	}

	if user.Nombre != "Admin" {
		t.Errorf("user.Nombre = %q, want %q", user.Nombre, "Admin")
	}
	if user.PINHash != "03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4" {
		t.Errorf("user.PINHash = %q, want sha256 hash of 1234", user.PINHash)
	}
	if user.Rol != "admin" {
		t.Errorf("user.Rol = %q, want %q", user.Rol, "admin")
	}
}

func TestSQLiteUserRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteUserRepository(db.RW)

	_, err := repo.FindByID(context.Background(), 999)
	if err == nil {
		t.Fatal("FindByID(999) expected error, got nil")
	}
}

func TestSQLiteUserRepository_FindByPINHash_ExistingUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteUserRepository(db.RW)

	// Maria Cajera's pin_hash is sha256("123")
	user, err := repo.FindByPINHash(context.Background(), "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3")
	if err != nil {
		t.Fatalf("FindByPINHash() error = %v", err)
	}

	if user.Nombre != "Maria Cajera" {
		t.Errorf("user.Nombre = %q, want %q", user.Nombre, "Maria Cajera")
	}
	if user.Rol != "cajero" {
		t.Errorf("user.Rol = %q, want %q", user.Rol, "cajero")
	}
}

func TestSQLiteUserRepository_IncrementFailedAttempts(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteUserRepository(db.RW)
	ctx := context.Background()

	// User starts with 0 attempts
	user, _ := repo.FindByID(ctx, 1)
	if user.IntentosFallidos != 0 {
		t.Fatalf("initial IntentosFallidos = %d, want 0", user.IntentosFallidos)
	}

	// Increment
	if err := repo.IncrementFailedAttempts(ctx, 1); err != nil {
		t.Fatalf("IncrementFailedAttempts() error = %v", err)
	}

	// Verify incremented
	user, _ = repo.FindByID(ctx, 1)
	if user.IntentosFallidos != 1 {
		t.Errorf("IntentosFallidos after increment = %d, want 1", user.IntentosFallidos)
	}

	// Increment again
	if err := repo.IncrementFailedAttempts(ctx, 1); err != nil {
		t.Fatalf("IncrementFailedAttempts() second call error = %v", err)
	}

	user, _ = repo.FindByID(ctx, 1)
	if user.IntentosFallidos != 2 {
		t.Errorf("IntentosFallidos after second increment = %d, want 2", user.IntentosFallidos)
	}
}

func TestSQLiteUserRepository_Lock_SetsBloqueadoHasta(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteUserRepository(db.RW)
	ctx := context.Background()

	lockUntil := time.Now().Add(5 * time.Minute).Truncate(time.Second)

	if err := repo.Lock(ctx, 1, lockUntil); err != nil {
		t.Fatalf("Lock() error = %v", err)
	}

	user, err := repo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("FindByID() after lock error = %v", err)
	}

	if user.BloqueadoHasta.IsZero() {
		t.Fatal("BloqueadoHasta is zero after Lock()")
	}

	// Compare with second precision
	if !user.BloqueadoHasta.Equal(lockUntil) {
		t.Errorf("BloqueadoHasta = %v, want %v", user.BloqueadoHasta, lockUntil)
	}

	// The user should be locked
	if !user.IsLocked() {
		t.Error("user.IsLocked() = false after Lock(), want true")
	}
}

func TestSQLiteUserRepository_ResetAttempts_ClearsCounterAndLockout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteUserRepository(db.RW)
	ctx := context.Background()

	// First increment and lock the user
	_ = repo.IncrementFailedAttempts(ctx, 1)
	_ = repo.IncrementFailedAttempts(ctx, 1)
	_ = repo.Lock(ctx, 1, time.Now().Add(5*time.Minute))

	// Verify locked state
	user, _ := repo.FindByID(ctx, 1)
	if user.IntentosFallidos != 2 {
		t.Fatalf("setup: IntentosFallidos = %d, want 2", user.IntentosFallidos)
	}
	if user.BloqueadoHasta.IsZero() {
		t.Fatal("setup: BloqueadoHasta is zero after Lock()")
	}

	// Reset
	if err := repo.ResetAttempts(ctx, 1); err != nil {
		t.Fatalf("ResetAttempts() error = %v", err)
	}

	// Verify reset
	user, err := repo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("FindByID() after reset error = %v", err)
	}

	if user.IntentosFallidos != 0 {
		t.Errorf("IntentosFallidos after reset = %d, want 0", user.IntentosFallidos)
	}
	if !user.BloqueadoHasta.IsZero() {
		t.Errorf("BloqueadoHasta after reset = %v, want zero", user.BloqueadoHasta)
	}
	if user.IsLocked() {
		t.Error("user.IsLocked() = true after ResetAttempts(), want false")
	}
}

func TestSQLiteUserRepository_IncrementFailedAttempts_NonExistentUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteUserRepository(db.RW)

	err := repo.IncrementFailedAttempts(context.Background(), 999)
	if err == nil {
		t.Fatal("IncrementFailedAttempts(999) expected error, got nil")
	}
}
