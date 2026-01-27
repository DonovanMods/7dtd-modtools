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
