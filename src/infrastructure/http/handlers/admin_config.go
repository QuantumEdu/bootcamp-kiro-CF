package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/application/services"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
)

// AdminConfigHandler handles the admin configuration page for managing API keys.
type AdminConfigHandler struct {
	configRepo ports.ConfigRepository
	crypto     *services.CryptoService
	tmpl       *template.Template
}

// NewAdminConfigHandler creates a new AdminConfigHandler.
func NewAdminConfigHandler(repo ports.ConfigRepository, crypto *services.CryptoService, tmpl *template.Template) *AdminConfigHandler {
	return &AdminConfigHandler{configRepo: repo, crypto: crypto, tmpl: tmpl}
}

// Show handles GET /admin/config — reads encrypted key, decrypts, masks, renders template.
func (h *AdminConfigHandler) Show(w http.ResponseWriter, r *http.Request) {
	encrypted, _ := h.configRepo.Get(r.Context(), "openrouter_api_key")
	masked := ""
	if encrypted != "" {
		decrypted, err := h.crypto.Decrypt(encrypted)
		if err == nil {
			masked = ports.MaskAPIKey(decrypted)
		}
	}

	data := WithUserContext(r, map[string]interface{}{
		"PageTitle": "Configuración",
		"MaskedKey": masked,
		"HasKey":    encrypted != "",
	})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.tmpl.ExecuteTemplate(w, "layout.html", data)
}

// Update handles POST /admin/config — validates non-empty key, encrypts, stores, redirects.
func (h *AdminConfigHandler) Update(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	apiKey := strings.TrimSpace(r.FormValue("api_key"))
	if apiKey == "" {
		data := WithUserContext(r, map[string]interface{}{
			"PageTitle": "Configuración",
			"Error":     "La API key no puede estar vacía",
		})
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.tmpl.ExecuteTemplate(w, "layout.html", data)
		return
	}

	encrypted, err := h.crypto.Encrypt(apiKey)
	if err != nil {
		http.Error(w, "Error al cifrar la clave", http.StatusInternalServerError)
		return
	}

	if err := h.configRepo.Set(r.Context(), "openrouter_api_key", encrypted); err != nil {
		http.Error(w, "Error al guardar la configuración", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/config", http.StatusSeeOther)
}
