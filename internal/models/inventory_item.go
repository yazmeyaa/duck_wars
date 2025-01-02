package models

type InventoryItem struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	IsSoldable    bool   `json:"soldable"`
	IsCraftable   bool   `json:"craftable"`
	IsUsedInCraft bool   `json:"userInCraft"`
}
