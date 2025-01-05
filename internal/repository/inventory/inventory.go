package inventory_repository

type InventoryRepository struct{}

func New() *InventoryRepository {
	return &InventoryRepository{}
}
