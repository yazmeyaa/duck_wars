package items_service

import (
	"context"
	"errors"

	"github.com/yazmeyaa/duck_wars/internal/game_data"
	"github.com/yazmeyaa/duck_wars/internal/models"
	inventory_repository "github.com/yazmeyaa/duck_wars/internal/repository/inventory"
	inventory_service "github.com/yazmeyaa/duck_wars/internal/service/inventory"
)

var (
	ErrItemNotFound         = errors.New("item not found")
	ErrNotCraftable         = errors.New("item is not craftable")
	ErrRecipeNotFound       = errors.New("not found reciepe")
	ErrNotEnoughIngredients = errors.New("not enought ingredients")
)

type ItemsService interface {
	Craft(ctx context.Context, characterID int, recipeID string) error
	GetRecipes(itemID string) ([]game_data.CraftRecipe, error)
	Transfer(ctx context.Context, itemID string, quantity uint16, characterID int, targetCharacterID int) error
	Destroy(ctx context.Context, characterID int, itemID string, quantity uint16) error
	Add(ctx context.Context, characterID int, itemID string, quantity uint16) error
	GetInfo(ctx context.Context, characterID int, itemID string) (*models.InventoryItem, error)
}

type itemsService struct {
	inventoryService    inventory_service.InventoryService
	inventoryRepository inventory_repository.InventoryRepository
	gameData            *game_data.GameData
}

// GetRecipes implements ItemsService.
func (i *itemsService) GetRecipes(itemID string) ([]game_data.CraftRecipe, error) {
	panic("unimplemented")
}

// Add implements ItemsService.
func (i *itemsService) Add(ctx context.Context, characterID int, itemID string, quantity uint16) error {
	return i.inventoryService.AddItem(ctx, itemID, characterID, quantity)
}

// Craft implements ItemsService.
func (i *itemsService) Craft(ctx context.Context, characterID int, recipeID string) error {
	reciepe, exist := i.gameData.Crafts[recipeID]
	if !exist {
		return ErrRecipeNotFound
	}

	itemsToRemove := make([]inventory_repository.SingleItemDTO, len(reciepe.Items))
	for idx, item := range reciepe.Items {
		itemsToRemove[idx] = inventory_repository.SingleItemDTO{
			ItemID:   item.ItemID,
			Quantity: item.Quantity,
		}
	}
	itemToCraft := inventory_repository.SingleItemDTO{
		ItemID:   reciepe.ItemID,
		Quantity: uint16(reciepe.Quantity),
	}
	err := i.inventoryRepository.CraftItems(ctx, characterID, itemsToRemove, []inventory_repository.SingleItemDTO{itemToCraft})
	if err != nil {
		return err
	}

	return nil
}

// Destroyц implements ItemsService.
func (i *itemsService) Destroy(ctx context.Context, characterID int, itemID string, quantity uint16) error {
	panic("unimplemented")
}

// GetInfo implements ItemsService.
func (i *itemsService) GetInfo(ctx context.Context, characterID int, itemID string) (*models.InventoryItem, error) {
	panic("unimplemented")
}

// Transfer implements ItemsService.
func (i *itemsService) Transfer(ctx context.Context, itemID string, quantity uint16, characterID int, targetCharacterID int) error {
	panic("unimplemented")
}

func New(inventoryService inventory_service.InventoryService, inventoryRepository inventory_repository.InventoryRepository, gameData *game_data.GameData) ItemsService {
	return &itemsService{
		inventoryService:    inventoryService,
		inventoryRepository: inventoryRepository,
		gameData:            gameData,
	}
}
