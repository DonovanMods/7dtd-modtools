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
package modlet

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/donovanmods/7dtd-modtools/lib/logger"
	"github.com/spf13/afero"
)

// FS is the filesystem interface used for file operations
var FS = &afero.Afero{Fs: afero.NewOsFs()}

type CmdArgs struct {
	Name     string
	Input    []string
	Output   string
	Gamedir  string
	Compress bool
	Force    bool
}

func (CA *CmdArgs) Sanitize() CmdArgs {
	for i, input := range CA.Input {
		CA.Input[i] = filepath.Clean(input)
	}
	CA.Output = filepath.Clean(CA.Output)
	CA.Gamedir = filepath.Clean(CA.Gamedir)

	return *CA
}

func NewCmd(args CmdArgs) error {
	return New(args)
}

func PackCmd(args CmdArgs) error {
	return Pack(args)
}

func UnpackCmd(args CmdArgs) error {
	// gamedir string, output string) error {
	for _, t := range args.Input {
		if err := Unpack(t, args); err != nil {
			// t, opts.Gamedir, opts.Output, opts.Force); err != nil {
			return fmt.Errorf("error building modlet from template %s: %w", t, err)
		}
	}
	return nil
}

func ValidateOutputFile(file string) (string, error) {
	file = filepath.Clean(file)

	if file == "" {
		return "", fmt.Errorf("no output file not specified")
	}

	if dir, err := FS.IsDir(file); err != nil {
		if !errors.Is(err, afero.ErrFileNotFound) {
			return "", fmt.Errorf("error checking output %s: %w", file, err)
		}
	} else if dir {
		return "", fmt.Errorf("output %s is a directory, want a file", file)
	}

	return file, nil
}

func ValidateOutputDir(dir string) (string, error) {
	dir = filepath.Clean(dir)

	if d, err := FS.IsDir(dir); err != nil {
		return "", fmt.Errorf("error checking output %s: %w", dir, err)
	} else if !d {
		return "", fmt.Errorf("output %s is not a directory", dir)
	}

	logger.Debug("Output directory: %s", dir)

	if dir == "." {
		return dir, nil
	}

	exists, err := FS.Exists(dir)
	if err != nil {
		return "", fmt.Errorf("error checking output directory %s: %w", dir, err)
	}

	if !exists {
		if err := FS.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("error creating output directory %s: %w", dir, err)
		}
	}

	return dir, nil
}

func CheckErr(err error) {
	if err != nil {
		logger.Fatal("Error: %v", err)
	}
}

func CheckValue(value any, err error) any {
	CheckErr(err)

	return value
}
