package auth_middleware

import (
	"log/slog"
	"net/http"

	authorization_service "github.com/yazmeyaa/duck_wars/internal/service/authorization"
	"github.com/yazmeyaa/duck_wars/internal/transport/http/middlewares"
)

func AuthMiddleware(logger *slog.Logger, a authorization_service.AuthorizationService) middlewares.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "Missed authorization data.", http.StatusBadRequest)
				return
			}

			err := a.ValidateTelegramData(auth)
			if err != nil {
				logger.Error("Error while validation request.", "error", err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
