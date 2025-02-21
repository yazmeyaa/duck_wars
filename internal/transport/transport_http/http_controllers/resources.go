package http_controllers

import (
	"log/slog"
	"net/http"

	"github.com/yazmeyaa/ton_game/internal/game_data"
	"github.com/yazmeyaa/ton_game/internal/transport/transport_http"
)

type ResourcesController interface {
	Items(w http.ResponseWriter, r *http.Request)
	Crafts(w http.ResponseWriter, r *http.Request)
	Adventures(w http.ResponseWriter, r *http.Request)
}

type resourcesController struct {
	gd     *game_data.GameData
	logger *slog.Logger
}

// Adventures implements ResourcesController.
func (rc *resourcesController) Adventures(w http.ResponseWriter, r *http.Request) {
	if err := transport_http.WriteJSON(w, rc.gd.Adventures); err != nil {
		transport_http.WriteInternalError(w)
	}
}

// Crafts implements ResourcesController.
func (rc *resourcesController) Crafts(w http.ResponseWriter, r *http.Request) {
	if err := transport_http.WriteJSON(w, rc.gd.Crafts); err != nil {
		transport_http.WriteInternalError(w)
	}
}

// Items implements ResourcesController.
func (rc *resourcesController) Items(w http.ResponseWriter, r *http.Request) {
	if err := transport_http.WriteJSON(w, rc.gd.Items); err != nil {
		transport_http.WriteInternalError(w)
	}
}

func NewResourcesController(gd *game_data.GameData, logger *slog.Logger) ResourcesController {
	return &resourcesController{
		gd:     gd,
		logger: logger,
	}
}
