package items_service

import (
	"context"
	"errors"

	"github.com/yazmeyaa/ton_game/internal/game_data"
	"github.com/yazmeyaa/ton_game/internal/models"
	inventoryrepository "github.com/yazmeyaa/ton_game/internal/repository/inventory"
	inventoryservice "github.com/yazmeyaa/ton_game/internal/service/inventory"
)

var (
	ErrItemNotFound         = errors.New("item not found")
	ErrRecipeNotFound       = errors.New("not found recipe")
	ErrNotEnoughIngredients = errors.New("not enough ingredients")
)

type ItemsService interface {
	Craft(ctx context.Context, characterID int, recipeID string) error
	GetRecipes(itemID string) []game_data.CraftRecipe
	Transfer(ctx context.Context, item inventoryrepository.SingleItemDTO, characterID int, targetCharacterID int) error
	Destroy(ctx context.Context, characterID int, itemID string, quantity uint16) error
	Add(ctx context.Context, characterID int, itemID string, quantity uint16) error
	GetInfo(ctx context.Context, characterID int, itemID string) (*models.InventoryItem, error)
}

type itemsService struct {
	inventoryService    inventoryservice.InventoryService
	inventoryRepository inventoryrepository.InventoryRepository
	gameData            *game_data.GameData
}

// GetRecipes implements ItemsService.
func (i *itemsService) GetRecipes(itemID string) []game_data.CraftRecipe {
	var recipes []game_data.CraftRecipe

	for _, recipe := range i.gameData.Crafts {
		if recipe.ItemID == itemID {
			recipes = append(recipes, recipe)
		}
	}

	return recipes
}

// Add implements ItemsService.
func (i *itemsService) Add(ctx context.Context, characterID int, itemID string, quantity uint16) error {
	return i.inventoryService.AddItem(ctx, itemID, characterID, quantity)
}

// Craft implements ItemsService.
func (i *itemsService) Craft(ctx context.Context, characterID int, recipeID string) error {
	recipe, exist := i.gameData.Crafts[recipeID]
	if !exist {
		return ErrRecipeNotFound
	}

	itemsToRemove := make([]inventoryrepository.SingleItemDTO, len(recipe.Items))
	for idx, item := range recipe.Items {
		itemsToRemove[idx] = inventoryrepository.SingleItemDTO{
			ItemID:   item.ItemID,
			Quantity: item.Quantity,
		}
	}
	itemToCraft := inventoryrepository.SingleItemDTO{
		ItemID:   recipe.ItemID,
		Quantity: uint16(recipe.Quantity),
	}
	err := i.inventoryRepository.CraftItems(ctx, characterID, itemsToRemove, []inventoryrepository.SingleItemDTO{itemToCraft})
	if err != nil {
		return err
	}

	return nil
}

// Destroy implements ItemsService.
func (i *itemsService) Destroy(ctx context.Context, characterID int, itemID string, quantity uint16) error {
	// TODO: IMPLEMENT THIS
	panic("unimplemented")
}

// GetInfo implements ItemsService.
func (i *itemsService) GetInfo(ctx context.Context, characterID int, itemID string) (*models.InventoryItem, error) {
	return i.inventoryRepository.GetItem(ctx, characterID, itemID)
}

// Transfer implements ItemsService.
func (i *itemsService) Transfer(ctx context.Context, item inventoryrepository.SingleItemDTO, characterID int, targetCharacterID int) error {
	return i.inventoryRepository.TransferItems(ctx, characterID, targetCharacterID, []inventoryrepository.SingleItemDTO{item})
}

func New(inventoryService inventoryservice.InventoryService, inventoryRepository inventoryrepository.InventoryRepository, gameData *game_data.GameData) ItemsService {
	return &itemsService{
		inventoryService:    inventoryService,
		inventoryRepository: inventoryRepository,
		gameData:            gameData,
	}
}
