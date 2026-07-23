// Package value_objects contains domain value objects for the POS system.
package value_objects

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// Domain errors for PIN validation and comparison.
var (
	ErrPINMismatch   = errors.New("PIN does not match")
	ErrPINTooShort   = errors.New("PIN must be at least 4 digits")
	ErrPINTooLong    = errors.New("PIN must be at most 6 digits")
	ErrPINNotNumeric = errors.New("PIN must contain only digits")
)

// bcryptCost defines the bcrypt hashing cost. Standard for PIN in POS context.
const bcryptCost = 10

// HashPIN takes a plain PIN string and returns its bcrypt hash.
// The caller should validate the PIN format before hashing.
func HashPIN(plainPIN string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPIN), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePIN verifies that a plain PIN matches the provided bcrypt hash.
// Returns nil on success, or ErrPINMismatch if they do not match.
func ComparePIN(hashedPIN, plainPIN string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPIN), []byte(plainPIN))
	if err != nil {
		return ErrPINMismatch
	}
	return nil
}

// ValidatePINFormat checks that a PIN is 4-6 digits only.
// Returns a specific domain error describing the validation failure.
func ValidatePINFormat(pin string) error {
	if len(pin) < 4 {
		return ErrPINTooShort
	}
	if len(pin) > 6 {
		return ErrPINTooLong
	}
	for _, r := range pin {
		if !unicode.IsDigit(r) {
			return ErrPINNotNumeric
		}
	}
	return nil
}
