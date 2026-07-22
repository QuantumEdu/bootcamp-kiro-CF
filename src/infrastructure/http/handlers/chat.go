package handlers

import (
	"html/template"
	"net/http"
	"strings"

	usecases "github.com/QuantumEdu/bootcamp-kiro-CF/src/application/use_cases"
)

// ChatHandler handles chat-related HTTP requests.
type ChatHandler struct {
	chatService *usecases.ChatService
	tmpl        *template.Template
}

// NewChatHandler creates a new chat handler.
func NewChatHandler(chatService *usecases.ChatService, tmpl *template.Template) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		tmpl:        tmpl,
	}
}

// HandleChat processes a chat message and returns an HTMX fragment.
func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := strings.TrimSpace(r.FormValue("query"))
	if query == "" {
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}

	result := h.chatService.ProcessQuery(r.Context(), query)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := map[string]interface{}{
		"Query":       result.Query,
		"SQL":         result.SQL,
		"Explanation": result.Explanation,
		"Columns":     result.Columns,
		"Results":     result.Results,
		"Error":       result.Error,
	}

	if err := h.tmpl.ExecuteTemplate(w, "chat/message.html", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
