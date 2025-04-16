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
package unpack

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
	"github.com/donovanmods/7dtd-modtools/lib/logger"
	"github.com/donovanmods/7dtd-modtools/modlet/lib/common"
	"github.com/donovanmods/7dtd-modtools/modlet/lib/functions"
)

func Run(tmpl string, gamedir string, outdir string) error {
	var modInfo modinfo.ModInfo

	fBufMap := make(common.FileBufferMap)
	gBuffer := bytes.NewBuffer(nil)
	fBuffer := &common.FileBuffer{
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

	fargs := common.FuncArgs{Outdir: outdir, Gamedir: gamedir, ModInfo: &modInfo, FBuffer: fBuffer, GBuffer: gBuffer, FBufMap: fBufMap}
	t, err := template.New(templateName).
		Funcs(template.FuncMap{
			"modlet":    functions.ModletFunc(fargs),
			"mult":      functions.MultFunc(fargs),
			"output":    functions.OutputFunc(fargs),
			"set":       functions.SetFunc(fargs),
			"write":     functions.WriteFunc(fargs),
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

func writeBuf(path string, fBuffer common.FileBuffer) error {
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
