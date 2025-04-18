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
package common

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
	"github.com/donovanmods/7dtd-modtools/lib/logger"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// FS is the filesystem interface used for file operations
var FS = &afero.Afero{Fs: afero.NewOsFs()}

type FileBuffer struct {
	Buffer *bytes.Buffer
	Writer io.WriteCloser
}

// Holds buffers and io.WriteCloser for each output file
type FileBufferMap map[string]FileBuffer

type FuncArgs struct {
	Output  string
	Gamedir string
	ModInfo *modinfo.ModInfo
	FBuffer *FileBuffer
	GBuffer *bytes.Buffer
	FBufMap FileBufferMap
	Options map[string]string
}

func OutputFile() (string, error) {
	output := viper.GetString("output")
	if output == "" {
		return "", fmt.Errorf("no output file not specified")
	}

	return output, nil
}

func OutputDir(overwrite bool) (string, error) {
	output := filepath.Clean(viper.GetString("output"))

	logger.Debug("Output directory: %s", output)

	if output == "." {
		return output, nil
	}

	exists, err := FS.Exists(output)
	if err != nil {
		return "", fmt.Errorf("error checking output directory %s: %w", output, err)
	}

	if exists {
		if !overwrite {
			return "", fmt.Errorf("output directory %s already exists", output)
		}

		if err := FS.RemoveAll(output); err != nil {
			return "", fmt.Errorf("error removing output directory %s: %w", output, err)
		}
	}

	if err := FS.MkdirAll(output, 0755); err != nil {
		return "", fmt.Errorf("error creating output directory %s: %w", output, err)
	}

	return output, nil
}
