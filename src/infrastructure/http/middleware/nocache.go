package middleware

import "net/http"

// NoCache sets Cache-Control: no-store on responses to prevent HTMX fragments
// from being served stale after navigation or server restart.
func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
