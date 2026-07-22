package middleware

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// AuthMiddleware handles PIN-based session authentication.
type AuthMiddleware struct {
	sessionSecret string
	maxAttempts   int
	lockoutMins   int
	mu            sync.RWMutex
	attempts      map[string]int
	lockedUntil   map[string]time.Time
}

// NewAuthMiddleware creates a new auth middleware.
func NewAuthMiddleware(secret string, maxAttempts, lockoutMins int) *AuthMiddleware {
	return &AuthMiddleware{
		sessionSecret: secret,
		maxAttempts:   maxAttempts,
		lockoutMins:   lockoutMins,
		attempts:      make(map[string]int),
		lockedUntil:   make(map[string]time.Time),
	}
}

// RequireAuth checks for a valid session cookie.
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("pos_session")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// Simple validation: token must match expected format
		if len(cookie.Value) < 10 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CheckPIN validates a PIN and returns a session token or error.
func (m *AuthMiddleware) CheckPIN(pin, clientIP string) (string, error) {
	m.mu.RLock()
	if until, ok := m.lockedUntil[clientIP]; ok && time.Now().Before(until) {
		m.mu.RUnlock()
		remaining := time.Until(until).Minutes()
		return "", fmt.Errorf("bloqueado. Espera %.0f minutos", remaining)
	}
	m.mu.RUnlock()

	pinHash := HashPIN(pin)

	// Known PINs (in production: query DB)
	validPins := map[string]string{
		"03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4": "Admin",
		"a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3": "Maria Cajera",
	}

	nombre, ok := validPins[pinHash]
	if !ok {
		m.mu.Lock()
		m.attempts[clientIP]++
		if m.attempts[clientIP] >= m.maxAttempts {
			m.lockedUntil[clientIP] = time.Now().Add(time.Duration(m.lockoutMins) * time.Minute)
			m.attempts[clientIP] = 0
		}
		m.mu.Unlock()
		return "", fmt.Errorf("PIN incorrecto")
	}

	// Reset attempts on success
	m.mu.Lock()
	delete(m.attempts, clientIP)
	delete(m.lockedUntil, clientIP)
	m.mu.Unlock()

	token := m.createToken(nombre)
	return token, nil
}

func (m *AuthMiddleware) createToken(nombre string) string {
	data := fmt.Sprintf("%s|%d", nombre, time.Now().Unix())
	hash := sha256.Sum256([]byte(data + m.sessionSecret))
	return fmt.Sprintf("%s|%x", data, hash[:16])
}

// HashPIN returns SHA-256 hash of a PIN string.
func HashPIN(pin string) string {
	h := sha256.Sum256([]byte(pin))
	return fmt.Sprintf("%x", h)
}
