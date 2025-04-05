package builder_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"testing"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
	"github.com/donovanmods/7dtd-modtools/lib/builder"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var (
	testTMP = "test_temp"
	FS      = &afero.Afero{Fs: afero.NewMemMapFs()}
)

type nopIO struct {
	io.Reader
	io.Writer
}

func (nopIO) Close() error { return nil }

func NewNopIO(t *testing.T, buf *bytes.Buffer) io.WriteCloser {
	t.Helper()

	return nopIO{buf, buf}
}

func mkTempDir(t *testing.T) string {
	t.Helper()

	if exists, err := afero.DirExists(FS, testTMP); err != nil {
		t.Fatalf("Error checking for temp directory: %v", err)
	} else if !exists {
		err := FS.Mkdir(testTMP, 0700)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("Created temp directory: %s", testTMP)
	}

	return testTMP
}

func setup(t *testing.T) *assert.Assertions {
	t.Helper()

	log.Print("running setup")

	// Use MemMapFs for testing
	builder.FS = FS

	mkTempDir(t)

	return assert.New(t)
}

func TestFuncModlet(t *testing.T) {
	assert := setup(t)

	funcArgs := builder.FuncArgs{
		Outdir:  testTMP,
		ModInfo: &modinfo.ModInfo{},
	}

	fn := builder.FuncModlet(funcArgs)
	modletName := "TestModlet"

	fn(modletName)

	assert.Equal(modletName, funcArgs.ModInfo.GetValue("name"))
	assert.Equal(filepath.Join(funcArgs.Outdir, modletName), funcArgs.ModInfo.Path())

	exists, err := FS.DirExists(funcArgs.ModInfo.Path())
	assert.NoError(err, "Error checking for Modlet directory")
	assert.True(exists, "Modlet directory should exist")
}

func TestFuncOutput(t *testing.T) {
	assert := setup(t)

	funcArgs := builder.FuncArgs{
		Outdir:  testTMP,
		ModInfo: &modinfo.ModInfo{},
		FBuffer: &builder.FileBuffer{},
		GBuffer: bytes.NewBuffer(nil),
		FBufMap: make(builder.FileBufferMap),
	}

	funcArgs.ModInfo.SetPath(funcArgs.Outdir)

	fn := builder.FuncOutput(funcArgs)
	outputPath := "testfile.txt"

	fn(outputPath)

	fullPath := filepath.Join(funcArgs.Outdir, outputPath)
	assert.Contains(funcArgs.FBufMap, fullPath, "File buffer map should contain the output path")

	exists, err := FS.Exists(funcArgs.ModInfo.Path())
	assert.NoError(err, "Error checking for output file")
	assert.True(exists, "Output file should exist")
}

func TestFuncSet(t *testing.T) {
	assert := setup(t)

	funcArgs := builder.FuncArgs{
		ModInfo: &modinfo.ModInfo{},
	}

	fn := builder.FuncSet(funcArgs)

	xpath := `//block[@name='terrStone']/drop[@event='Harvest' and @name='resourceRockSmall']/@count`
	value := "999"

	expected := fmt.Sprintf("<set xpath=\"%s\">%s</set>", xpath, value)
	actual := fn(xpath, value)

	assert.Equal(expected, actual, "FuncSet should return the correct XML string")
}

func TestFuncWrite(t *testing.T) {
	assert := setup(t)

	buf := bytes.NewBuffer(nil)

	funcArgs := builder.FuncArgs{
		FBuffer: &builder.FileBuffer{
			Buffer: buf,
			Writer: NewNopIO(t, buf), // Using stdout for testing
		},
		GBuffer: bytes.NewBufferString("Test content"),
	}

	fn := builder.FuncWrite(funcArgs)
	fn()

	assert.Equal("Test content", funcArgs.FBuffer.Buffer.String(), "Buffer content should match")
}
