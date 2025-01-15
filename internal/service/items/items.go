package items_service

import (
	"context"
	"errors"

	"github.com/yazmeyaa/duck_wars/internal/game_data"
	"github.com/yazmeyaa/duck_wars/internal/models"
	inventory_service "github.com/yazmeyaa/duck_wars/internal/service/inventory"
)

var (
	ErrItemNotFound         = errors.New("item not found")
	ErrNotCraftable         = errors.New("item is not craftable")
	ErrRecipeNotFound       = errors.New("not found reciepe")
	ErrNotEnoughIngredients = errors.New("not enought ingredients")
)

type ItemsService interface {
	Craft(ctx context.Context, characterID int, itemID string, recipeID string, quantity uint16) error
	GetRecipes(itemID string) ([]game_data.CraftRecipe, error)
	Transfer(ctx context.Context, itemID string, quantity uint16, characterID int, targetCharacterID int) error
	Destroy(ctx context.Context, characterID int, itemID string, quantity uint16) error
	Add(ctx context.Context, characterID int, itemID string, quantity uint16) error
	GetInfo(ctx context.Context, characterID int, itemID string) (*models.InventoryItem, error)
}

type itemsService struct {
	inventoryService inventory_service.InventoryService
	gameData         *game_data.GameData
}

// GetRecipes implements ItemsService.
func (i *itemsService) GetRecipes(itemID string) ([]game_data.CraftRecipe, error) {
	panic("unimplemented")
}

// Add implements ItemsService.
func (i *itemsService) Add(ctx context.Context, characterID int, itemID string, quantity uint16) error {
	return i.inventoryService.AddItem(ctx, itemID, characterID, quantity)
}

type userCraftResource struct {
	itemID   string
	required uint16
}

// Craft implements ItemsService.
func (i *itemsService) Craft(ctx context.Context, characterID int, itemID string, recipeID string, quantity uint16) error {
	item, itemExist := i.gameData.Items[itemID]
	if !itemExist {
		return ErrItemNotFound
	}
	if !item.IsCraftable {
		return ErrNotCraftable
	}

	itemRecipes, itemRecipesExist := i.gameData.Crafts[itemID]
	if !itemRecipesExist {
		return ErrRecipeNotFound
	}

	recipe, recipeExist := itemRecipes.Recipes[recipeID]
	if !recipeExist {
		return ErrRecipeNotFound
	}

	requiredResources := make(map[string]userCraftResource)
	for _, recipeItem := range recipe.Items {
		requiredResources[recipeItem.ItemID] = userCraftResource{
			itemID:   recipeItem.ItemID,
			required: uint16(quantity) * recipeItem.Quantity,
		}
	}

	userInventory, err := i.inventoryService.GetInventory(ctx, characterID)
	if err != nil {
		return err
	}

	var itemsToUpdate []inventory_service.ItemDTO
	for _, item := range userInventory {
		if resource, exists := requiredResources[item.ItemID]; exists {
			if resource.required > item.Quantity {
				return ErrNotEnoughIngredients
			}
			itemsToUpdate = append(itemsToUpdate, inventory_service.ItemDTO{
				ItemID:   item.ItemID,
				Quantity: item.Quantity - resource.required,
			})
			delete(requiredResources, item.ItemID)
		}
	}

	if len(requiredResources) > 0 {
		return ErrNotEnoughIngredients
	}

	if err := i.inventoryService.UpdateItemsBatch(ctx, characterID, itemsToUpdate); err != nil {
		return err
	}

	if err := i.inventoryService.AddItem(ctx, itemID, characterID, quantity); err != nil {
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

func New(inventoryService inventory_service.InventoryService, gameData *game_data.GameData) ItemsService {
	return &itemsService{
		inventoryService: inventoryService,
		gameData:         gameData,
	}
}
