package data_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/yazmeyaa/ton_game/internal/game_data"
)

func TestGameData(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	gd := game_data.NewGameData(logger)

	testDataPath := filepath.Join("..", "testdata")
	craftFilePath := filepath.Join(testDataPath, "crafts.json")
	itemsFilePath := filepath.Join(testDataPath, "items.json")
	weaponAttributesFilePath := filepath.Join(testDataPath, "weapon_attributes.json")

	err := gd.Read(game_data.GameDataConstructor{
		ItemsFile:            itemsFilePath,
		CraftFile:            craftFilePath,
		WeaponAttributesFile: weaponAttributesFilePath,
	})

	if err != nil {
		t.Error(err.Error())
		return
	}

	stick, has := gd.Items["tree_stick"]

	if !has {
		t.Error("Missed item data \"tree_stick\"")
		return
	}

	if stick.Name != "Tree stick" {
		t.Error("Wrong item name value")
		return
	}
}
