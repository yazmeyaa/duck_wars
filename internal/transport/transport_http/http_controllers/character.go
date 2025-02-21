package http_controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/yazmeyaa/ton_game/internal/models"
	character_service "github.com/yazmeyaa/ton_game/internal/service/character"
	user_service "github.com/yazmeyaa/ton_game/internal/service/user"
	"github.com/yazmeyaa/ton_game/internal/transport/transport_http"
	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/middlewares/auth_middleware"
)

type getCharacterResponse struct {
	ID         int     `json:"id"`
	UserID     int     `json:"userId"`
	Name       string  `json:"name"`
	Coins      float64 `json:"coins"`
	MaxEnergy  uint8   `json:"maxEnergy"`
	Experience float64 `json:"experience"`
}

type getCharacterStateResponse struct {
	Energy        uint8     `json:"energy"`
	SpendEnergyAt time.Time `json:"spend_energy_at"`
	State         string    `json:"state"`
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	char, err := c.characterService.GetByUserID(userID)
	if err != nil {
		c.logger.Error("Cannot get character", "error", err.Error())
		transport_http.WriteError(w, http.StatusNotFound, "Cannot find character linked to user.")
		return
	}

	state, err := c.userService.GetState(ctx, char.ID)
	if err != nil {
		c.logger.Error("Cannot get user state", "error", err.Error())
		transport_http.WriteInternalError(w)
		return
	}

	writeUserState(w, state)
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

func writeUserState(w http.ResponseWriter, state *models.CharacterState) {
	response := getCharacterStateResponse{
		Energy:        state.Energy,
		SpendEnergyAt: state.SpendEnergyAt,
		State:         string(state.State),
	}

	if err := transport_http.WriteJSON(w, &response); err != nil {
		transport_http.WriteInternalError(w)
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

func NewCharacterController(userService user_service.UserService, characterService character_service.CharacterService, logger *slog.Logger) CharacterController {
	return &characterController{
		userService:      userService,
		characterService: characterService,
		logger:           logger,
	}
}
