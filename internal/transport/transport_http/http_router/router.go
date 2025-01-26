package http_router

import (
	"log/slog"
	"net/http"

	authorization_service "github.com/yazmeyaa/duck_wars/internal/service/authorization"
	"github.com/yazmeyaa/duck_wars/internal/transport/transport_http/http_controllers"
	"github.com/yazmeyaa/duck_wars/internal/transport/transport_http/middlewares"
	"github.com/yazmeyaa/duck_wars/internal/transport/transport_http/middlewares/cors_middleware"
	"github.com/yazmeyaa/duck_wars/internal/transport/transport_http/middlewares/logger_middleware"
	"github.com/yazmeyaa/duck_wars/internal/transport/transport_http/middlewares/panic_recover_middleware"
)

type RouterConstructor struct {
	CharacterController      http_controllers.CharacterController
	AuthenticationController http_controllers.AuthenticationController
	Logger                   *slog.Logger
	AuthorizationService     authorization_service.AuthorizationService
}

func New(c RouterConstructor) http.Handler {
	characterRouter := newCharactersHandler(c.CharacterController, c.AuthorizationService, c.Logger)
	v1 := http.NewServeMux()
	v1.Handle("/characters/", http.StripPrefix("/characters", characterRouter))

	authRouter := setupAuthRouter(c.AuthenticationController)
	v1.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	root := http.NewServeMux()
	root.Handle("/v1/", http.StripPrefix("/v1", v1))
	return middlewares.ApplyMiddleware(
		root,
		logger_middleware.LoggerMiddleware(c.Logger),
		cors_middleware.CorsMiddleware(),
		panic_recover_middleware.PanicRecoveryMiddleware(c.Logger),
	)
}
