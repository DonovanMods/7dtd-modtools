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
package builder

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
)

/*
// Helper functions
*/

type FuncArgs struct {
	Outdir        string
	Gamedir       string
	Templates     []string
	ModInfo       *modinfo.ModInfo
	FBuffer       *FileBuffer
	GBuffer       *bytes.Buffer
	FBufMap       FileBufferMap
	IoReader      io.Reader
	IoWriteCloser io.WriteCloser
}

// Helper function to create a directory if it doesn't exist
func mkPath(path string) error {
	if !fs.ValidPath(path) {
		return fmt.Errorf("invalid path %q", path)
	}

	_, err := FS.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		log.Printf("creating directory: %q", path)

		if err := FS.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("unable to create directory %q: %w", path, err)
		}
	}

	return nil
}

/*
// Functions for building modlets
*/

// modlet sets up the modlet name and path
// func FuncModlet(outdir string, modInfo *modinfo.ModInfo) func(string) null {
func FuncModlet(args FuncArgs) func(string) null {
	return func(name string) null {
		if name == "" {
			log.Fatal("modlet name must be provided")
		}

		path := filepath.Join(args.Outdir, name)

		*args.ModInfo = *modinfo.NewModInfo(name)
		args.ModInfo.SetPath(path)

		if err := mkPath(args.ModInfo.Path()); err != nil {
			log.Fatal(err)
		}

		log.Printf("creating modlet %q\n", args.ModInfo.GetValue("name"))

		return null("")
	}
}

// output sets up the output file and creates a uniq buffer
// func FuncOutput(fBuffer *fileBuffer, gBuffer *bytes.Buffer, fBufMap fileBufferMap, modInfo *modinfo.ModInfo) func(string) null {
func FuncOutput(args FuncArgs) func(string) null {
	return func(path string) null {
		var (
			err error
			f   io.WriteCloser
		)

		if path == "" {
			log.Fatal("output file not provided")
		}

		if args.ModInfo.Path() == "" {
			log.Fatal("please set the modlet using {{ modlet <name> }}")
		}

		cleanPath := filepath.Clean(path)
		fullPath := filepath.Join(args.ModInfo.Path(), cleanPath)

		log.Printf("buffering output for %q\n", fullPath)

		if err := mkPath(filepath.Dir(fullPath)); err != nil {
			log.Fatal(err)
		}

		f, err = FS.Create(fullPath)
		if err != nil {
			log.Fatalf("error creating output file %s: %v", fullPath, err)
		}

		args.GBuffer.Reset()

		args.FBufMap[fullPath] = FileBuffer{
			Buffer: bytes.NewBuffer(nil),
			Writer: f,
		}
		*args.FBuffer = args.FBufMap[fullPath]

		outputFound = true

		return null("")
	}
}

// write writes the contents of the buffer to the output file
// func FuncWrite(fBuffer *fileBuffer, gBuffer *bytes.Buffer) func() null {
func FuncWrite(args FuncArgs) func() null {
	return func() null {
		if !outputFound || (*args.FBuffer).Writer == nil {
			log.Fatal("you've called `write` without providing an output file, please use `output <filepath>` before `write`")
		}

		log.Println("saving fileBuffer")

		// Copy the current buffer to the output buffer
		(*args.FBuffer).Buffer.Write(args.GBuffer.Bytes())

		// Reset the output state
		outputFound = false

		return null("")
	}
}
