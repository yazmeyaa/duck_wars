package logger_middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/middlewares"
)

func LoggerMiddleware(logger *slog.Logger) middlewares.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			logger.Info("Handled HTTP request", "path", r.URL.Path, "method", r.Method, "time", time.Since(start).Milliseconds())
		})
	}
}
