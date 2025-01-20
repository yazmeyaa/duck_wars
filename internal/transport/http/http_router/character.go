package http_router

import (
	"log/slog"
	"net/http"

	authorization_service "github.com/yazmeyaa/duck_wars/internal/service/authorization"
	"github.com/yazmeyaa/duck_wars/internal/transport/http/http_controllers"
	"github.com/yazmeyaa/duck_wars/internal/transport/http/middlewares"
	"github.com/yazmeyaa/duck_wars/internal/transport/http/middlewares/auth_middleware"
)

func newCharactersHandler(characterController http_controllers.CharacterController, authorizationService authorization_service.AuthorizationService, logger *slog.Logger) http.Handler {
	router := http.NewServeMux()

	router.Handle("GET /me", middlewares.ApplyMiddleware(http.HandlerFunc(characterController.Get), auth_middleware.AuthMiddleware(logger, authorizationService)))
	router.HandleFunc("GET /{id}", characterController.GetByID)

	return router
}
