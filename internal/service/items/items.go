package items_service

import (
	"context"
	"errors"

	"github.com/yazmeyaa/ton_game/internal/game_data"
	"github.com/yazmeyaa/ton_game/internal/models"
	inventory_repository "github.com/yazmeyaa/ton_game/internal/repository/inventory"
	inventory_service "github.com/yazmeyaa/ton_game/internal/service/inventory"
)

var (
	ErrItemNotFound         = errors.New("item not found")
	ErrRecipeNotFound       = errors.New("not found reciepe")
	ErrNotEnoughIngredients = errors.New("not enough ingredients")
)

type ItemsService interface {
	Craft(ctx context.Context, characterID int, recipeID string) error
	GetRecipes(itemID string) []game_data.CraftRecipe
	Transfer(ctx context.Context, item inventory_repository.SingleItemDTO, characterID int, targetCharacterID int) error
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
func (i *itemsService) GetRecipes(itemID string) []game_data.CraftRecipe {
	var reciepes []game_data.CraftRecipe

	for _, reciepe := range i.gameData.Crafts {
		if reciepe.ItemID == itemID {
			reciepes = append(reciepes, reciepe)
		}
	}

	return reciepes
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
func (i *itemsService) Transfer(ctx context.Context, item inventory_repository.SingleItemDTO, characterID int, targetCharacterID int) error {
	return i.inventoryRepository.TransferItems(ctx, characterID, targetCharacterID, []inventory_repository.SingleItemDTO{item})
}

func New(inventoryService inventory_service.InventoryService, inventoryRepository inventory_repository.InventoryRepository, gameData *game_data.GameData) ItemsService {
	return &itemsService{
		inventoryService:    inventoryService,
		inventoryRepository: inventoryRepository,
		gameData:            gameData,
	}
}
