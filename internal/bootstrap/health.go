package bootstrap

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// healthHandler creates a health check handler that reports application and
// dependency status. It accepts an optional pgxpool.Pool (nil for SQLite/local mode)
// and the current appEnv string.
func healthHandler(pool *pgxpool.Pool, appEnv string) http.HandlerFunc {
	var coldStart = true
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		resp := map[string]interface{}{
			"status":  "ok",
			"app_env": appEnv,
		}

		// Database check
		if pool != nil {
			if err := pool.Ping(ctx); err != nil {
				resp["database"] = "error"
				resp["status"] = "degraded"
				w.WriteHeader(http.StatusServiceUnavailable)
			} else {
				resp["database"] = "ok"
			}
		} else {
			resp["database"] = "ok" // SQLite mode — always ok
		}

		// Lambda metadata
		if appEnv == "lambda" {
			resp["memory_mb"] = os.Getenv("AWS_LAMBDA_FUNCTION_MEMORY_SIZE")
			resp["cold_start"] = coldStart
			coldStart = false
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
