package inventory_service

import (
	"context"
	"log/slog"

	"github.com/yazmeyaa/duck_wars/internal/models"
	inventory_repository "github.com/yazmeyaa/duck_wars/internal/repository/inventory"
)

type ItemDTO struct {
	ItemID   string
	Quantity uint16
}

type InventoryService interface {
	AddItem(ctx context.Context, itemID string, characterID int, quantity uint16) error
	AddItemsBatch(ctx context.Context, characterID int, items []ItemDTO) error
	RemoveItem(ctx context.Context, itemID string, characterID int) error
	RemoveItemsBatch(ctx context.Context, characterID int, itemsToDelete []ItemDTO) error
	GetInventory(ctx context.Context, characterID int) ([]models.InventoryItem, error)
	GetItemByID(ctx context.Context, characterID int, itemID string) (*models.InventoryItem, error)
	UpdateItemsBatch(ctx context.Context, characterID int, items []ItemDTO) error
}

type inventoryService struct {
	repo   inventory_repository.InventoryRepository
	logger *slog.Logger
}

// UpdateItemsBatch implements InventoryService.
func (i *inventoryService) UpdateItemsBatch(ctx context.Context, characterID int, items []ItemDTO) error {
	s := make([]inventory_repository.SingleItemDTO, len(items))
	for _, item := range items {
		s = append(s, inventory_repository.SingleItemDTO{
			ItemID:   item.ItemID,
			Quantity: item.Quantity,
		})
	}

	return i.repo.UpdateItemsBatch(ctx, characterID, s)
}

func New(repo inventory_repository.InventoryRepository, logger *slog.Logger) InventoryService {
	return &inventoryService{
		repo:   repo,
		logger: logger,
	}
}

// AddItemsBatch implements InventoryService.
func (i *inventoryService) AddItemsBatch(ctx context.Context, characterID int, items []ItemDTO) error {
	panic("unimplemented")
}

// GetItemByID implements InventoryService.
func (i *inventoryService) GetItemByID(ctx context.Context, characterID int, itemID string) (*models.InventoryItem, error) {
	panic("unimplemented")
}

// RemoveItemsBatch implements InventoryService.
func (i *inventoryService) RemoveItemsBatch(ctx context.Context, characterID int, itemsToDelete []ItemDTO) error {
	panic("unimplemented")
}

func (i *inventoryService) AddItem(ctx context.Context, itemID string, characterID int, quantity uint16) error {
	err := i.repo.AddItem(ctx, itemID, characterID, quantity)

	if err != nil {
		i.logger.Error("Error while add item to character", "error", err.Error(), "character_id", characterID, "item_id", itemID)
	}

	return err
}

func (i *inventoryService) RemoveItem(ctx context.Context, itemID string, characterID int) error {
	return i.repo.RemoveItem(ctx, itemID, characterID)
}

func (i *inventoryService) GetInventory(ctx context.Context, characterID int) ([]models.InventoryItem, error) {
	return i.repo.GetInventory(ctx, characterID)
}
