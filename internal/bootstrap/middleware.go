package bootstrap

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// PanicRecovery is a middleware that recovers from panics, logs the stack trace,
// and returns HTTP 500 with a JSON error body. This prevents Lambda from
// terminating on unhandled panics.
func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := debug.Stack()
				LogJSON(LogEntry{
					Level:     "error",
					RequestID: r.Header.Get("X-Amzn-Trace-Id"),
					Method:    r.Method,
					Path:      r.URL.Path,
					Error:     fmt.Sprintf("panic: %v\n%s", rec, stack),
				})
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
