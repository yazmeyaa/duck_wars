package handlers

import (
	api "github.com/yazmeyaa/ton_game/internal/api/gen"
	"github.com/yazmeyaa/ton_game/internal/game_data"
	characterservice "github.com/yazmeyaa/ton_game/internal/service/character"
	inventoryservice "github.com/yazmeyaa/ton_game/internal/service/inventory"
	items_service "github.com/yazmeyaa/ton_game/internal/service/items"
	user_service "github.com/yazmeyaa/ton_game/internal/service/user"
	"net/http"
)

type ServerHandlers struct {
	userSvc      user_service.UserService
	characterSvc characterservice.CharacterService
	inventorySvc inventoryservice.InventoryService
	itemsSvc     items_service.ItemsService
	gameData     *game_data.GameData
}

var _ api.ServerInterface = (*ServerHandlers)(nil)

func NewServerHandler(
	userSvc user_service.UserService,
	characterSvc characterservice.CharacterService,
	inventorySvc inventoryservice.InventoryService,
	itemsSvc items_service.ItemsService,
	gameData *game_data.GameData,
) *ServerHandlers {
	return &ServerHandlers{
		userSvc:      userSvc,
		characterSvc: characterSvc,
		inventorySvc: inventorySvc,
		itemsSvc:     itemsSvc,
		gameData:     gameData,
	}
}

func (s ServerHandlers) GetV1Adventures(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s ServerHandlers) PostV1AdventuresAdventureId(w http.ResponseWriter, r *http.Request, adventureId string) {
	//TODO implement me
	panic("implement me")
}

func (s ServerHandlers) PostV1AuthLoginTelegram(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s ServerHandlers) GetV1Characters(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s ServerHandlers) GetV1CharactersState(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s ServerHandlers) GetV1CharactersCharacterId(w http.ResponseWriter, r *http.Request, characterId string) {
	//TODO implement me
	panic("implement me")
}

func (s ServerHandlers) GetV1Items(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
