package http_controllers

import (
	"net/http"

	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
)

type CharacterController interface {
	Get(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
}

type characterController struct {
	userService      user_service.UserService
	characterService character_service.CharacterService
}

// Get implements CharacterController.
func (c *characterController) Get(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// GetByID implements CharacterController.
func (c *characterController) GetByID(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func NewCharacterController(userService user_service.UserService, characterService character_service.CharacterService) CharacterController {
	return &characterController{
		userService:      userService,
		characterService: characterService,
	}
}
