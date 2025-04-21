package http_server

import (
	api "github.com/yazmeyaa/ton_game/internal/api/gen"
	"github.com/yazmeyaa/ton_game/internal/game_data"
	"github.com/yazmeyaa/ton_game/internal/http_server/handlers"
	characterservice "github.com/yazmeyaa/ton_game/internal/service/character"
	inventoryservice "github.com/yazmeyaa/ton_game/internal/service/inventory"
	itemsservice "github.com/yazmeyaa/ton_game/internal/service/items"
	userservice "github.com/yazmeyaa/ton_game/internal/service/user"
	"net/http"
)

type NewServerOptions struct {
	userSvc      userservice.UserService
	characterSvc characterservice.CharacterService
	inventorySvc inventoryservice.InventoryService
	itemsSvc     itemsservice.ItemsService
	gameData     *game_data.GameData
}

func NewServer(opts NewServerOptions) http.Handler {
	sh := handlers.NewServerHandler(opts.userSvc, opts.characterSvc, opts.inventorySvc, opts.itemsSvc, opts.gameData)

	server := api.HandlerWithOptions(sh, api.StdHTTPServerOptions{
		Middlewares: []api.MiddlewareFunc{},
	})
	return server
}
