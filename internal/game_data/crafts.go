package game_data

type ItemType string

const (
	RegularItem ItemType = "item"
	WeaponType  ItemType = "weapon"
)

type Craft struct {
	ItemID  string                 `json:"item_id"`
	Recipes map[string]CraftRecipe `json:"recipes"`
}

type RecipeItem struct {
	ItemID   string `json:"item_id"`
	Quantity uint16 `json:"quantity"`
}

type CraftRecipe struct {
	Quantity   int          `json:"quantity"`
	EnergyCost int          `json:"energy_cost"`
	Items      []RecipeItem `json:"items"`
}
