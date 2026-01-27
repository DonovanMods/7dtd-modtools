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

func TestTemplateHelperFunctions(t *testing.T) {
	assert := setup(t)
	defer cleanup(t)

	// Create minimal game data
	gameDir := filepath.Join(testTMP, "gamedata")
	require.NoError(t, FS.MkdirAll(gameDir, 0755))
	require.NoError(t, FS.WriteFile(filepath.Join(gameDir, "blocks.xml"), []byte(`<blocks></blocks>`), 0644))
	require.NoError(t, FS.WriteFile(filepath.Join(gameDir, "entityclasses.xml"), []byte(`<entity_classes></entity_classes>`), 0644))

	// Template using helper functions
	tmplContent := `{{- modlet "test-helpers" -}}
{{- output "test.txt" -}}
hasPrefix: {{ hasPrefix "terrStone" "terr" }}
hasSuffix: {{ hasSuffix "terrStone" "Stone" }}
match: {{ match "plantedCorn" "^planted" }}
notMatch: {{ notMatch "terrStone" "^planted" }}
multValue: {{ multValue "55" 1.5 }}
probMult: {{ probMult "0.8" 1.5 }}
{{- write -}}`

	tmplPath := filepath.Join(testTMP, "helpers.tmpl")
	require.NoError(t, FS.WriteFile(tmplPath, []byte(tmplContent), 0644))

	args := modlet.CmdArgs{
		Input:   []string{tmplPath},
		Output:  testTMP,
		Gamedir: gameDir,
	}

	err := modlet.Unpack(tmplPath, args)
	require.NoError(t, err)

	content, err := FS.ReadFile(filepath.Join(testTMP, "test-helpers", "test.txt"))
	require.NoError(t, err)

	output := string(content)
	assert.Contains(output, "hasPrefix: true")
	assert.Contains(output, "hasSuffix: true")
	assert.Contains(output, "match: true")
	assert.Contains(output, "notMatch: true")
	assert.Contains(output, "multValue: 83")
	assert.Contains(output, "probMult: 1")
}
