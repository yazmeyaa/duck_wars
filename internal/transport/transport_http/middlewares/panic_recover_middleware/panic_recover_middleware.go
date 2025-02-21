package panic_recover_middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/middlewares"
)

func PanicRecoveryMiddleware(logger *slog.Logger) middlewares.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					logger.Error("Panic recovered!", slog.Any("error", err))
					jsonBody, _ := json.Marshal(map[string]string{
						"error": "There was an internal server error",
					})
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, err := w.Write(jsonBody)
					if err != nil {
						return
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
