package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a sliding window rate limiter per IP.
type RateLimiter struct {
	mu       sync.RWMutex
	clients  map[string]*clientWindow
	limit    int
	window   time.Duration
	cleanup  time.Duration
}

type clientWindow struct {
	timestamps []time.Time
	lastSeen   time.Time
}

// NewRateLimiter creates a new rate limiter.
// limit: max requests per window
// window: time window duration
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*clientWindow),
		limit:   limit,
		window:  window,
		cleanup: window * 2,
	}

	// Background cleanup of stale entries
	go rl.cleanupLoop()

	return rl
}

// Middleware returns an HTTP middleware that enforces rate limiting.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		if !rl.allow(ip) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Retry-After", fmt.Sprintf("%d", int(rl.window.Seconds())))
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintf(w, `<p class="text-sm text-red-600">⚠️ Demasiadas solicitudes. Espera %d segundos.</p>`, int(rl.window.Seconds()))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	client, exists := rl.clients[ip]

	if !exists {
		rl.clients[ip] = &clientWindow{
			timestamps: []time.Time{now},
			lastSeen:   now,
		}
		return true
	}

	// Remove expired timestamps
	windowStart := now.Add(-rl.window)
	valid := make([]time.Time, 0, len(client.timestamps))
	for _, ts := range client.timestamps {
		if ts.After(windowStart) {
			valid = append(valid, ts)
		}
	}

	if len(valid) >= rl.limit {
		client.timestamps = valid
		client.lastSeen = now
		return false
	}

	client.timestamps = append(valid, now)
	client.lastSeen = now
	return true
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, client := range rl.clients {
			if now.Sub(client.lastSeen) > rl.cleanup {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For first (reverse proxy)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	return r.RemoteAddr
}
