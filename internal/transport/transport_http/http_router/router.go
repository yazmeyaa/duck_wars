package http_router

import (
	"log/slog"
	"net/http"

	authorizationservice "github.com/yazmeyaa/ton_game/internal/service/authorization"
	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/http_controllers"
	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/middlewares"
	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/middlewares/cors_middleware"
	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/middlewares/logger_middleware"
	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/middlewares/panic_recover_middleware"
)

type RouterConstructor struct {
	CharacterController      http_controllers.CharacterController
	AuthenticationController http_controllers.AuthenticationController
	ResourcesController      http_controllers.ResourcesController
	Logger                   *slog.Logger
	AuthorizationService     authorizationservice.AuthorizationService
}

func New(c RouterConstructor) http.Handler {
	characterRouter := newCharactersHandler(c.CharacterController, c.AuthorizationService, c.Logger)
	v1 := http.NewServeMux()
	v1.Handle("/characters/", http.StripPrefix("/characters", characterRouter))

	authRouter := setupAuthRouter(c.AuthenticationController)
	v1.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	resources := setupResourcesRouter(c.ResourcesController)
	v1.Handle("/resources/", http.StripPrefix("/resources", resources))

	root := http.NewServeMux()
	root.Handle("/v1/", http.StripPrefix("/v1", v1))
	return middlewares.ApplyMiddleware(
		root,
		logger_middleware.LoggerMiddleware(c.Logger),
		cors_middleware.CorsMiddleware(),
		panic_recover_middleware.PanicRecoveryMiddleware(c.Logger),
	)
}
