package game_data

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type GameData struct {
	Items            map[string]Item
	Crafts           map[string]CraftRecipe
	WeaponAttributes map[string]WeaponAttribute
	Adventures       map[string]Adventure
	AdventureRewards map[string]AdventureRewards
	logger           *slog.Logger
}

type GameDataConstructor struct {
	CraftFile            string
	ItemsFile            string
	WeaponAttributesFile string
	AdventuresFile       string
	AdventureRewardsFile string
}

func (gd *GameData) Read(props GameDataConstructor) error {
	craftData, err := readJSONFile[map[string]CraftRecipe](props.CraftFile)
	if err != nil {
		return err
	}
	gd.logger.Info("Craft data loaded successfully", "entries", len(craftData))

	itemsData, err := readJSONFile[map[string]Item](props.ItemsFile)
	if err != nil {
		return err
	}
	gd.logger.Info("Items data loaded successfully", "entries", len(itemsData))

	weaponAttributesData, err := readJSONFile[map[string]WeaponAttribute](props.WeaponAttributesFile)
	if err != nil {
		return err
	}
	gd.logger.Info("Weapon attributes loaded successfully", "entries", len(weaponAttributesData))

	adventuresData, err := readJSONFile[map[string]Adventure](props.AdventuresFile)
	if err != nil {
		return err
	}
	gd.logger.Info("Adventures loaded successfully", "entries", len(adventuresData))

	adventuresRewardsData, err := readJSONFile[map[string]AdventureRewards](props.AdventureRewardsFile)
	if err != nil {
		return err
	}
	gd.logger.Info("Adventures rewards loaded successfully", "entries", len(adventuresRewardsData))

	gd.Crafts = craftData
	gd.Items = itemsData
	gd.WeaponAttributes = weaponAttributesData
	gd.Adventures = adventuresData
	gd.AdventureRewards = adventuresRewardsData

	gd.logger.Info("All game data loaded successfully!")

	return nil
}

func NewGameData(logger *slog.Logger) *GameData {
	return &GameData{
		Items:            make(map[string]Item),
		Crafts:           make(map[string]CraftRecipe),
		WeaponAttributes: make(map[string]WeaponAttribute),
		Adventures:       make(map[string]Adventure),
		logger:           logger,
	}
}

func readJSONFile[T any](pathToFile string) (T, error) {
	var data T
	file, err := os.Open(pathToFile)
	if err != nil {
		return data, fmt.Errorf("failed to open file %s: %w", pathToFile, err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return data, fmt.Errorf("failed to decode JSON from file %s: %w", pathToFile, err)
	}
	return data, nil
}
