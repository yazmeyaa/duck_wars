package http_router

import (
	"net/http"

	"github.com/yazmeyaa/duck_wars/internal/transport/http/http_controllers"
)

func newCharactersHandler(characterController http_controllers.CharacterController) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /me", characterController.Get)
	router.HandleFunc("GET /{id}", characterController.GetByID)

	return router
}
