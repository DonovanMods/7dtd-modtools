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
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
	"github.com/donovanmods/7dtd-gamedata/modlet"
	"github.com/donovanmods/7dtd-modtools/lib/logger"
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
func ModletFunc(fargs FuncArgs) func(string) null {
	return func(name string) null {
		name = strings.TrimSpace(name)

		if name == "" {
			logger.Fatal("modlet name must be provided")
		}

		path := filepath.Join(fargs.Output, name)

		if e, err := FS.Exists(path); err != nil {
			logger.Fatal("error checking for modlet %q: %w", name, err)
		} else if e {
			if fargs.Options["force"] == strconv.FormatBool(true) {
				logger.Warn("modlet %q already exists, overwriting", name)
				if err := FS.RemoveAll(path); err != nil {
					logger.Fatal("error removing existing modlet %q: %w", name, err)
				}
			} else {
				logger.Fatal("modlet %q already exists, use --force to overwrite", name)
			}
		}

		*fargs.ModInfo = *modinfo.NewModInfo(name)
		fargs.ModInfo.SetPath(path)

		if err := MkPath(fargs.ModInfo.Path()); err != nil {
			logger.Panic(err)
		}

		logger.Info("creating modlet %q", fargs.ModInfo.GetValue("name"))

		return null("")
	}
}

// MultFunc creates a multiplier instruction for modlet templates.
// Supports by=N (required), min=N (optional), max=N (optional)
func MultFunc(fargs FuncArgs) func(string, ...string) string {
	return func(xpath string, args ...string) string {
		var multiplier, minVal, maxVal float64
		var hasMin, hasMax bool
		var err error

		xpath = strings.TrimSpace(xpath)
		if xpath == "" {
			logger.Fatal("xpath must be provided to the mult command")
		}

		pargs := ParseArgs(args)
		if len(pargs) == 0 {
			logger.Fatal("mult requires additional argument (by= at least)")
		}

		// Parse by (required)
		if by, ok := pargs["by"]; ok {
			if by == "" {
				logger.Fatal("mult requires a valid by= argument")
			}
			if multiplier, err = strconv.ParseFloat(by, 64); err != nil {
				logger.Fatal("error parsing multiplier %q: %w", by, err)
			}
		} else {
			logger.Fatal("mult requires by= argument")
		}

		// Parse min (optional)
		if min, ok := pargs["min"]; ok && min != "" {
			if minVal, err = strconv.ParseFloat(min, 64); err != nil {
				logger.Fatal("error parsing min %q: %w", min, err)
			}
			hasMin = true
		}

		// Parse max (optional)
		if max, ok := pargs["max"]; ok && max != "" {
			if maxVal, err = strconv.ParseFloat(max, 64); err != nil {
				logger.Fatal("error parsing max %q: %w", max, err)
			}
			hasMax = true
		}

		// Apply bounds
		result := multiplier
		if hasMin && result < minVal {
			result = minVal
		}
		if hasMax && result > maxVal {
			result = maxVal
		}

		return must(mkSet(xpath, strconv.FormatFloat(result, 'f', -1, 64)))
	}
}

// output sets up the output file and creates a uniq buffer
// func FuncOutput(fBuffer *fileBuffer, gBuffer *bytes.Buffer, fBufMap fileBufferMap, modInfo *modinfo.ModInfo) func(string) null {
func OutputFunc(fargs FuncArgs) func(string) null {
	return func(path string) null {
		path = strings.TrimSpace(path)

		if path == "" {
			logger.Fatal("output file not provided")
		}

		if fargs.ModInfo.Path() == "" {
			logger.Fatal("you must set the modlet using {{ modlet <name> }} prior to any other operations")
		}

		cleanPath := filepath.Clean(path)
		fullPath := filepath.Join(fargs.ModInfo.Path(), cleanPath)

		logger.Trace("buffering output for %q", fullPath)

		if err := MkPath(filepath.Dir(fullPath)); err != nil {
			logger.Panic(err)
		}

		f, err := FS.Create(fullPath)
		if err != nil {
			logger.Fatal("error creating output file %s: %w", fullPath, err)
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
func SetFunc(fargs FuncArgs) func(string, string) string {
	return func(xpath string, value string) string {
		return must(mkSet(xpath, value))
	}
}

// CommentFunc returns an XML comment string
func CommentFunc() func(string) string {
	return func(text string) string {
		return fmt.Sprintf("<!-- %s -->", text)
	}
}

// write writes the contents of the buffer to the output file
// func FuncWrite(fBuffer *fileBuffer, gBuffer *bytes.Buffer) func() null {
func WriteFunc(fargs FuncArgs) func() null {
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

// mkSet creates a Set modlet instruction XML string
func mkSet(xpath string, value string) (string, error) {
	m := modlet.Modlet{
		XMLName: xml.Name{Local: "set"},
		XPath:   xpath,
		Value:   value,
	}
	b, err := xml.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("error marshaling set instruction: %w", err)
	}
	// Unescape single quotes that xml.Marshal escapes (7DTD expects unescaped quotes in XPath)
	result := strings.ReplaceAll(string(b), "&#39;", "'")
	return result, nil
}

// Helper function to create a directory if it doesn't exist
func MkPath(path string) error {
	path = filepath.Clean(path)

	if !fs.ValidPath(path) {
		return fmt.Errorf("invalid path %q", path)
	}

	_, err := FS.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		logger.Info("creating directory: %q", path)

		if err := FS.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("unable to create directory %q: %w", path, err)
		}
	}

	return nil
}
