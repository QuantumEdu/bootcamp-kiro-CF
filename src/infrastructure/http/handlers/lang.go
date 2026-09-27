package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/i18n"
)

// LangHandler handles language switching via POST /lang.
// Sets a "lang" cookie and redirects back to the referring page.
type LangHandler struct{}

// NewLangHandler creates a new LangHandler.
func NewLangHandler() *LangHandler {
	return &LangHandler{}
}

// Switch sets the lang cookie and redirects back to the referer.
// POST /lang with form field lang=en or lang=es.
func (h *LangHandler) Switch(w http.ResponseWriter, r *http.Request) {
	lang := r.FormValue("lang")
	if lang != i18n.LangEN && lang != i18n.LangES {
		http.Error(w, "Unsupported language", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     i18n.CookieName,
		Value:    lang,
		Path:     "/",
		MaxAge:   365 * 24 * 3600, // 1 year
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect back to referring page, or home if no referer
	referer := localReferer(r)
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

// localReferer allows only root-relative paths or the request's own origin.
func localReferer(r *http.Request) string {
	u, err := url.Parse(r.Header.Get("Referer"))
	if err != nil || u.User != nil || strings.Contains(u.Path, "\\") {
		return "/"
	}
	if u.IsAbs() {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		if u.Scheme != scheme || !strings.EqualFold(u.Host, r.Host) {
			return "/"
		}
	} else if u.Host != "" {
		return "/"
	}
	if !strings.HasPrefix(u.Path, "/") || strings.HasPrefix(u.Path, "//") {
		return "/"
	}
	return u.RequestURI()
}
