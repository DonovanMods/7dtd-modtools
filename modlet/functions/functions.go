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
package functions

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
	"github.com/donovanmods/7dtd-gamedata/modlet"
	"github.com/donovanmods/7dtd-modtools/lib/logger"
	"github.com/donovanmods/7dtd-modtools/modlet/common"
)

/*
// Helper functions
*/

type (
	null string // we use this when we need to output nothing
	Key  string
)

const (
	By  Key = "by"
	Min Key = "min"
	Max Key = "max"
)

var outputFound bool

func (k Key) String() string {
	return string(k)
}

func (k Key) IsValid() bool {
	validArgs := []Key{By, Min, Max}

	for _, v := range validArgs {
		if v == k {
			return true
		}
	}
	return false
}

/*
// Functions for building modlets
*/

// modlet sets up the modlet name and path
// func FuncModlet(outdir string, modInfo *modinfo.ModInfo) func(string) null {
func ModletFunc(fargs common.FuncArgs) func(string) null {
	return func(name string) null {
		name = strings.TrimSpace(name)

		if name == "" {
			logger.Fatal("modlet name must be provided")
		}

		path := filepath.Join(fargs.Outdir, name)

		*fargs.ModInfo = *modinfo.NewModInfo(name)
		fargs.ModInfo.SetPath(path)

		if err := mkPath(fargs.ModInfo.Path()); err != nil {
			logger.Panic(err)
		}

		logger.Info("creating modlet %q", fargs.ModInfo.GetValue("name"))

		return null("")
	}
}

func MultFunc(fargs common.FuncArgs) func(string, ...string) string {
	return func(xpath string, args ...string) string {
		var multiplier float64
		var err error

		xpath = strings.TrimSpace(xpath)
		if xpath == "" {
			logger.Fatal("xpath must be provided to the mult command")
		}

		pargs := ParseArgs(args)
		if len(pargs) == 0 {
			logger.Fatal("mult requires additional argument (by= at least)")
		}

		if by, ok := pargs["by"]; ok {
			if !ok || by == "" {
				logger.Fatal("mult requires a valid by= argument")
			}

			if multiplier, err = strconv.ParseFloat(by, 64); err != nil {
				logger.Fatal("error parsing multiplier %q: %w", by, err)
			}
		}

		return must(modlet.MkSet(xpath, strconv.FormatFloat(multiplier, 'f', -1, 64)))
	}
}

// output sets up the output file and creates a uniq buffer
// func FuncOutput(fBuffer *fileBuffer, gBuffer *bytes.Buffer, fBufMap fileBufferMap, modInfo *modinfo.ModInfo) func(string) null {
func OutputFunc(fargs common.FuncArgs) func(string) null {
	return func(path string) null {
		path = strings.TrimSpace(path)

		if path == "" {
			logger.Fatal("output file not provided")
		}

		if fargs.ModInfo.Path() == "" {
			logger.Fatal("please set the modlet using {{ modlet <name> }}")
		}

		cleanPath := filepath.Clean(path)
		fullPath := filepath.Join(fargs.ModInfo.Path(), cleanPath)

		logger.Trace("buffering output for %q", fullPath)

		if err := mkPath(filepath.Dir(fullPath)); err != nil {
			logger.Panic(err)
		}

		f, err := common.FS.Create(fullPath)
		if err != nil {
			logger.Fatal("error creating output file %s: %w", fullPath, err)
		}

		fargs.GBuffer.Reset()

		fargs.FBufMap[fullPath] = common.FileBuffer{
			Buffer: bytes.NewBuffer(nil),
			Writer: f,
		}
		*fargs.FBuffer = fargs.FBufMap[fullPath]

		outputFound = true

		return null("")
	}
}

// set produces a Set modlet instruction with the given xpath and value
func SetFunc(fargs common.FuncArgs) func(string, string) string {
	return func(xpath string, value string) string {
		return must(modlet.MkSet(xpath, value))
	}
}

// write writes the contents of the buffer to the output file
// func FuncWrite(fBuffer *fileBuffer, gBuffer *bytes.Buffer) func() null {
func WriteFunc(fargs common.FuncArgs) func() null {
	return func() null {
		if !outputFound || (*fargs.FBuffer).Writer == nil {
			logger.Fatal("you've called `write` without providing an output file, please use `output <filepath>` before `write`")
		}

		logger.Trace("saving fileBuffer")

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

func ParseArgs(args []string) map[Key]string {
	pargs := make(map[Key]string, len(args))

	re := regexp.MustCompile(`(?P<key>[^=]+)=(?P<value>.+)`)

	for _, arg := range args {
		if !re.MatchString(arg) {
			logger.Fatal("invalid argument format: %q", arg)
		}

		matches := re.FindStringSubmatch(arg)
		if len(matches) != 3 {
			logger.Fatal("invalid argument format: %q", arg)
		}

		key := Key(strings.ToLower(strings.TrimSpace(matches[1])))
		value := strings.TrimSpace(matches[2])

		if key == "" || value == "" {
			logger.Fatal("invalid input (key or value is empty)")
		}

		if !key.IsValid() {
			logger.Error("unknown argument key: %q", key)
			continue
		}

		pargs[Key(key)] = value
	}

	return pargs
}

// must is a helper function to handle errors
func must(output string, err error) string {
	if err != nil {
		logger.Fatal("error creating function: %v", err)
	}

	return output
}

// Helper function to create a directory if it doesn't exist
func mkPath(path string) error {
	if !fs.ValidPath(path) {
		return fmt.Errorf("invalid path %q", path)
	}

	_, err := common.FS.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		logger.Info("creating directory: %q", path)

		if err := common.FS.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("unable to create directory %q: %w", path, err)
		}
	}

	return nil
}
