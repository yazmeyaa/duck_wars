package auth_middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	authorization_service "github.com/yazmeyaa/duck_wars/internal/service/authorization"
	"github.com/yazmeyaa/duck_wars/internal/transport/transport_http"
	"github.com/yazmeyaa/duck_wars/internal/transport/transport_http/middlewares"
)

type RequestContextKey string

const (
	userIDContextKey RequestContextKey = "user_id"
)

var (
	ErrUserIDNotInjected = errors.New("cannot get user id from request context; make sure that you use GetUserIDFromRequest only with AuthMiddleware middleware")
	ErrParsingUserID     = errors.New("cannot parse user id from token")
)

func GetUserIDFromRequest(r *http.Request) (int, error) {
	v := r.Context().Value(userIDContextKey)
	if userID, ok := v.(string); ok {
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			return 0, ErrParsingUserID
		}
		return userIDInt, nil
	}

	return 0, ErrUserIDNotInjected
}

func AuthMiddleware(logger *slog.Logger, a authorization_service.AuthorizationService) middlewares.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				transport_http.WriteError(w, http.StatusBadRequest, "Missed authorization data.")
				return
			}

			claims, err := a.ParseJWTClaims(auth)
			if err != nil {
				logger.Error("Error while validation request.", "error", err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			sub, err := claims.GetSubject()
			if err != nil {
				logger.Error("Error while get issuer from JWT claims...", "error", err.Error())
				transport_http.WriteError(w, http.StatusInternalServerError, "Internal error, please try again later...")
				return
			}

			logger.Info("Subject claim", "sub", sub)

			ctx := context.WithValue(r.Context(), userIDContextKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
