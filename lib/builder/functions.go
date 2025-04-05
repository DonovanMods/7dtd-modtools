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
	"encoding/xml"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
	"github.com/donovanmods/7dtd-gamedata/modlet"
	"github.com/donovanmods/7dtd-gamedata/xmltools"
)

/*
// Helper functions
*/

type FuncArgs struct {
	Outdir  string
	Gamedir string
	ModInfo *modinfo.ModInfo
	FBuffer *FileBuffer
	GBuffer *bytes.Buffer
	FBufMap FileBufferMap
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
func FuncModlet(fargs FuncArgs) func(string) null {
	return func(name string) null {
		name = strings.TrimSpace(name)

		if name == "" {
			log.Fatal("modlet name must be provided")
		}

		path := filepath.Join(fargs.Outdir, name)

		*fargs.ModInfo = *modinfo.NewModInfo(name)
		fargs.ModInfo.SetPath(path)

		if err := mkPath(fargs.ModInfo.Path()); err != nil {
			log.Fatal(err)
		}

		log.Printf("creating modlet %q\n", fargs.ModInfo.GetValue("name"))

		return null("")
	}
}

func FuncMult(fargs FuncArgs) func(string, ...string) string {
	return func(string, ...string) string {
		return ""
	}
}

// output sets up the output file and creates a uniq buffer
// func FuncOutput(fBuffer *fileBuffer, gBuffer *bytes.Buffer, fBufMap fileBufferMap, modInfo *modinfo.ModInfo) func(string) null {
func FuncOutput(fargs FuncArgs) func(string) null {
	return func(path string) null {
		path = strings.TrimSpace(path)

		if path == "" {
			log.Fatal("output file not provided")
		}

		if fargs.ModInfo.Path() == "" {
			log.Fatal("please set the modlet using {{ modlet <name> }}")
		}

		cleanPath := filepath.Clean(path)
		fullPath := filepath.Join(fargs.ModInfo.Path(), cleanPath)

		log.Printf("buffering output for %q\n", fullPath)

		if err := mkPath(filepath.Dir(fullPath)); err != nil {
			log.Fatal(err)
		}

		f, err := FS.Create(fullPath)
		if err != nil {
			log.Fatalf("error creating output file %s: %v", fullPath, err)
		}

		fargs.GBuffer.Reset()

		fargs.FBufMap[fullPath] = FileBuffer{
			Buffer: bytes.NewBuffer(nil),
			Writer: f,
		}
		*fargs.FBuffer = fargs.FBufMap[fullPath]

		outputFound = true

		return null("")
	}
}

// set produces a Set modlet instruction with the given xpath and value
func FuncSet(fargs FuncArgs) func(string, string) string {
	return mkSet
}

// write writes the contents of the buffer to the output file
// func FuncWrite(fBuffer *fileBuffer, gBuffer *bytes.Buffer) func() null {
func FuncWrite(fargs FuncArgs) func() null {
	return func() null {
		if !outputFound || (*fargs.FBuffer).Writer == nil {
			log.Fatal("you've called `write` without providing an output file, please use `output <filepath>` before `write`")
		}

		log.Println("saving fileBuffer")

		// Copy the current buffer to the output buffer
		(*fargs.FBuffer).Buffer.Write(fargs.GBuffer.Bytes())

		// Reset the output state
		outputFound = false

		return null("")
	}
}

/*
// Helper functions
*/

// mkSet creates a Set modlet instruction with the given xpath and value
func mkSet(xpath, value string) string {
	xpath = strings.TrimSpace(xpath)
	value = strings.TrimSpace(value)

	if xpath == "" || value == "" {
		log.Fatal("xpath and value must be provided")
	}

	log.Printf("creating set modlet for xpath %q with value %q\n", xpath, value)

	modlet := modlet.Modlet{
		XMLName: xml.Name{Local: "set"},
		XPath:   xpath,
		Value:   value,
	}

	set, err := xml.Marshal(modlet)
	if err != nil {
		log.Fatalf("error marshalling modlet: %v", err)
	}

	return string(xmltools.UnescapeXML(set))
}
