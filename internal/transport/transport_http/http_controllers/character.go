package http_controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
)

type getCharacterResponse struct {
	ID         int     `json:"id"`
	UserID     int     `json:"userId"`
	Name       string  `json:"name"`
	Coins      float64 `json:"coins"`
	MaxEnergy  uint8   `json:"maxEnergy"`
	Experience float64 `json:"experience"`
}

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\":\"Hello world!\"}"))
}

// GetByID implements CharacterController.
func (c *characterController) GetByID(w http.ResponseWriter, r *http.Request) {
	idValue := r.PathValue("id")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	char, err := c.characterService.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := getCharacterResponse{
		ID:         char.ID,
		UserID:     char.UserID,
		Name:       char.Name,
		MaxEnergy:  char.MaxEnergy,
		Coins:      char.Coins,
		Experience: char.Experience,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func NewCharacterController(userService user_service.UserService, characterService character_service.CharacterService) CharacterController {
	return &characterController{
		userService:      userService,
		characterService: characterService,
	}
}
