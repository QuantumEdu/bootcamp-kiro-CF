package i18n

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestLanguage(t *testing.T) {
	for _, tc := range []struct{ name, cookie, want string }{
		{"default English", "", "en"}, {"Spanish cookie", "es", "es"},
		{"English cookie", "en", "en"}, {"invalid cookie", "xx", "en"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.cookie != "" {
				r.AddCookie(&http.Cookie{Name: CookieName, Value: tc.cookie})
			}
			if got := FromRequest(r).Lang(); got != tc.want {
				t.Fatalf("language = %q, want %q", got, tc.want)
			}
		})
	}
}
