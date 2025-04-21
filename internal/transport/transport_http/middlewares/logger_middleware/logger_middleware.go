package logger_middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/middlewares"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	logger     *slog.Logger
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(logger *slog.Logger) middlewares.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrappedWriter := &responseWriterWrapper{
				ResponseWriter: w,
				logger:         logger,
				statusCode:     0,
			}
			start := time.Now()
			h.ServeHTTP(wrappedWriter, r)
			logger.Info(
				"Handled HTTP request",
				"path", r.URL.Path,
				"method", r.Method,
				"time", time.Since(start).Milliseconds(),
				"statusCode", wrappedWriter.statusCode,
			)
		})
	}
}
