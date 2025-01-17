package http_router

import (
	"net/http"

	"github.com/yazmeyaa/duck_wars/internal/transport/http/http_controllers"
)

type RouterConstructor struct {
	CharacterController http_controllers.CharacterController
}

func New(c RouterConstructor) http.Handler {
	root := http.NewServeMux()

	v1 := http.NewServeMux()
	root.Handle("/v1", v1)

	characterRouter := setupCharacterRouter(c.CharacterController)
	v1.Handle("/character", http.StripPrefix("/character", characterRouter))

	return root
}
