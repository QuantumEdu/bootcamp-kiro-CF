package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/application/use_cases"
)

// ClientHandler handles client listing and creation pages.
type ClientHandler struct {
	createUC *use_cases.CreateClient
	listUC   *use_cases.ListClients
	tmpl     *template.Template
}

// NewClientHandler creates a new ClientHandler.
func NewClientHandler(create *use_cases.CreateClient, list *use_cases.ListClients, tmpl *template.Template) *ClientHandler {
	return &ClientHandler{createUC: create, listUC: list, tmpl: tmpl}
}

// List handles GET /clientes — renders the client list page.
func (h *ClientHandler) List(w http.ResponseWriter, r *http.Request) {
	clients, err := h.listUC.Execute(r.Context())
	if err != nil {
		http.Error(w, "Error al cargar clientes", http.StatusInternalServerError)
		return
	}

	data := WithUserContext(r, map[string]interface{}{
		"PageTitle": "Clientes",
		"Clients":   clients,
	})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, "Error de template: "+err.Error(), http.StatusInternalServerError)
	}
}

// CreateForm handles GET /clientes/new — renders the client creation form.
func (h *ClientHandler) CreateForm(w http.ResponseWriter, r *http.Request) {
	data := WithUserContext(r, map[string]interface{}{
		"PageTitle": "Nuevo Cliente",
	})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, "Error de template: "+err.Error(), http.StatusInternalServerError)
	}
}

// Create handles POST /clientes — parses form, validates, and creates a new client.
func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	input := use_cases.CreateClientInput{
		Nombre:    strings.TrimSpace(r.FormValue("nombre")),
		Telefono:  strings.TrimSpace(r.FormValue("telefono")),
		Direccion: strings.TrimSpace(r.FormValue("direccion")),
	}

	_, err := h.createUC.Execute(r.Context(), input)
	if err != nil {
		data := WithUserContext(r, map[string]interface{}{
			"PageTitle": "Nuevo Cliente",
			"Error":     "El nombre del cliente es obligatorio",
			"Input":     input,
		})
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if tmplErr := h.tmpl.ExecuteTemplate(w, "layout.html", data); tmplErr != nil {
			http.Error(w, "Error de template: "+tmplErr.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/clientes", http.StatusSeeOther)
}
