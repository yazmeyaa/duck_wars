package http_controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/yazmeyaa/duck_wars/internal/models"
	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
	"github.com/yazmeyaa/duck_wars/internal/transport/transport_http"
	"github.com/yazmeyaa/duck_wars/internal/transport/transport_http/middlewares/auth_middleware"
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
	GetState(w http.ResponseWriter, r *http.Request)
}

type characterController struct {
	userService      user_service.UserService
	characterService character_service.CharacterService
	logger           *slog.Logger
}

// GetState implements CharacterController.
func (c *characterController) GetState(w http.ResponseWriter, r *http.Request) {
	userID, err := auth_middleware.GetUserIDFromRequest(r)
	if err != nil {
		c.logger.Error("Error while get user id from request context", "error", err.Error())
		transport_http.WriteInternalError(w)
		return
	}

	c.logger.Info("", "", userID)
}

func writeUserCharacter(w http.ResponseWriter, char *models.Character) {
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

// Get implements CharacterController.
func (c *characterController) Get(w http.ResponseWriter, r *http.Request) {
	userID, err := auth_middleware.GetUserIDFromRequest(r)
	if err != nil {
		c.logger.Error("Error while get user id from request context", "error", err.Error())
		transport_http.WriteInternalError(w)
		return
	}
	char, err := c.characterService.GetByUserID(userID)
	if err != nil {
		transport_http.WriteError(w, http.StatusInternalServerError, "Cannot get user character.")
		return
	}

	writeUserCharacter(w, char)
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

	writeUserCharacter(w, char)
}

func NewCharacterController(userService user_service.UserService, characterService character_service.CharacterService) CharacterController {
	return &characterController{
		userService:      userService,
		characterService: characterService,
	}
}
