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
package common_test

import (
	"testing"

	"github.com/donovanmods/7dtd-modtools/lib/logger"
	"github.com/donovanmods/7dtd-modtools/modlet/lib/common"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
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
	common.FS = FS

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

func TestOutputFile(t *testing.T) {
	// Setup
	assert := setup(t)
	defer cleanup(t)

	// Test cases
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"valid", "test.txt", "test.txt", false},
		{"empty", "", "", true},
	}

	for _, test := range tests {
		viper.Set("output", test.input)

		t.Run(test.name, func(t *testing.T) {
			result, err := common.OutputFile()

			if test.wantErr {
				assert.Error(err)
				return
			}

			assert.NoError(err)
			assert.Equal(test.expected, result)
		})
	}
}

func TestOutputDir(t *testing.T) {
	// Setup
	assert := setup(t)
	defer cleanup(t)

	// Test cases
	tests := []struct {
		name      string
		input     string
		expected  string
		overwrite bool
		wantErr   bool
	}{
		{"valid overwrite", testTMP, testTMP, true, false},
		{"valid error", testTMP, testTMP, false, true},
		{"empty", "", ".", false, false},
	}

	for _, test := range tests {
		viper.Set("output", test.input)

		t.Run(test.name, func(t *testing.T) {
			actual, err := common.OutputDir(test.overwrite)

			if test.wantErr {
				assert.Error(err)
				return
			}

			assert.NoError(err)
			assert.Equal(test.expected, actual)
		})
	}
}
