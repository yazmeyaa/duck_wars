package http_router

import (
	"net/http"

	"github.com/yazmeyaa/duck_wars/internal/transport/http/http_controllers"
)

func setupCharacterRouter(characterController http_controllers.CharacterController) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /my", characterController.Get)
	router.HandleFunc("GET /{id}", characterController.GetByID)

	return router
}
