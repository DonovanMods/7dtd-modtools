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
	"testing"

	"github.com/donovanmods/7dtd-modtools/gamedata"
	"github.com/donovanmods/7dtd-modtools/lib/logger"
	"github.com/donovanmods/7dtd-modtools/modlet"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var (
	testTMP = "test_temp"
	FS      = &afero.Afero{Fs: afero.NewMemMapFs()}
)

func mkTempDir(t *testing.T) string {
	t.Helper()

	if exists, err := afero.DirExists(FS, testTMP); err != nil {
		t.Fatalf("Error checking for temp directory: %v", err)
	} else if !exists {
		err := FS.Mkdir(testTMP, 0700)
		if err != nil {
			t.Fatal(err)
		}
		logger.Info("Created temp directory: %s", testTMP)
	}

	return testTMP
}

func setup(t *testing.T) *assert.Assertions {
	t.Helper()

	logger.Testing = true

	// Use MemMapFs for testing
	modlet.FS = FS
	gamedata.FS = FS

	mkTempDir(t)

	return assert.New(t)
}

func cleanup(t *testing.T) {
	t.Helper()

	// Remove the temp directory
	if err := FS.RemoveAll(testTMP); err != nil {
		t.Fatalf("Error removing temp directory: %v", err)
	}
	logger.Info("Removed temp directory: %s", testTMP)
}

func TestCmdArgsSanitize(t *testing.T) {
	// Setup
	assert := setup(t)
	defer cleanup(t)

	// Test cases
	tests := []struct {
		name     string
		input    modlet.CmdArgs
		expected modlet.CmdArgs
	}{
		{
			name: "valid input",
			input: modlet.CmdArgs{
				Output: "//test.txt/",
				Force:  true,
			},
			expected: modlet.CmdArgs{
				Output:   "/test.txt",
				Gamedir:  ".",
				Compress: false,
				Force:    true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.input.Sanitize()
			assert.Equal(test.expected, test.input)
		})
	}
}
