package http_controllers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/yazmeyaa/ton_game/internal/game_data"
	"github.com/yazmeyaa/ton_game/internal/transport/transport_http"
	"log/slog"
	"net/http"
	"time"
)

type ResourcesController interface {
	Items(w http.ResponseWriter, r *http.Request)
	Crafts(w http.ResponseWriter, r *http.Request)
	Adventures(w http.ResponseWriter, r *http.Request)
	All(w http.ResponseWriter, r *http.Request)
}

type resourcesController struct {
	gd     *game_data.GameData
	rdb    *redis.Client
	logger *slog.Logger
}

type adventureResponse struct {
	AdventureID string `json:"adventureId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	EnergyCost  int    `json:"energyCost"`
}

func adventureToResponse(adventure game_data.Adventure) adventureResponse {
	return adventureResponse{
		AdventureID: adventure.AdventureID,
		Title:       adventure.Title,
		Description: adventure.Description,
		Duration:    adventure.Duration,
		EnergyCost:  adventure.EnergyCost,
	}
}

func adventuresMapToResponseMap(gd *game_data.GameData) map[string]adventureResponse {
	adventures := make(map[string]adventureResponse, len(gd.Adventures))

	for _, adventure := range gd.Adventures {
		adventures[adventure.AdventureID] = adventureToResponse(adventure)
	}

	return adventures
}

const (
	GetAllResourcesCacheKey = "resources.all"
)

func (rc *resourcesController) All(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	cacheRes, err := rc.rdb.Get(ctx, GetAllResourcesCacheKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			adventures := adventuresMapToResponseMap(rc.gd)

			res := map[string]interface{}{
				"adventures":       adventures,
				"crafts":           rc.gd.Crafts,
				"adventureRewards": rc.gd.AdventureRewards,
				"items":            rc.gd.Items,
				"weaponAttributes": rc.gd.WeaponAttributes,
			}

			jsonBytes, err := json.Marshal(res)
			if err != nil {
				err := transport_http.WriteInternalError(w)
				if err != nil {
					return
				}
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(jsonBytes)
			if err != nil {
				rc.logger.Error("Error while writing response", "error", err.Error())
				return
			}

			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
				defer cancel()
				if err := rc.rdb.Set(ctx, GetAllResourcesCacheKey, string(jsonBytes), 0).Err(); err != nil {
					rc.logger.Error("Cannot write data to cache...", "error", err.Error())
					return
				}
			}()
			return
		}

		rc.logger.Error("Failed to get cached data", "cache_key", GetAllResourcesCacheKey, "error", err.Error())
		err := transport_http.WriteInternalError(w)
		if err != nil {
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(cacheRes))
	if err != nil {
		rc.logger.Error("Error while write cached data.", "error", err.Error())
		err := transport_http.WriteInternalError(w)
		if err != nil {
			return
		}
		return
	}
}

// Adventures implements ResourcesController.
func (rc *resourcesController) Adventures(w http.ResponseWriter, r *http.Request) {
	if err := transport_http.WriteJSON(w, rc.gd.Adventures); err != nil {
		err := transport_http.WriteInternalError(w)
		if err != nil {
			return
		}
	}
}

// Crafts implements ResourcesController.
func (rc *resourcesController) Crafts(w http.ResponseWriter, r *http.Request) {
	if err := transport_http.WriteJSON(w, rc.gd.Crafts); err != nil {
		err := transport_http.WriteInternalError(w)
		if err != nil {
			return
		}
	}
}

// Items implements ResourcesController.
func (rc *resourcesController) Items(w http.ResponseWriter, r *http.Request) {
	if err := transport_http.WriteJSON(w, rc.gd.Items); err != nil {
		err := transport_http.WriteInternalError(w)
		if err != nil {
			return
		}
	}
}

func NewResourcesController(gd *game_data.GameData, rdb *redis.Client, logger *slog.Logger) ResourcesController {
	return &resourcesController{
		gd:     gd,
		rdb:    rdb,
		logger: logger,
	}
}
