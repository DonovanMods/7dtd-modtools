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
package pack

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
	"github.com/donovanmods/7dtd-modtools/lib/logger"
	"github.com/donovanmods/7dtd-modtools/modlet/lib/common"
	"github.com/spf13/viper"
)

type BufferMap map[string]*bytes.Buffer

var (
	validExtensions []string = []string{".md", ".txt", ".xml"}
	modletName      string   = "modlet.tmpl"
	modletPath      string
	bufferMap       BufferMap = make(BufferMap)
)

func Run(moddir string, output string) error {
	moddir = filepath.Clean(moddir)
	output = filepath.Clean(output)

	logger.Info("Packing modlet from %s to %s", moddir, output)

	if isDir, _ := common.FS.DirExists(moddir); !isDir {
		return fmt.Errorf("modlet directory %s does not exist", moddir)
	}

	if err := walkDir(moddir); err != nil {
		return fmt.Errorf("error walking modlet directory: %w", err)
	}

	if err := common.FS.MkdirAll(output, 0755); err != nil {
		return fmt.Errorf("error creating output directory %s: %w", output, err)
	}

	if err := writeBufferMapToFile(output); err != nil {
		return fmt.Errorf("error writing modlet file: %w", err)
	}

	logger.Info("Modlet packed successfully to %s", modletPath)

	return nil
}

func walkDir(dir string) error {
	err := common.FS.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking modlet directory: %w", err)
		}

		if !info.IsDir() {

			logger.Trace("Walking %#v", path)

			if !slices.Contains(validExtensions, filepath.Ext(info.Name())) {
				logger.Warn("Skipping unsupported file: %s", path)
				return nil
			}

			relPath := stripPath(path, dir)

			logger.Info("Packing file: %s", relPath)

			if strings.EqualFold(info.Name(), "modinfo.xml") {
				miXML, err := modinfo.Parse(path)
				if err != nil {
					return fmt.Errorf("error parsing modinfo.xml: %w", err)
				}
				modletName = miXML.GetValue("name")
				if modletName == "" {
					return fmt.Errorf("modlet name not found in modinfo.xml")
				}
			}

			if err := mkFileBlock(path, relPath); err != nil {
				return fmt.Errorf("error adding section for %s: %w", relPath, err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking modlet directory: %w", err)
	}

	return nil
}

func mkFileBlock(path string, relpath string) error {
	buf := bytes.NewBuffer(nil)

	logger.Trace("creating file block for %q", relpath)

	_, err := fmt.Fprintf(buf, "{{- output %q -}}\n", relpath)
	if err != nil {
		return fmt.Errorf("error writing {{ output }} tag for %q: %w", relpath, err)
	}

	c, err := common.FS.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", path, err)
	}

	if _, err := buf.Write(c); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	_, err = fmt.Fprint(buf, "{{- write -}}\n")
	if err != nil {
		return fmt.Errorf("error writing {{ write }} tag for %q: %w", relpath, err)
	}

	bufferMap[relpath] = buf

	return nil
}

func writeOutput(key string, file io.Writer) error {
	if key == "" {
		return errors.New("empty key provided to writeToOutput")
	}

	if _, ok := bufferMap[key]; !ok {
		return fmt.Errorf("no buffer found for key %q", key)
	}

	logger.Debug("Writing %q to output buffer", key)

	if _, err := fmt.Fprint(file, bufferMap[key].String()); err != nil {
		return fmt.Errorf("error writing %s buffer to output: %w", key, err)
	}

	return nil
}

func writeBufferMapToFile(output string) error {
	var (
		// buffer = bytes.NewBuffer([]byte(fmt.Sprintf("{{- modlet %q -}}\n", modletName)))
		// buffer = bytes.NewBuffer(nil)
		gzWriter *gzip.Writer
		writer   *bufio.Writer
		compress = viper.GetBool("compress")
	)

	modletPath = filepath.Join(output, strings.ToLower(modletName+".tmpl"))

	if compress {
		modletPath += ".gz"
	}

	file, err := common.FS.Create(modletPath)
	if err != nil {
		return err
	}

	if compress {
		gzWriter = gzip.NewWriter(file)
		writer = bufio.NewWriter(gzWriter)
	} else {
		writer = bufio.NewWriter(file)
	}

	if exists, err := common.FS.Exists(modletPath); err == nil && exists {
		if !viper.GetBool("force") {
			return fmt.Errorf("modlet file %s already exists, use --force to overwrite", modletPath)
		}

		logger.Warn("Overwriting existing modlet file %s", modletPath)
	}

	// Write template modlet name
	logger.Debug("Writing modlet data to buffer for %q", modletName)

	if _, err := fmt.Fprintf(writer, "{{- modlet %q -}}\n", modletName); err != nil {
		return fmt.Errorf("error writing modlet name to buffer: %w", err)
	}

	// Write ModInfo data
	if err := writeOutput(findBufferKey("modinfo.xml"), writer); err != nil {
		return fmt.Errorf("error writing modinfo block: %w", err)
	}

	// Write the remaining buffers
	for key := range bufferMap {
		// Skip modinfo because it has already been written
		if strings.EqualFold(key, "modinfo.xml") {
			continue
		}

		if err := writeOutput(key, writer); err != nil {
			return fmt.Errorf("error writing modinfo block: %w", err)
		}
	}

	// Flush our buffers before writing to file
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	if compress {
		if err := gzWriter.Close(); err != nil {
			return fmt.Errorf("error closing gzip writer: %w", err)
		}
	}
	// if err := common.FS.WriteFile(modletPath, buffer.Bytes(), 0644); err != nil {
	// 	return fmt.Errorf("error writing modlet file %s: %w", filepath.Join(output, modletPath), err)
	// }

	if err := file.Close(); err != nil {
		return fmt.Errorf("error closing file %s: %w", modletPath, err)
	}

	return nil
}

func stripPath(path string, dir string) string {
	path = strings.TrimPrefix(path, dir)
	path = strings.TrimPrefix(path, string(filepath.Separator))

	return filepath.Clean(path)
}

func findBufferKey(s string) string {
	for k := range bufferMap {
		if strings.EqualFold(k, s) {
			return k
		}
	}

	return ""
}
