package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/http/middleware"
)

// AuthHandler handles login/logout.
type AuthHandler struct {
	auth *middleware.AuthMiddleware
	tmpl *template.Template
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(auth *middleware.AuthMiddleware, tmpl *template.Template) *AuthHandler {
	return &AuthHandler{auth: auth, tmpl: tmpl}
}

// LoginPage renders the login page.
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.tmpl.ExecuteTemplate(w, "login.html", map[string]interface{}{
		"Error": "",
	})
}

// Login processes login form submission.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	pin := r.FormValue("pin")
	if pin == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		h.tmpl.ExecuteTemplate(w, "login.html", map[string]interface{}{
			"Error": "Ingresa tu PIN",
		})
		return
	}

	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}

	token, err := h.auth.Login(pin, clientIP)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		h.tmpl.ExecuteTemplate(w, "login.html", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400, // 24 hours
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout clears the session.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
