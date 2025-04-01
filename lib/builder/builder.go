package builder

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
)

type Buffer struct {
	buffer *bytes.Buffer
	writer io.WriteCloser
}

type BufferMap map[string]Buffer

func BuildModlet(templates []string, gamedir string, outdir string) error {
	var modInfo *modinfo.ModInfo
	bufmap := BufferMap{}
	buffer := &bytes.Buffer{}
	startCalled := false

	if len(templates) == 0 {
		return errors.New("no templates provided")
	}

	if gamedir == "" {
		return errors.New("gamedir not provided")
	}

	if outdir == "" {
		return errors.New("outdir not provided")
	}

	templateFile := filepath.Base(templates[0])

	modlet := func(path string) string {
		if path == "" {
			log.Fatal("modlet path not provided")
		}

		cleanPath := filepath.Clean(path)

		modInfo = modinfo.NewModInfo(filepath.Base(cleanPath))
		modInfo.SetPath(cleanPath)

		return ""
	}

	start := func() string {
		log.Println("Start function called")

		if modInfo == nil || modInfo.Path() == "" {
			log.Fatal("Please set the modlet path using {{ modlet <path> }}")
		}

		buffer.Reset()
		startCalled = true

		return ""
	}

	writeTo := func(path string) string {
		if path == "" {
			log.Fatal("output file not provided")
		}

		if !startCalled {
			log.Fatal("start function not called before writeTo")
		}

		cleanPath := filepath.Clean(path)
		fullPath := filepath.Join(modInfo.Path(), cleanPath)

		log.Printf("Output file: %s\n", fullPath)

		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			log.Fatalf("error creating modlet directory %s: %v", filepath.Dir(fullPath), err)
		}

		f, err := os.Create(fullPath)
		if err != nil {
			log.Fatalf("error creating output file %s: %v", fullPath, err)
		}

		bufmap[fullPath] = Buffer{
			buffer: &bytes.Buffer{},
			writer: f,
		}
		bufmap[fullPath].buffer.Write(buffer.Bytes())
		buffer.Reset()

		return ""
	}

	log.Printf("Processing template: %s\n", templateFile)

	t, err := template.New(templateFile).
		Funcs(template.FuncMap{
			"modlet":    modlet,
			"start":     start,
			"writeTo":   writeTo,
			"xmlHeader": func() string { return xml.Header },
		}).
		ParseFiles(templates...)
	if err != nil {
		return fmt.Errorf("error parsing template %s: %w", templateFile, err)
	}

	log.Printf("Executing %s\n", templateFile)
	if err := t.ExecuteTemplate(buffer, templateFile, nil); err != nil {
		return fmt.Errorf("error executing template %s: %w", templateFile, err)
	}

	for path, buffer := range bufmap {
		log.Printf("Writing to %s\n", path)

		if buffer.writer != nil {
			defer func() {
				if err := buffer.writer.Close(); err != nil {
					log.Fatalf("error closing output file: %v", err)
				}
			}()

			if _, err := buffer.buffer.WriteTo(buffer.writer); err != nil {
				return fmt.Errorf("error writing to output file %s: %w", path, err)
			}
		}
	}

	return nil
}
