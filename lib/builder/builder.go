package builder

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
)

type fileBuffer struct {
	buffer *bytes.Buffer
	writer io.WriteCloser
}

type null string // we use this when we need to output nothing

// Holds buffers and io.WriteCloser for each output file
type fileBufferMap map[string]fileBuffer

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

	fBufMap := make(fileBufferMap)
	gBuffer := bytes.NewBuffer(nil)
	fBuffer := &fileBuffer{
		buffer: gBuffer,
		writer: nil,
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

	log.Printf("processing template: %s\n", templateName)

	t, err := template.New(templateName).
		Funcs(template.FuncMap{
			"modlet":    modletFunc(outdir, &modInfo),
			"output":    outputFunc(fBuffer, gBuffer, fBufMap, &modInfo),
			"write":     writeFunc(fBuffer, gBuffer),
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
			log.Fatal(err)
		}
	}

	return nil
}

func mkPath(path string) error {
	if !fs.ValidPath(path) {
		return fmt.Errorf("invalid path %q", path)
	}

	_, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		log.Printf("creating directory: %q", path)

		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("unable to create directory %q: %w", path, err)
		}
	}

	return nil
}

func writeBuf(path string, fBuffer fileBuffer) error {
	log.Printf("writing %q\n", path)

	if fBuffer.writer != nil {
		defer func() {
			if err := fBuffer.writer.Close(); err != nil {
				log.Fatalf("error closing output file: %v", err)
			}
		}()

		if _, err := fBuffer.buffer.WriteTo(fBuffer.writer); err != nil {
			return fmt.Errorf("error writing to output file %s: %w", path, err)
		}
	}

	return nil
}

func modletFunc(outdir string, modInfo *modinfo.ModInfo) func(string) null {
	return func(name string) null {
		if name == "" {
			log.Fatal("modlet name must be provided")
		}

		path := filepath.Join(outdir, name)

		*modInfo = *modinfo.NewModInfo(name)
		modInfo.SetPath(path)

		if err := mkPath(modInfo.Path()); err != nil {
			log.Fatal(err)
		}

		log.Printf("creating modlet %q\n", modInfo.GetValue("name"))

		return null("")
	}
}

func outputFunc(fBuffer *fileBuffer, gBuffer *bytes.Buffer, fBufMap fileBufferMap, modInfo *modinfo.ModInfo) func(string) null {
	return func(path string) null {
		if path == "" {
			log.Fatal("output file not provided")
		}

		if modInfo.Path() == "" {
			log.Fatal("please set the modlet using {{ modlet <name> }}")
		}

		cleanPath := filepath.Clean(path)
		fullPath := filepath.Join(modInfo.Path(), cleanPath)

		log.Printf("buffering output for %q\n", fullPath)

		if err := mkPath(filepath.Dir(fullPath)); err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(fullPath)
		if err != nil {
			log.Fatalf("error creating output file %s: %v", fullPath, err)
		}

		gBuffer.Reset()

		fBufMap[fullPath] = fileBuffer{
			buffer: bytes.NewBuffer(nil),
			writer: f,
		}
		*fBuffer = fBufMap[fullPath]

		outputFound = true

		return null("")
	}
}

func writeFunc(fBuffer *fileBuffer, gBuffer *bytes.Buffer) func() null {
	return func() null {
		if !outputFound || (*fBuffer).writer == nil {
			log.Fatal("you've called `write` without providing an output file, please use `output <filepath>` before `write`")
		}

		log.Println("saving fileBuffer")

		// Copy the current buffer to the output buffer
		(*fBuffer).buffer.Write(gBuffer.Bytes())

		// Reset the output state
		outputFound = false

		return null("")
	}
}
