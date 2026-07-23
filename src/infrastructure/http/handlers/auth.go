package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/application/use_cases"
	"github.com/alexedwards/scs/v2"
)

// AuthHandler handles login, logout, and session management.
type AuthHandler struct {
	authUC   *use_cases.AuthenticateUser
	tmpl     *template.Template
	sessions *scs.SessionManager
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authUC *use_cases.AuthenticateUser, tmpl *template.Template, sessions *scs.SessionManager) *AuthHandler {
	return &AuthHandler{
		authUC:   authUC,
		tmpl:     tmpl,
		sessions: sessions,
	}
}

// LoginPage renders the login form (GET /login).
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	h.renderLogin(w, "")
}

// Login handles PIN authentication (POST /login).
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderLogin(w, "Error al procesar el formulario")
		return
	}

	pin := r.FormValue("pin")
	if pin == "" {
		h.renderLogin(w, "Ingrese su PIN")
		return
	}

	user, err := h.authUC.Execute(r.Context(), pin)
	if err != nil {
		msg := "PIN incorrecto"
		if err == use_cases.ErrAuthAccountLocked {
			msg = "Cuenta bloqueada temporalmente. Intente más tarde."
		}
		h.renderLogin(w, msg)
		return
	}

	// Store user info in session.
	h.sessions.Put(r.Context(), "user_id", strconv.FormatInt(user.ID, 10))
	h.sessions.Put(r.Context(), "user_name", user.Nombre)
	h.sessions.Put(r.Context(), "user_role", user.Rol)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout destroys the session and redirects to login (POST /logout).
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_ = h.sessions.Destroy(r.Context())
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// renderLogin renders the login page with an optional error message.
func (h *AuthHandler) renderLogin(w http.ResponseWriter, errMsg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := map[string]interface{}{
		"PageTitle": "Iniciar Sesión",
		"Error":     errMsg,
	}

	// Try the login template; fall back to inline HTML if not found.
	if t := h.tmpl.Lookup("login.html"); t != nil {
		if err := t.Execute(w, data); err != nil {
			http.Error(w, "Error de template", http.StatusInternalServerError)
		}
		return
	}

	// Fallback: minimal login form when template is not yet created.
	w.WriteHeader(http.StatusOK)
	fallback := `<!DOCTYPE html>
<html lang="es">
<head><meta charset="UTF-8"><title>Login - POS</title></head>
<body style="display:flex;justify-content:center;align-items:center;min-height:100vh;font-family:sans-serif;">
<form method="POST" action="/login" style="text-align:center;">
<h1>POS AI-First</h1>
<p>Ingrese su PIN para acceder</p>
{{if .Error}}<p style="color:red;">{{.Error}}</p>{{end}}
<input type="password" name="pin" maxlength="6" inputmode="numeric" pattern="[0-9]*" placeholder="PIN" autofocus
  style="font-size:2rem;width:200px;text-align:center;padding:10px;margin:10px 0;"/>
<br/><button type="submit" style="padding:10px 30px;font-size:1rem;cursor:pointer;">Entrar</button>
</form>
</body>
</html>`

	// Parse and execute the fallback template with data.
	ft, err := template.New("fallback-login").Parse(fallback)
	if err != nil {
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}
	_ = ft.Execute(w, data)
}
