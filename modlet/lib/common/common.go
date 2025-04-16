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
	"io"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
	"github.com/spf13/afero"
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
	Outdir  string
	Gamedir string
	ModInfo *modinfo.ModInfo
	FBuffer *FileBuffer
	GBuffer *bytes.Buffer
	FBufMap FileBufferMap
}
