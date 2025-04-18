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
package functions_test

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"testing"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
	"github.com/donovanmods/7dtd-modtools/lib/logger"
	"github.com/donovanmods/7dtd-modtools/modlet/lib/common"
	"github.com/donovanmods/7dtd-modtools/modlet/lib/functions"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testTMP = "test_temp"
	FS      = &afero.Afero{Fs: afero.NewMemMapFs()}
)

type nopIO struct {
	io.Reader
	io.Writer
}

func (nopIO) Close() error { return nil }

func NewNopIO(t *testing.T, buf *bytes.Buffer) io.WriteCloser {
	t.Helper()

	return nopIO{buf, buf}
}

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

func TestModletFunc(t *testing.T) {
	assert := setup(t)

	funcArgs := common.FuncArgs{
		Output:  testTMP,
		ModInfo: &modinfo.ModInfo{},
	}

	fn := functions.ModletFunc(funcArgs)
	modletName := "TestModlet"

	fn(modletName)

	assert.Equal(modletName, funcArgs.ModInfo.GetValue("name"))
	assert.Equal(filepath.Join(funcArgs.Output, modletName), funcArgs.ModInfo.Path())

	exists, err := FS.DirExists(funcArgs.ModInfo.Path())
	assert.NoError(err, "Error checking for Modlet directory")
	assert.True(exists, "Modlet directory should exist")
}

func TestOutputFunc(t *testing.T) {
	assert := setup(t)

	funcArgs := common.FuncArgs{
		Output:  testTMP,
		ModInfo: &modinfo.ModInfo{},
		FBuffer: &common.FileBuffer{},
		GBuffer: bytes.NewBuffer(nil),
		FBufMap: make(common.FileBufferMap),
	}

	funcArgs.ModInfo.SetPath(funcArgs.Output)

	fn := functions.OutputFunc(funcArgs)
	outputPath := "testfile.txt"

	fn(outputPath)

	fullPath := filepath.Join(funcArgs.Output, outputPath)
	assert.Contains(funcArgs.FBufMap, fullPath, "File buffer map should contain the output path")

	exists, err := FS.Exists(funcArgs.ModInfo.Path())
	assert.NoError(err, "Error checking for output file")
	assert.True(exists, "Output file should exist")
}

func TestSetFunc(t *testing.T) {
	assert := setup(t)

	funcArgs := common.FuncArgs{
		ModInfo: &modinfo.ModInfo{},
	}

	fn := functions.SetFunc(funcArgs)

	xpath := `//block[@name='terrStone']/drop[@event='Harvest' and @name='resourceRockSmall']/@count`
	value := "999"

	expected := fmt.Sprintf("<set xpath=\"%s\">%s</set>", xpath, value)
	actual := fn(xpath, value)

	assert.Equal(expected, actual, "FuncSet should return the correct XML string")
}

func TestWriteFunc(t *testing.T) {
	assert := setup(t)

	buf := bytes.NewBuffer(nil)

	funcArgs := common.FuncArgs{
		FBuffer: &common.FileBuffer{
			Buffer: buf,
			Writer: NewNopIO(t, buf), // Using stdout for testing
		},
		GBuffer: bytes.NewBufferString("Test content"),
	}

	fn := functions.WriteFunc(funcArgs)
	fn()

	assert.Equal("Test content", funcArgs.FBuffer.Buffer.String(), "Buffer content should match")
}

func TestParseArgs(t *testing.T) {
	assert := setup(t)

	inputs := []struct {
		Args     []string
		expected map[functions.Key]string
	}{
		{
			Args: []string{"by=2.25", "min=4", "max=25"},
			expected: map[functions.Key]string{
				"by":  "2.25",
				"max": "25",
				"min": "4",
			},
		},
		{
			Args: []string{"By=2.25", "Min=4"},
			expected: map[functions.Key]string{
				"by":  "2.25",
				"min": "4",
			},
		},
		{
			Args: []string{"BY=2.25", "MIN=4"},
			expected: map[functions.Key]string{
				"by":  "2.25",
				"min": "4",
			},
		},
		{
			Args:     []string{"by=1", "invalid=foo"},
			expected: map[functions.Key]string{"by": "1"},
		},
		{
			Args:     []string{},
			expected: map[functions.Key]string{},
		},
	}

	for _, input := range inputs {
		actual := functions.ParseArgs(input.Args)
		assert.Equal(input.expected, actual, "ParseArgs should return the correct string")
	}
}

func TestParseArgsInvalid(t *testing.T) {
	_ = setup(t)

	inputs := []struct {
		Args     []string
		expected map[functions.Key]string
	}{
		{
			Args: []string{"by="},
		},
	}

	for _, input := range inputs {
		require.Panics(t, func() {
			functions.ParseArgs(input.Args)
		}, "invalid input (key or value is empty)")
	}
}

func TestMkPath(t *testing.T) {
	assert := setup(t)

	path := "test/path/to/dir"

	err := functions.MkPath(path)
	assert.NoError(err, "Error creating directory")

	exists, err := FS.DirExists(path)
	assert.NoError(err, "Error checking for directory")
	assert.True(exists, "Directory should exist")
}
