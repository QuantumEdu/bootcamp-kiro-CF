package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
)

// CryptoService handles AES-GCM encryption/decryption for sensitive data
// such as API keys stored in the configuracion table.
type CryptoService struct {
	key []byte // 32-byte AES-256 key derived from SESSION_SECRET
}

// NewCryptoService derives a 256-bit key from the secret using SHA-256.
func NewCryptoService(secret string) *CryptoService {
	hash := sha256.Sum256([]byte(secret))
	return &CryptoService{key: hash[:]}
}

// Encrypt encrypts plaintext using AES-GCM. Returns hex-encoded ciphertext.
func (s *CryptoService) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(sealed), nil
}

// Decrypt decrypts hex-encoded AES-GCM ciphertext.
func (s *CryptoService) Decrypt(ciphertextHex string) (string, error) {
	data, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
