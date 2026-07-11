package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/appcontext"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// RequestLogger logs the beginning and end of each request
// It also sets a request-scoped logger on the context
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := chiMiddleware.GetReqID(r.Context())

			requestLogger := logger.With(
				"request_id", requestID,
				"method", r.Method,
				"path", r.URL.Path,
			)
			ctx := appcontext.SetLogger(r.Context(), requestLogger)

			start := time.Now()
			next.ServeHTTP(w, r.WithContext(ctx))
			requestLogger.Info("request completed", "duration_ms", time.Since(start).Milliseconds())
		})
	}
}
