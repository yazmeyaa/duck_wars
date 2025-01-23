package http_router

import (
	"net/http"

	"github.com/yazmeyaa/duck_wars/internal/transport/transport_http/http_controllers"
)

func setupAuthRouter(a http_controllers.AuthenticationController) http.Handler {
	authMux := http.NewServeMux()

	authMux.Handle("POST /login/telegram", http.HandlerFunc(a.LoginViaTelegramInitData))
	authMux.Handle("POST /validate", http.HandlerFunc(a.ValidateToken))

	return authMux
}
