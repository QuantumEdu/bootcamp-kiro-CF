package bootstrap

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestRouterHTTPSCookiePolicy(t *testing.T) {
	oldTemplates, oldStatic := TemplatesDir, StaticDir
	TemplatesDir, StaticDir = "../../templates", "../../static"
	t.Cleanup(func() { TemplatesDir, StaticDir = oldTemplates, oldStatic })
	for _, secure := range []bool{false, true} {
		name := "local HTTP"
		if secure {
			name = "HTTPS deployment behind proxy"
		}
		t.Run(name, func(t *testing.T) {
			router, cleanup, err := BuildRouter(Config{DatabasePath: filepath.Join(t.TempDir(), "demo.db"), SessionSecret: "test-only-secret", SessionCookieSecure: secure, QueryTimeoutSeconds: 5})
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()
			for _, form := range []struct{ path, body, cookie string }{{"/login", "pin=1234", "pos_session"}, {"/lang", "lang=es", "lang"}} {
				req := httptest.NewRequest(http.MethodPost, "http://demo.example"+form.path, strings.NewReader(form.body))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.Header.Set("X-Forwarded-Proto", "https")
				scheme := "http"
				if secure {
					scheme = "https"
				}
				req.Header.Set("Referer", scheme+"://demo.example/clientes")
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				if rec.Code != http.StatusSeeOther {
					t.Fatalf("%s response = %d", form.path, rec.Code)
				}
				if form.path == "/lang" && rec.Header().Get("Location") != "/clientes" {
					t.Errorf("language redirect = %q, want same-origin clients page", rec.Header().Get("Location"))
				}
				found := false
				for _, cookie := range rec.Result().Cookies() {
					if cookie.Name != form.cookie {
						continue
					}
					found = true
					if cookie.Secure != secure || !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode {
						t.Errorf("%s cookie policy = %+v, secure want %v", form.cookie, cookie, secure)
					}
				}
				if !found {
					t.Errorf("missing %s cookie", form.cookie)
				}
			}
		})
	}
}
