package handlers

import (
	"html/template"
	"net/http"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/database"
)

// ClientHandler handles client-related HTTP requests.
type ClientHandler struct {
	repo *database.ClientRepo
	tmpl *template.Template
}

// NewClientHandler creates a new client handler.
func NewClientHandler(repo *database.ClientRepo, tmpl *template.Template) *ClientHandler {
	return &ClientHandler{repo: repo, tmpl: tmpl}
}

// List renders the clients page.
func (h *ClientHandler) List(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("q")
	clients, err := h.repo.List(search)
	if err != nil {
		http.Error(w, "Error loading clients", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"PageTitle": "Clientes",
		"Clients":  clients,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Header.Get("HX-Request") == "true" {
		h.tmpl.ExecuteTemplate(w, "clients/index.html", data)
		return
	}
	h.tmpl.ExecuteTemplate(w, "layout.html", data)
}

// Create handles client creation.
func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	client := &entities.Client{
		Nombre:    r.FormValue("nombre"),
		Telefono:  r.FormValue("telefono"),
		Direccion: r.FormValue("direccion"),
	}

	if client.Nombre == "" {
		http.Error(w, "Nombre es requerido", http.StatusBadRequest)
		return
	}

	if _, err := h.repo.Create(client); err != nil {
		http.Error(w, "Error creating client", http.StatusInternalServerError)
		return
	}

	// Re-render the list
	h.List(w, r)
}
