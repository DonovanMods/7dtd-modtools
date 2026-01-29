/*
Copyright © 2025 Donovan C. Young <dyoung522@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
*/
package modlet_test

import (
	"path/filepath"
	"testing"

	"github.com/donovanmods/7dtd-modtools/modlet"
	"github.com/stretchr/testify/require"
)

func TestUnpackWithGameData(t *testing.T) {
	assert := setup(t)
	defer cleanup(t)

	// Create game data directory with test XML
	gameDir := filepath.Join(testTMP, "gamedata")
	require.NoError(t, FS.MkdirAll(gameDir, 0755))

	blocksXML := `<?xml version="1.0" encoding="UTF-8"?>
<blocks>
	<block name="terrStone">
		<property name="Material" value="Mite"/>
		<drop event="Harvest" name="resourceRockSmall" count="55"/>
	</block>
</blocks>`

	entityXML := `<?xml version="1.0" encoding="UTF-8"?>
<entity_classes>
	<entity_class name="animalChicken">
		<drop event="Harvest" name="foodRawMeat" count="5" tag="butcherHarvest"/>
	</entity_class>
</entity_classes>`

	// Minimal stubs so gamedata.Load() succeeds (it requires all five files)
	itemsXML := `<?xml version="1.0" encoding="UTF-8"?><items></items>`
	recipesXML := `<?xml version="1.0" encoding="UTF-8"?><recipes></recipes>`
	lootXML := `<?xml version="1.0" encoding="UTF-8"?><lootcontainers></lootcontainers>`

	require.NoError(t, FS.WriteFile(filepath.Join(gameDir, "blocks.xml"), []byte(blocksXML), 0644))
	require.NoError(t, FS.WriteFile(filepath.Join(gameDir, "entityclasses.xml"), []byte(entityXML), 0644))
	require.NoError(t, FS.WriteFile(filepath.Join(gameDir, "items.xml"), []byte(itemsXML), 0644))
	require.NoError(t, FS.WriteFile(filepath.Join(gameDir, "recipes.xml"), []byte(recipesXML), 0644))
	require.NoError(t, FS.WriteFile(filepath.Join(gameDir, "loot.xml"), []byte(lootXML), 0644))

	// Create template that uses game data
	tmplContent := `{{- modlet "test-gamedata" -}}
{{- output "Config/blocks.xml" -}}
{{ xmlHeader }}
<configs>
{{- range .GameData.Blocks }}
{{ comment .Name }}
{{- end }}
</configs>
{{- write -}}`

	tmplPath := filepath.Join(testTMP, "test.tmpl")
	require.NoError(t, FS.WriteFile(tmplPath, []byte(tmplContent), 0644))

	// Unpack with game data
	args := modlet.CmdArgs{
		Input:   []string{tmplPath},
		Output:  testTMP,
		Gamedir: gameDir,
	}

	err := modlet.Unpack(tmplPath, args)
	require.NoError(t, err)

	// Verify output contains block name in comment
	outputPath := filepath.Join(testTMP, "test-gamedata", "Config", "blocks.xml")
	content, err := FS.ReadFile(outputPath)
	require.NoError(t, err)

	assert.Contains(string(content), "<!-- terrStone -->")
}
