package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLanguageSwitch(t *testing.T) {
	for _, tc := range []struct {
		name, lang, referer, target string
		status                      int
	}{
		{"local page", "es", "/products?q=tea", "/products?q=tea", 303},
		{"same origin", "en", "http://example.com/login", "/login", 303},
		{"external", "en", "https://evil.test/", "/", 303},
		{"network path", "en", "//evil.test/", "/", 303},
		{"wrong scheme", "en", "https://example.com/login", "/", 303},
		{"invalid language", "xx", "/", "", 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "http://example.com/lang", strings.NewReader("lang="+tc.lang))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Set("Referer", tc.referer)
			w := httptest.NewRecorder()
			NewLangHandler().Switch(w, r)
			if w.Code != tc.status || w.Header().Get("Location") != tc.target {
				t.Fatalf("response = %d %q", w.Code, w.Header().Get("Location"))
			}
			cookies := w.Result().Cookies()
			if tc.status == 400 {
				if len(cookies) != 0 {
					t.Fatal("invalid language must not set cookie")
				}
				return
			}
			if len(cookies) != 1 || cookies[0].Value != tc.lang || !cookies[0].HttpOnly || cookies[0].SameSite != http.SameSiteLaxMode {
				t.Fatalf("cookie = %+v", cookies)
			}
		})
	}
}

func TestLoginPageLanguage(t *testing.T) {
	tmpl := template.Must(template.ParseFiles("../../../../templates/login.html"))
	for _, tc := range []struct{ lang, label string }{{"", "Sign In"}, {"es", "Entrar"}} {
		t.Run(tc.lang, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/login", nil)
			if tc.lang != "" {
				r.AddCookie(&http.Cookie{Name: "lang", Value: tc.lang})
			}
			w := httptest.NewRecorder()
			NewAuthHandler(nil, tmpl, nil).LoginPage(w, r)
			if w.Code != 200 || !strings.Contains(w.Body.String(), tc.label) {
				t.Fatalf("login response %d: %s", w.Code, w.Body.String())
			}
		})
	}
}
