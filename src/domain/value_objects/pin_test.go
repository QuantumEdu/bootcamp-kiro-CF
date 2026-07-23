package value_objects

import (
	"errors"
	"testing"
)

func TestHashPIN_ProducesValidHash(t *testing.T) {
	pin := "1234"
	hash, err := HashPIN(pin)
	if err != nil {
		t.Fatalf("HashPIN(%q) returned unexpected error: %v", pin, err)
	}
	if hash == "" {
		t.Fatal("HashPIN returned empty hash")
	}
	if hash == pin {
		t.Fatal("HashPIN returned the plain PIN instead of a hash")
	}
}

func TestHashPIN_DifferentHashesForSamePIN(t *testing.T) {
	pin := "5678"
	hash1, err := HashPIN(pin)
	if err != nil {
		t.Fatalf("first HashPIN call failed: %v", err)
	}
	hash2, err := HashPIN(pin)
	if err != nil {
		t.Fatalf("second HashPIN call failed: %v", err)
	}
	if hash1 == hash2 {
		t.Error("expected different hashes for same PIN (bcrypt uses random salt)")
	}
}

func TestComparePIN_ValidMatch(t *testing.T) {
	pin := "4321"
	hash, err := HashPIN(pin)
	if err != nil {
		t.Fatalf("HashPIN failed: %v", err)
	}

	if err := ComparePIN(hash, pin); err != nil {
		t.Errorf("ComparePIN should succeed for matching PIN, got: %v", err)
	}
}

func TestComparePIN_WrongPIN(t *testing.T) {
	pin := "1234"
	hash, err := HashPIN(pin)
	if err != nil {
		t.Fatalf("HashPIN failed: %v", err)
	}

	err = ComparePIN(hash, "9999")
	if err == nil {
		t.Fatal("ComparePIN should return error for wrong PIN")
	}
	if !errors.Is(err, ErrPINMismatch) {
		t.Errorf("expected ErrPINMismatch, got: %v", err)
	}
}

func TestValidatePINFormat(t *testing.T) {
	tests := []struct {
		name    string
		pin     string
		wantErr error
	}{
		{"valid 4 digits", "1234", nil},
		{"valid 5 digits", "12345", nil},
		{"valid 6 digits", "123456", nil},
		{"too short - 3 digits", "123", ErrPINTooShort},
		{"too short - 1 digit", "1", ErrPINTooShort},
		{"too short - empty", "", ErrPINTooShort},
		{"too long - 7 digits", "1234567", ErrPINTooLong},
		{"too long - 10 digits", "1234567890", ErrPINTooLong},
		{"non-numeric - letters", "abcd", ErrPINNotNumeric},
		{"non-numeric - mixed", "12ab", ErrPINNotNumeric},
		{"non-numeric - special chars", "12#4", ErrPINNotNumeric},
		{"non-numeric - spaces", "12 4", ErrPINNotNumeric},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePINFormat(tt.pin)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("ValidatePINFormat(%q) = %v, want nil", tt.pin, err)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidatePINFormat(%q) = %v, want %v", tt.pin, err, tt.wantErr)
			}
		})
	}
}
