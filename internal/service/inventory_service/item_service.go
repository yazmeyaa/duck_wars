package inventory_service

import (
	"context"
	"log/slog"

	"github.com/yazmeyaa/duck_wars/internal/models"
	inventory_repository "github.com/yazmeyaa/duck_wars/internal/repository/inventory"
)

type InventoryService struct {
	repo   *inventory_repository.InventoryRepository
	logger *slog.Logger
}

func New(repo *inventory_repository.InventoryRepository, logger *slog.Logger) *InventoryService {
	return &InventoryService{
		repo:   repo,
		logger: logger,
	}
}

func (i *InventoryService) AddItem(ctx context.Context, itemID string, characterID int, quantity uint16) error {
	err := i.repo.AddItem(ctx, itemID, characterID, quantity)

	if err != nil {
		i.logger.Error("Error while add item to character", "error", err.Error(), "character_id", characterID, "item_id", itemID)
	}

	return err
}

func (i *InventoryService) RemoveItem(ctx context.Context, itemID string, characterID int) error {
	return i.repo.RemoveItem(ctx, itemID, characterID)
}

func (i *InventoryService) GetInventory(ctx context.Context, characterID int) ([]models.InventoryItem, error) {
	return i.repo.GetInventory(ctx, characterID)
}
