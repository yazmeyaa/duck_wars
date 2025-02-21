package http_router

import (
	"log/slog"
	"net/http"

	authorization_service "github.com/yazmeyaa/ton_game/internal/service/authorization"
	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/http_controllers"
	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/middlewares"
	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/middlewares/auth_middleware"
)

func newCharactersHandler(characterController http_controllers.CharacterController, authorizationService authorization_service.AuthorizationService, logger *slog.Logger) http.Handler {
	router := http.NewServeMux()

	router.Handle("GET /me", middlewares.ApplyMiddleware(
		http.HandlerFunc(characterController.Get),
		auth_middleware.AuthMiddleware(logger, authorizationService)),
	)
	router.Handle("GET /state", middlewares.ApplyMiddleware(
		http.HandlerFunc(characterController.GetState),
		auth_middleware.AuthMiddleware(logger, authorizationService)),
	)
	router.HandleFunc("GET /{id}", characterController.GetByID)

	return router
}
