package builder_test

import (
	"bytes"
	"io"
	"log"
	"path/filepath"
	"testing"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
	"github.com/donovanmods/7dtd-modtools/lib/builder"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
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

var (
	testTMP = "test_temp"
	FS      = &afero.Afero{Fs: afero.NewMemMapFs()}
)

func mkTempDir(t *testing.T) {
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
}

type FunctionsTestSuite struct {
	suite.Suite
}

func (suite *FunctionsTestSuite) SetupSuite() {
	log.Printf("setting up test using %s", testTMP)

	builder.FS = FS

	mkTempDir(suite.T())
}

func (suite *FunctionsTestSuite) TestFuncModlet() {
	funcArgs := builder.FuncArgs{
		Outdir:  testTMP,
		ModInfo: &modinfo.ModInfo{},
	}

	fn := builder.FuncModlet(funcArgs)
	modletName := "TestModlet"

	fn(modletName)

	suite.Equal(modletName, funcArgs.ModInfo.GetValue("name"))
	suite.Equal(filepath.Join(funcArgs.Outdir, modletName), funcArgs.ModInfo.Path())
	exists, err := FS.DirExists(funcArgs.ModInfo.Path())
	suite.NoError(err, "Error checking for Modlet directory")
	suite.Assert().True(exists, "Modlet directory should exist")
}

func (suite *FunctionsTestSuite) TestFuncOutput() {
	bufIO := bytes.NewBuffer(nil)
	funcArgs := builder.FuncArgs{
		Outdir:        testTMP,
		ModInfo:       &modinfo.ModInfo{},
		FBuffer:       &builder.FileBuffer{},
		GBuffer:       bytes.NewBuffer(nil),
		FBufMap:       make(builder.FileBufferMap),
		IoReader:      bufIO,
		IoWriteCloser: NewNopIO(suite.T(), bufIO),
	}

	funcArgs.ModInfo.SetPath(funcArgs.Outdir)

	fn := builder.FuncOutput(funcArgs)
	outputPath := "testfile.txt"

	fn(outputPath)

	fullPath := filepath.Join(funcArgs.Outdir, outputPath)
	suite.Contains(funcArgs.FBufMap, fullPath, "File buffer map should contain the output path")
	exists, err := FS.Exists(funcArgs.ModInfo.Path())
	suite.NoError(err, "Error checking for output file")
	suite.Assert().True(exists, "Output file should exist")
}

func (suite *FunctionsTestSuite) TestFuncWrite() {
	funcArgs := builder.FuncArgs{
		FBuffer: &builder.FileBuffer{
			Buffer: bytes.NewBuffer(nil),
			Writer: NewNopIO(suite.T(), bytes.NewBuffer(nil)), // Using stdout for testing
		},
		GBuffer: bytes.NewBufferString("Test content"),
	}

	fn := builder.FuncWrite(funcArgs)
	fn()

	suite.Equal("Test content", funcArgs.FBuffer.Buffer.String(), "Buffer content should match")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFunctionRunSuite(t *testing.T) {
	suite.Run(t, new(FunctionsTestSuite))
}
