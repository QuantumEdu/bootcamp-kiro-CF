package services

import (
	"testing"
)

func TestNewCryptoService_DerivesKey(t *testing.T) {
	svc := NewCryptoService("test-secret")
	if len(svc.key) != 32 {
		t.Errorf("expected 32-byte key, got %d bytes", len(svc.key))
	}
}

func TestCryptoService_EncryptDecrypt_RoundTrip(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		plaintext string
	}{
		{"simple text", "my-secret", "hello world"},
		{"api key format", "session-secret-123", "sk-or-v1-abc123def456"},
		{"unicode", "secret", "contraseña con ñ y ü"},
		{"long string", "secret", "a]b[c{d}e(f)g*h+i.j?k^l$m|n\\o/p"},
		{"single char", "s", "x"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCryptoService(tt.secret)

			encrypted, err := svc.Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}
			if encrypted == "" {
				t.Fatal("Encrypt returned empty string")
			}
			if encrypted == tt.plaintext {
				t.Error("Encrypt returned plaintext unchanged")
			}

			decrypted, err := svc.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}
			if decrypted != tt.plaintext {
				t.Errorf("round-trip failed: got %q, want %q", decrypted, tt.plaintext)
			}
		})
	}
}

func TestCryptoService_Decrypt_InvalidHex(t *testing.T) {
	svc := NewCryptoService("secret")
	_, err := svc.Decrypt("not-valid-hex!")
	if err == nil {
		t.Error("expected error for invalid hex, got nil")
	}
}

func TestCryptoService_Decrypt_TooShort(t *testing.T) {
	svc := NewCryptoService("secret")
	_, err := svc.Decrypt("aabb")
	if err == nil {
		t.Error("expected error for too-short ciphertext, got nil")
	}
	if err != nil && err.Error() != "ciphertext too short" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCryptoService_Decrypt_WrongKey(t *testing.T) {
	svc1 := NewCryptoService("secret-one")
	svc2 := NewCryptoService("secret-two")

	encrypted, err := svc1.Encrypt("sensitive data")
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = svc2.Decrypt(encrypted)
	if err == nil {
		t.Error("expected error when decrypting with wrong key, got nil")
	}
}

func TestCryptoService_Encrypt_ProducesDifferentCiphertexts(t *testing.T) {
	svc := NewCryptoService("secret")
	plaintext := "same input"

	enc1, err := svc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("first Encrypt failed: %v", err)
	}
	enc2, err := svc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("second Encrypt failed: %v", err)
	}

	if enc1 == enc2 {
		t.Error("two encryptions of the same plaintext produced identical ciphertext (nonce reuse)")
	}
}
