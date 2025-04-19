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
	"compress/gzip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
	"github.com/donovanmods/7dtd-modtools/lib/logger"
)

func Unpack(tmpl string, args CmdArgs) error {
	args.Sanitize()

	var (
		gamedir = args.Gamedir
		output  = args.Output
		force   = args.Force
		modInfo modinfo.ModInfo
	)

	tmpl = filepath.Clean(tmpl)
	if tmpl == "" {
		return errors.New("no templates provided")
	}

	templateName := filepath.Base(tmpl)
	if templateName == "" {
		return errors.New("no template name provided")
	}

	if err := ValidateTemplate(tmpl); err != nil {
		return fmt.Errorf("error validating template %s: %w", tmpl, err)
	}

	gamedir = filepath.Clean(gamedir)
	output = filepath.Clean(output)

	fBufMap := make(FileBufferMap)
	gBuffer := bytes.NewBuffer(nil)
	fBuffer := &FileBuffer{
		Buffer: gBuffer,
		Writer: nil,
	}

	logger.Debug("processing template: %s", templateName)

	fargs := FuncArgs{
		Output:  output,
		Gamedir: gamedir,
		ModInfo: &modInfo,
		FBuffer: fBuffer,
		GBuffer: gBuffer,
		FBufMap: fBufMap,
		Options: map[string]string{
			"force": strconv.FormatBool(force),
		},
	}

	t, err := NewTemplate(tmpl, templateName, fargs)
	if err != nil {
		return fmt.Errorf("error parsing template %s: %w", templateName, err)
	}

	if err := t.ExecuteTemplate(gBuffer, templateName, nil); err != nil {
		return fmt.Errorf("error executing template %s: %w", templateName, err)
	}

	// Write our fBuffer to disk
	for path, fBuffer := range fBufMap {
		if err := WriteBuf(path, fBuffer); err != nil {
			logger.Panic(err)
		}
	}

	return nil
}

func WriteBuf(path string, fBuffer FileBuffer) error {
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

func ValidateTemplate(tmpl string) error {
	tmpl = filepath.Clean(tmpl)
	templateName := filepath.Base(tmpl)

	if templateName == "" {
		return errors.New("no template name provided")
	}

	re := regexp.MustCompile(`(?i)\.tmpl(\.gz)?$`)
	if !re.MatchString(templateName) {
		return fmt.Errorf("%q does not appear to be a valid template file (must end in `.tmpl` or `.tmpl.gz`)", templateName)
	}

	if _, err := FS.Stat(tmpl); err != nil {
		return fmt.Errorf("error validating template %s: %w", tmpl, err)
	}

	return nil
}

func NewTemplate(tmpl string, name string, fargs FuncArgs) (*template.Template, error) {
	var data []byte

	rawData, err := FS.ReadFile(tmpl)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(tmpl, ".gz") {
		gzipReader, err := gzip.NewReader(bytes.NewReader(rawData))
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := gzipReader.Close(); err != nil {
				logger.Fatal("error closing gzip reader: %w", err)
			}
		}()

		data, err = io.ReadAll(gzipReader)
		if err != nil {
			return nil, err
		}
	} else {
		data = rawData
	}

	return template.New(name).
		Funcs(template.FuncMap{
			"modlet":    ModletFunc(fargs),
			"mult":      MultFunc(fargs),
			"output":    OutputFunc(fargs),
			"set":       SetFunc(fargs),
			"write":     WriteFunc(fargs),
			"xmlHeader": func() string { return xml.Header },
		}).
		Parse(string(data))
}
