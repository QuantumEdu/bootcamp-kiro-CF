package middleware

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/database"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const UserRolKey contextKey = "user_rol"
const UserNameKey contextKey = "user_name"

// AuthMiddleware handles PIN-based authentication.
type AuthMiddleware struct {
	userRepo       *database.UserRepo
	sessionSecret  string
	maxAttempts    int
	lockoutMinutes int
	mu             sync.RWMutex
	attempts       map[string]*loginAttempt
}

type loginAttempt struct {
	count    int
	lockedAt time.Time
}

// NewAuthMiddleware creates a new auth middleware.
func NewAuthMiddleware(userRepo *database.UserRepo, sessionSecret string, maxAttempts, lockoutMinutes int) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo:       userRepo,
		sessionSecret:  sessionSecret,
		maxAttempts:    maxAttempts,
		lockoutMinutes: lockoutMinutes,
		attempts:       make(map[string]*loginAttempt),
	}
}

// RequireAuth is the middleware that checks for a valid session.
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Validate session (simple HMAC-based token)
		userID, rol, nombre, ok := m.validateSession(cookie.Value)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserRolKey, rol)
		ctx = context.WithValue(ctx, UserNameKey, nombre)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Login attempts to authenticate a user with their PIN.
func (m *AuthMiddleware) Login(pin, clientIP string) (string, error) {
	// Check lockout
	m.mu.RLock()
	attempt, exists := m.attempts[clientIP]
	m.mu.RUnlock()

	if exists && attempt.count >= m.maxAttempts {
		lockoutEnd := attempt.lockedAt.Add(time.Duration(m.lockoutMinutes) * time.Minute)
		if time.Now().Before(lockoutEnd) {
			remaining := time.Until(lockoutEnd).Minutes()
			return "", fmt.Errorf("cuenta bloqueada. Intenta en %.0f minutos", remaining)
		}
		// Lockout expired, reset
		m.mu.Lock()
		delete(m.attempts, clientIP)
		m.mu.Unlock()
	}

	// Hash PIN
	pinHash := HashPIN(pin)

	// Find user
	user, err := m.userRepo.GetByPinHash(pinHash)
	if err != nil {
		return "", fmt.Errorf("error interno")
	}
	if user == nil {
		// Record failed attempt
		m.recordFailedAttempt(clientIP)
		return "", fmt.Errorf("PIN incorrecto")
	}

	// Reset attempts on success
	m.mu.Lock()
	delete(m.attempts, clientIP)
	m.mu.Unlock()

	// Create session token
	token := m.createSession(user.ID, user.Rol, user.Nombre)
	return token, nil
}

func (m *AuthMiddleware) recordFailedAttempt(clientIP string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	attempt, exists := m.attempts[clientIP]
	if !exists {
		m.attempts[clientIP] = &loginAttempt{count: 1, lockedAt: time.Now()}
		return
	}
	attempt.count++
	if attempt.count >= m.maxAttempts {
		attempt.lockedAt = time.Now()
	}
}

func (m *AuthMiddleware) createSession(userID int64, rol, nombre string) string {
	data := fmt.Sprintf("%d|%s|%s|%d", userID, rol, nombre, time.Now().Unix())
	hash := sha256.Sum256([]byte(data + m.sessionSecret))
	return fmt.Sprintf("%s|%x", data, hash[:16])
}

func (m *AuthMiddleware) validateSession(token string) (int64, string, string, bool) {
	// Token format: "userID|rol|nombre|timestamp|hash"
	var userID int64
	var rol, nombre string
	var ts int64

	// Find last | (hash separator)
	lastPipe := -1
	for i := len(token) - 1; i >= 0; i-- {
		if token[i] == '|' {
			lastPipe = i
			break
		}
	}
	if lastPipe <= 0 {
		return 0, "", "", false
	}

	data := token[:lastPipe]
	providedHash := token[lastPipe+1:]

	// Verify hash
	expectedHash := sha256.Sum256([]byte(data + m.sessionSecret))
	expected := fmt.Sprintf("%x", expectedHash[:16])
	if providedHash != expected {
		return 0, "", "", false
	}

	// Parse data
	n, err := fmt.Sscanf(data, "%d|", &userID)
	if err != nil || n < 1 {
		return 0, "", "", false
	}

	// Extract parts manually for fields with potential |
	parts := splitN(data, '|', 4)
	if len(parts) < 4 {
		return 0, "", "", false
	}

	fmt.Sscanf(parts[0], "%d", &userID)
	rol = parts[1]
	nombre = parts[2]
	fmt.Sscanf(parts[3], "%d", &ts)

	// Check session age (24h max)
	sessionAge := time.Since(time.Unix(ts, 0))
	if sessionAge > 24*time.Hour {
		return 0, "", "", false
	}

	return userID, rol, nombre, true
}

// HashPIN creates a SHA-256 hash of a PIN.
func HashPIN(pin string) string {
	hash := sha256.Sum256([]byte(pin))
	return fmt.Sprintf("%x", hash)
}

func splitN(s string, sep byte, n int) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s) && len(parts) < n-1; i++ {
		if s[i] == sep {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}
