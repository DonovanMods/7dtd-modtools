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
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"testing"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
	"github.com/donovanmods/7dtd-modtools/modlet"
	"github.com/stretchr/testify/require"
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

func TestModletFunc(t *testing.T) {
	assert := setup(t)

	funcArgs := modlet.FuncArgs{
		Output:  testTMP,
		ModInfo: &modinfo.ModInfo{},
	}

	fn := modlet.ModletFunc(funcArgs)
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

	funcArgs := modlet.FuncArgs{
		Output:  testTMP,
		ModInfo: &modinfo.ModInfo{},
		FBuffer: &modlet.FileBuffer{},
		GBuffer: bytes.NewBuffer(nil),
		FBufMap: make(modlet.FileBufferMap),
	}

	funcArgs.ModInfo.SetPath(funcArgs.Output)

	fn := modlet.OutputFunc(funcArgs)
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

	funcArgs := modlet.FuncArgs{
		ModInfo: &modinfo.ModInfo{},
	}

	fn := modlet.SetFunc(funcArgs)

	xpath := `//block[@name='terrStone']/drop[@event='Harvest' and @name='resourceRockSmall']/@count`
	value := "999"

	expected := fmt.Sprintf("<set xpath=\"%s\">%s</set>", xpath, value)
	actual := fn(xpath, value)

	assert.Equal(expected, actual, "FuncSet should return the correct XML string")
}

func TestWriteFunc(t *testing.T) {
	assert := setup(t)

	buf := bytes.NewBuffer(nil)

	funcArgs := modlet.FuncArgs{
		FBuffer: &modlet.FileBuffer{
			Buffer: buf,
			Writer: NewNopIO(t, buf), // Using stdout for testing
		},
		GBuffer: bytes.NewBufferString("Test content"),
	}

	fn := modlet.WriteFunc(funcArgs)
	fn()

	assert.Equal("Test content", funcArgs.FBuffer.Buffer.String(), "Buffer content should match")
}

func TestParseArgs(t *testing.T) {
	assert := setup(t)

	inputs := []struct {
		Args     []string
		expected map[modlet.Key]string
	}{
		{
			Args: []string{"by=2.25", "min=4", "max=25"},
			expected: map[modlet.Key]string{
				"by":  "2.25",
				"max": "25",
				"min": "4",
			},
		},
		{
			Args: []string{"By=2.25", "Min=4"},
			expected: map[modlet.Key]string{
				"by":  "2.25",
				"min": "4",
			},
		},
		{
			Args: []string{"BY=2.25", "MIN=4"},
			expected: map[modlet.Key]string{
				"by":  "2.25",
				"min": "4",
			},
		},
		{
			Args:     []string{"by=1", "invalid=foo"},
			expected: map[modlet.Key]string{"by": "1"},
		},
		{
			Args:     []string{},
			expected: map[modlet.Key]string{},
		},
	}

	for _, input := range inputs {
		actual := modlet.ParseArgs(input.Args)
		assert.Equal(input.expected, actual, "ParseArgs should return the correct string")
	}
}

func TestParseArgsInvalid(t *testing.T) {
	_ = setup(t)

	inputs := []struct {
		Args     []string
		expected map[modlet.Key]string
	}{
		{
			Args: []string{"by="},
		},
	}

	for _, input := range inputs {
		require.Panics(t, func() {
			modlet.ParseArgs(input.Args)
		}, "invalid input (key or value is empty)")
	}
}

func TestMkPath(t *testing.T) {
	assert := setup(t)

	path := "test/path/to/dir"

	err := modlet.MkPath(path)
	assert.NoError(err, "Error creating directory")

	exists, err := FS.DirExists(path)
	assert.NoError(err, "Error checking for directory")
	assert.True(exists, "Directory should exist")
}

func TestMultFuncWithBounds(t *testing.T) {
	assert := setup(t)
	defer cleanup(t)

	funcArgs := modlet.FuncArgs{
		ModInfo: &modinfo.ModInfo{},
	}

	fn := modlet.MultFunc(funcArgs)

	// Basic multiplication (no bounds)
	result := fn("//test/@value", "by=2.0")
	assert.Contains(result, "<set")
	assert.Contains(result, "2")

	// With min bound - result should be clamped to min
	result = fn("//test/@value", "by=0.5", "min=1")
	assert.Contains(result, ">1<")

	// With max bound - result should be clamped to max
	result = fn("//test/@value", "by=10.0", "max=5")
	assert.Contains(result, ">5<")

	// Both bounds
	result = fn("//test/@value", "by=0.1", "min=2", "max=10")
	assert.Contains(result, ">2<")
}

func TestCommentFunc(t *testing.T) {
	_ = setup(t)

	fn := modlet.CommentFunc()

	result := fn("This is a comment")
	require.Equal(t, "<!-- This is a comment -->", result)

	result = fn("Multiple\nlines")
	require.Equal(t, "<!-- Multiple\nlines -->", result)
}

func TestProbFunc(t *testing.T) {
	assert := setup(t)

	funcArgs := modlet.FuncArgs{
		ModInfo: &modinfo.ModInfo{},
	}

	fn := modlet.ProbFunc(funcArgs)

	// Normal multiplication
	result := fn("//block/@prob", "0.3", "by=1.5")
	assert.Contains(result, "<set")
	assert.Contains(result, "0.45")

	// Capped at 1.0
	result = fn("//block/@prob", "0.8", "by=1.5")
	assert.Contains(result, ">1<")
}
