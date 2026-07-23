package http

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
)

// NewSessionManager creates a configured scs.SessionManager backed by SQLite.
// It uses the existing "sessions" table created by migration 002_sessions.sql.
func NewSessionManager(db *sql.DB) *scs.SessionManager {
	sm := scs.New()
	sm.Store = sqlite3store.New(db)
	sm.Lifetime = 8 * time.Hour
	sm.Cookie.Name = "pos_session"
	sm.Cookie.HttpOnly = true
	sm.Cookie.SameSite = http.SameSiteLaxMode
	sm.Cookie.Secure = false // Set to true in production with HTTPS
	return sm
}
