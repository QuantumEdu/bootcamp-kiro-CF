package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/application/nlsql"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/i18n"
)

// ChatHandler handles chat-related HTTP requests.
type ChatHandler struct {
	service *nlsql.Service
	tmpl    *template.Template
}

// NewChatHandler creates a new chat handler.
func NewChatHandler(service *nlsql.Service, tmpl *template.Template) *ChatHandler {
	return &ChatHandler{service: service, tmpl: tmpl}
}

// HandleChat processes a chat message and returns an HTMX fragment.
func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.FormValue("query"))
	if query == "" {
		http.Error(w, tr(r, "chat.query_required"), http.StatusBadRequest)
		return
	}

	result := h.service.ProcessQueryLanguage(r.Context(), query, i18n.FromRequest(r).Lang())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := WithUserContext(r, map[string]interface{}{
		"Query":       result.Query,
		"Explanation": result.Explanation,
		"Columns":     result.Columns,
		"Results":     result.Results,
		"Error":       result.Error,
	})
	if err := h.tmpl.ExecuteTemplate(w, "chat_message", data); err != nil {
		http.Error(w, tr(r, "error.template"), http.StatusInternalServerError)
	}
}
