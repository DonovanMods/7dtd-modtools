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
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
	"github.com/donovanmods/7dtd-modtools/lib/logger"
	"github.com/spf13/afero"
)

// FS is the filesystem interface used for file operations
var FS = &afero.Afero{Fs: afero.NewOsFs()}

type FileBuffer struct {
	Buffer *bytes.Buffer
	Writer io.WriteCloser
}

type null string // we use this when we need to output nothing

// Holds buffers and io.WriteCloser for each output file
type FileBufferMap map[string]FileBuffer

// Used to track output/write state
var outputFound bool

func BuildModlets(templates []string, gamedir string, outdir string) error {
	for _, t := range templates {
		if err := BuildModlet(t, gamedir, outdir); err != nil {
			return fmt.Errorf("error building modlet from template %s: %w", t, err)
		}
	}
	return nil
}

func BuildModlet(tmpl string, gamedir string, outdir string) error {
	var modInfo modinfo.ModInfo

	fBufMap := make(FileBufferMap)
	gBuffer := bytes.NewBuffer(nil)
	fBuffer := &FileBuffer{
		Buffer: gBuffer,
		Writer: nil,
	}

	if strings.TrimSpace(tmpl) == "" {
		return errors.New("no templates provided")
	}

	if strings.TrimSpace(gamedir) == "" {
		return errors.New("gamedir not provided")
	}

	if strings.TrimSpace(outdir) == "" {
		return errors.New("outdir not provided")
	}

	outdir = filepath.Clean(outdir)
	templateName := filepath.Base(tmpl)

	logger.Debug("processing template: %s", templateName)

	fargs := FuncArgs{outdir, gamedir, &modInfo, fBuffer, gBuffer, fBufMap}
	t, err := template.New(templateName).
		Funcs(template.FuncMap{
			"modlet":    ModletFunc(fargs),
			"mult":      MultFunc(fargs),
			"output":    OutputFunc(fargs),
			"set":       SetFunc(fargs),
			"write":     WriteFunc(fargs),
			"xmlHeader": func() string { return xml.Header },
		}).
		ParseFiles(tmpl)
	if err != nil {
		return fmt.Errorf("error parsing template %s: %w", templateName, err)
	}

	if err := t.ExecuteTemplate(gBuffer, templateName, nil); err != nil {
		return fmt.Errorf("error executing template %s: %w", templateName, err)
	}

	// Write our fBuffer to disk
	for path, fBuffer := range fBufMap {
		if err := writeBuf(path, fBuffer); err != nil {
			logger.Panic(err)
		}
	}

	return nil
}

func writeBuf(path string, fBuffer FileBuffer) error {
	path = strings.TrimSpace(path)

	logger.Info("writing %q", path)

	if fBuffer.Writer != nil {
		defer func() {
			if err := fBuffer.Writer.Close(); err != nil {
				logger.Fatal("error closing output file: %w", err)
			}
		}()

		if _, err := fBuffer.Buffer.WriteTo(fBuffer.Writer); err != nil {
			return fmt.Errorf("error writing to output file %s: %w", path, err)
		}
	}

	return nil
}
