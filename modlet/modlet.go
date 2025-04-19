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
	"fmt"
	"path/filepath"

	"github.com/donovanmods/7dtd-modtools/lib/logger"
	"github.com/spf13/afero"
)

// FS is the filesystem interface used for file operations
var FS = &afero.Afero{Fs: afero.NewOsFs()}

type CmdArgs struct {
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

func (CA CmdArgs) NewModlet() error {
	return NewModlet(CA)
}

func (CA CmdArgs) PackModlet() error {
	return Pack(CA)
}

func (CA CmdArgs) UnpackModlet() error {
	// gamedir string, output string) error {
	for _, t := range CA.Input {
		if err := Unpack(t, CA); err != nil {
			// t, opts.Gamedir, opts.Output, opts.Force); err != nil {
			return fmt.Errorf("error building modlet from template %s: %w", t, err)
		}
	}
	return nil
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
