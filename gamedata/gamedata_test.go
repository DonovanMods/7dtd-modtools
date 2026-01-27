package gamedata_test

import (
	"path/filepath"
	"testing"

	"github.com/donovanmods/7dtd-modtools/gamedata"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadGameData(t *testing.T) {
	// Use in-memory filesystem
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}
	gamedata.FS = fs

	tmpDir := "/tmp/gamedata"
	require.NoError(t, fs.MkdirAll(tmpDir, 0755))

	blocksXML := `<?xml version="1.0" encoding="UTF-8"?>
<blocks>
	<block name="terrStone">
		<property name="Material" value="Mite"/>
		<drop event="Harvest" name="resourceRockSmall" count="55"/>
	</block>
	<block name="terrDirt">
		<property name="Material" value="Mdirt"/>
	</block>
</blocks>`

	entityClassesXML := `<?xml version="1.0" encoding="UTF-8"?>
<entity_classes>
	<entity_class name="animalChicken">
		<property name="Class" value="EntityAnimal"/>
		<drop event="Harvest" name="foodRawMeat" count="5" tag="butcherHarvest"/>
	</entity_class>
</entity_classes>`

	require.NoError(t, fs.WriteFile(filepath.Join(tmpDir, "blocks.xml"), []byte(blocksXML), 0644))
	require.NoError(t, fs.WriteFile(filepath.Join(tmpDir, "entityclasses.xml"), []byte(entityClassesXML), 0644))

	gd, err := gamedata.Load(tmpDir)
	require.NoError(t, err)

	assert.Len(t, gd.Blocks, 2)
	assert.Equal(t, "terrStone", gd.Blocks[0].Name)
	assert.Equal(t, "Mite", gd.Blocks[0].GetProperty("Material"))

	assert.Len(t, gd.EntityClasses, 1)
	assert.Equal(t, "animalChicken", gd.EntityClasses[0].Name)
}
