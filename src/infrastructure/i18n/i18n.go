// Package i18n provides a lightweight translation system for the POS AI-First app.
// Supported languages: "es" (Spanish) and "en" (English, default).
// Language is persisted via a "lang" cookie and resolved on every request.
package i18n

import "net/http"

const (
	LangES      = "es"
	LangEN      = "en"
	CookieName  = "lang"
	DefaultLang = LangEN
)

// T is a translation function type used in Go templates: {{call .T "key"}}.
type T func(key string) string

// Translator holds the translations map for a resolved language.
type Translator struct {
	lang         string
	translations map[string]string
}

// New returns a Translator for the given language code.
// Falls back to English for unknown languages.
func New(lang string) *Translator {
	switch lang {
	case LangES:
		return &Translator{lang: LangES, translations: translationsES}
	default:
		return &Translator{lang: LangEN, translations: translationsEN}
	}
}

// FromRequest resolves the language from the request's "lang" cookie.
// Returns a Translator ready to use in templates.
func FromRequest(r *http.Request) *Translator {
	cookie, err := r.Cookie(CookieName)
	if err != nil || cookie.Value == "" {
		return New(DefaultLang)
	}
	return New(cookie.Value)
}

// Translate returns the translation for key in the current language.
// Falls back to the key itself if not found.
func (t *Translator) Translate(key string) string {
	if val, ok := t.translations[key]; ok {
		return val
	}
	return key
}

// Func returns a T function suitable for use in template data as .T.
// Usage in template: {{call .T "nav.dashboard"}}
func (t *Translator) Func() T {
	return t.Translate
}

// Lang returns the current language code ("es" or "en").
func (t *Translator) Lang() string {
	return t.lang
}

// IsEN returns true if the current language is English.
func (t *Translator) IsEN() bool {
	return t.lang == LangEN
}
