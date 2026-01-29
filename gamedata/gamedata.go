package gamedata

import (
	"encoding/xml"
	"fmt"
	"path/filepath"

	"github.com/donovanmods/7dtd-gamedata/gamexml"
	"github.com/spf13/afero"
)

// FS is the filesystem to use (can be swapped for testing)
var FS = &afero.Afero{Fs: afero.NewOsFs()}

// GameData holds parsed game XML data
type GameData struct {
	Blocks        []gamexml.Block
	EntityClasses []gamexml.EntityClass
	Items         []gamexml.Item
	Recipes       []gamexml.Recipe
	Loot          []gamexml.LootGroup
}

// Load parses game XML files from the given directory
func Load(gamedir string) (*GameData, error) {
	gd := &GameData{}

	// Load blocks.xml
	blocksPath := filepath.Join(gamedir, "blocks.xml")
	if err := gd.loadBlocks(blocksPath); err != nil {
		return nil, fmt.Errorf("loading blocks.xml: %w", err)
	}

	// Load entityclasses.xml
	entityPath := filepath.Join(gamedir, "entityclasses.xml")
	if err := gd.loadEntityClasses(entityPath); err != nil {
		return nil, fmt.Errorf("loading entityclasses.xml: %w", err)
	}

	// Load items.xml
	itemsPath := filepath.Join(gamedir, "items.xml")
	if err := gd.loadItems(itemsPath); err != nil {
		return nil, fmt.Errorf("loading items.xml: %w", err)
	}

	// Load recipes.xml
	recipesPath := filepath.Join(gamedir, "recipes.xml")
	if err := gd.loadRecipes(recipesPath); err != nil {
		return nil, fmt.Errorf("loading recipes.xml: %w", err)
	}

	// Load loot.xml
	lootPath := filepath.Join(gamedir, "loot.xml")
	if err := gd.loadLoot(lootPath); err != nil {
		return nil, fmt.Errorf("loading loot.xml: %w", err)
	}

	return gd, nil
}

func (gd *GameData) loadBlocks(path string) error {
	data, err := FS.ReadFile(path)
	if err != nil {
		return err
	}

	var blocks gamexml.Blocks
	if err := xml.Unmarshal(data, &blocks); err != nil {
		return err
	}

	gd.Blocks = blocks.Block
	return nil
}

func (gd *GameData) loadEntityClasses(path string) error {
	data, err := FS.ReadFile(path)
	if err != nil {
		return err
	}

	var entities gamexml.EntityClasses
	if err := xml.Unmarshal(data, &entities); err != nil {
		return err
	}

	gd.EntityClasses = entities.EntityClass
	return nil
}

func (gd *GameData) loadItems(path string) error {
	data, err := FS.ReadFile(path)
	if err != nil {
		return err
	}

	var items gamexml.Items
	if err := xml.Unmarshal(data, &items); err != nil {
		return err
	}

	gd.Items = items.Item
	return nil
}

func (gd *GameData) loadRecipes(path string) error {
	data, err := FS.ReadFile(path)
	if err != nil {
		return err
	}

	var recipes gamexml.Recipes
	if err := xml.Unmarshal(data, &recipes); err != nil {
		return err
	}

	gd.Recipes = recipes.Recipe
	return nil
}

func (gd *GameData) loadLoot(path string) error {
	data, err := FS.ReadFile(path)
	if err != nil {
		return err
	}

	var loot gamexml.Loot
	if err := xml.Unmarshal(data, &loot); err != nil {
		return err
	}

	gd.Loot = loot.LootGroup
	return nil
}
