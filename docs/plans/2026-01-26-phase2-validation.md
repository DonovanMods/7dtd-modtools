# Phase 2: Validation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Validate generated XPath expressions against game XML, reporting mismatches with diagnostic context to help modders catch stale or incorrect paths.

**Architecture:** Extend 7dtd-gamedata with raw XML loading and XPath query support. Create a validate package in modtools that extracts XPaths from generated configs, validates them against game XML, and produces summary reports with failure diagnostics.

**Tech Stack:** Go 1.25, antchfx/xmlquery, antchfx/xpath, 7dtd-gamedata library

---

## Design Decisions

| Decision            | Choice                      | Rationale                                        |
| ------------------- | --------------------------- | ------------------------------------------------ |
| Validation timing   | After all files written     | Consolidated report, cleaner output              |
| Output format       | Summary with counts         | Show pass/fail per file, list only failures      |
| XPath library       | antchfx/xpath + xmlquery    | Popular, well-maintained, full XPath 1.0 support |
| Default behavior    | On by default               | Encourages catching issues early                 |
| File mapping        | Convention-based            | Config/blocks.xml → gamedir/blocks.xml           |
| Game XML loading    | On demand via 7dtd-gamedata | Supports any game file, memory efficient         |
| Failure diagnostics | Parent exists check         | Shows which path segment failed                  |
| Similar matches     | Deferred to Phase 5         | YAGNI - keep initial implementation simple       |

---

## Architecture Overview

The validation system has two entry points:

1. **Integrated validation** - runs automatically after `modlet unpack` completes
2. **Standalone command** - `modlet validate <modlet-dir> [--gamedir <path>]`

Both share a core validation engine that:

- Parses generated modlet XML files to extract XPath expressions
- Maps each config file to its corresponding game XML file by naming convention
- Loads game XML files on demand via 7dtd-gamedata
- Executes each XPath against the game XML using `antchfx/xmlquery`
- Collects results and produces a summary report with diagnostics

**Package structure:**

```
# In 7dtd-gamedata
gamexml/
  loader.go        # XMLDocument, GameXML types with XPath support
  loader_test.go

# In 7dtd-modtools
validate/
  extract.go       # XPath extraction from modlet XML
  validate.go      # Core validation engine
  diagnose.go      # Failure diagnostics
  report.go        # Result formatting and output
  validate_test.go
```

---

## Critical Context

### Repository Locations

| Repository        | Path                                            | Purpose                      |
| ----------------- | ----------------------------------------------- | ---------------------------- |
| **7dtd-gamedata** | `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata` | Library for parsing game XML |
| **7dtd-modtools** | `/home/dyoung/Projects/mods/7dtd/7dtd-modtools` | CLI tool (this project)      |

**IMPORTANT:** Tasks 1-2 work in `7dtd-gamedata`. Tasks 3+ work in `7dtd-modtools`.

### Files to Read Before Starting

**In 7dtd-gamedata:**

- `gamexml/blocks.go` - Existing XML parsing patterns
- `gamexml/gamexml.go` - Common types

**In 7dtd-modtools:**

- `modlet/functions.go` - See how XPath instructions are generated (mkSet)
- `cmd/modlet/modlet.go` - Persistent flags pattern (gamedir)
- `cmd/modlet/unpack/unpack.go` - Existing unpack command structure
- `modlet/tests/modlet_test.go` - Test setup patterns

### New Dependency

Add `antchfx/xmlquery` to both repositories:

```bash
go get github.com/antchfx/xmlquery
```

This provides XPath query support over parsed XML documents.

---

## Task 1: Add XMLDocument and GameXML Loader to 7dtd-gamedata

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata`

**Files:**

- Create: `gamexml/loader.go`
- Create: `gamexml/loader_test.go`

**Step 1: Write the failing test**

```go
// gamexml/loader_test.go
package gamexml_test

import (
	"testing"

	"github.com/donovanmods/7dtd-gamedata/gamexml"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGameXML_Load(t *testing.T) {
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}
	gamexml.FS = fs

	gamedir := "/game/Config"
	require.NoError(t, fs.MkdirAll(gamedir, 0755))

	blocksXML := `<?xml version="1.0" encoding="UTF-8"?>
<blocks>
	<block name="terrStone">
		<property name="Material" value="Mite"/>
	</block>
	<block name="terrDirt">
		<property name="Material" value="Mdirt"/>
	</block>
</blocks>`

	require.NoError(t, fs.WriteFile(gamedir+"/blocks.xml", []byte(blocksXML), 0644))

	gx := gamexml.NewGameXML(gamedir)

	doc, err := gx.Load("blocks.xml")
	require.NoError(t, err)
	assert.Equal(t, "blocks.xml", doc.Filename)
	assert.NotNil(t, doc.Root)

	// Loading same file again should return cached version
	doc2, err := gx.Load("blocks.xml")
	require.NoError(t, err)
	assert.Same(t, doc, doc2)
}

func TestGameXML_LoadMissing(t *testing.T) {
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}
	gamexml.FS = fs

	gx := gamexml.NewGameXML("/nonexistent")

	_, err := gx.Load("blocks.xml")
	assert.Error(t, err)
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go test -v -run "TestGameXML" ./gamexml/...
```

Expected: FAIL with undefined types

**Step 3: Write minimal implementation**

```go
// gamexml/loader.go
package gamexml

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/antchfx/xmlquery"
)

// XMLDocument holds a parsed XML file with its root node for XPath queries
type XMLDocument struct {
	Filename string
	Root     *xmlquery.Node
}

// GameXML manages loading and caching of game XML files
type GameXML struct {
	gamedir   string
	documents map[string]*XMLDocument
}

// NewGameXML creates a new GameXML loader for the given game directory
func NewGameXML(gamedir string) *GameXML {
	return &GameXML{
		gamedir:   gamedir,
		documents: make(map[string]*XMLDocument),
	}
}

// Load loads an XML file by name, returning cached version if available
func (g *GameXML) Load(filename string) (*XMLDocument, error) {
	// Check cache first
	if doc, ok := g.documents[filename]; ok {
		return doc, nil
	}

	// Load from filesystem
	path := filepath.Join(g.gamedir, filename)
	data, err := FS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", filename, err)
	}

	// Parse XML
	root, err := xmlquery.Parse(strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", filename, err)
	}

	doc := &XMLDocument{
		Filename: filename,
		Root:     root,
	}

	// Cache and return
	g.documents[filename] = doc
	return doc, nil
}
```

**Step 4: Add go dependency**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go get github.com/antchfx/xmlquery
go mod tidy
```

**Step 5: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go test -v -run "TestGameXML" ./gamexml/...
```

Expected: PASS

**Step 6: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
git add gamexml/loader.go gamexml/loader_test.go go.mod go.sum
git commit -m "feat(gamexml): add XMLDocument and GameXML loader for XPath queries"
```

---

## Task 2: Add XPath Query Method to GameXML

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata`

**Files:**

- Modify: `gamexml/loader.go`
- Modify: `gamexml/loader_test.go`

**Step 1: Write the failing test**

Add to `gamexml/loader_test.go`:

```go
func TestGameXML_QueryXPath(t *testing.T) {
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}
	gamexml.FS = fs

	gamedir := "/game/Config"
	require.NoError(t, fs.MkdirAll(gamedir, 0755))

	blocksXML := `<?xml version="1.0" encoding="UTF-8"?>
<blocks>
	<block name="terrStone">
		<property name="Material" value="Mite"/>
		<drop event="Harvest" name="resourceRockSmall" count="55"/>
	</block>
	<block name="terrDirt">
		<property name="Material" value="Mdirt"/>
	</block>
</blocks>`

	require.NoError(t, fs.WriteFile(gamedir+"/blocks.xml", []byte(blocksXML), 0644))

	gx := gamexml.NewGameXML(gamedir)

	// Query that matches
	nodes, err := gx.QueryXPath("blocks.xml", "//block[@name='terrStone']")
	require.NoError(t, err)
	assert.Len(t, nodes, 1)

	// Query that matches multiple
	nodes, err = gx.QueryXPath("blocks.xml", "//block")
	require.NoError(t, err)
	assert.Len(t, nodes, 2)

	// Query that matches none
	nodes, err = gx.QueryXPath("blocks.xml", "//block[@name='nonexistent']")
	require.NoError(t, err)
	assert.Len(t, nodes, 0)

	// Query with nested path
	nodes, err = gx.QueryXPath("blocks.xml", "//block[@name='terrStone']/drop[@event='Harvest']/@count")
	require.NoError(t, err)
	assert.Len(t, nodes, 1)
}

func TestGameXML_QueryXPath_InvalidSyntax(t *testing.T) {
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}
	gamexml.FS = fs

	gamedir := "/game/Config"
	require.NoError(t, fs.MkdirAll(gamedir, 0755))
	require.NoError(t, fs.WriteFile(gamedir+"/blocks.xml", []byte(`<blocks></blocks>`), 0644))

	gx := gamexml.NewGameXML(gamedir)

	// Invalid XPath syntax
	_, err := gx.QueryXPath("blocks.xml", "//[invalid")
	assert.Error(t, err)
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go test -v -run "TestGameXML_QueryXPath" ./gamexml/...
```

Expected: FAIL with undefined method

**Step 3: Add QueryXPath method**

Add to `gamexml/loader.go`:

```go
// QueryXPath executes an XPath query against the specified file
// Returns matching nodes, or empty slice if no matches
func (g *GameXML) QueryXPath(filename, xpath string) ([]*xmlquery.Node, error) {
	doc, err := g.Load(filename)
	if err != nil {
		return nil, err
	}

	nodes, err := xmlquery.QueryAll(doc.Root, xpath)
	if err != nil {
		return nil, fmt.Errorf("invalid xpath %q: %w", xpath, err)
	}

	return nodes, nil
}
```

**Step 4: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go test -v -run "TestGameXML_QueryXPath" ./gamexml/...
```

Expected: PASS

**Step 5: Run all tests**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go test -v ./...
```

Expected: All tests PASS

**Step 6: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
git add gamexml/loader.go gamexml/loader_test.go
git commit -m "feat(gamexml): add QueryXPath method for XPath validation"
```

---

## CHECKPOINT A: Verify 7dtd-gamedata

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata`

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go test -v ./...
```

**Expected:** All tests PASS

**If any test fails:** STOP and diagnose before proceeding to Task 3.

---

## Task 3: Create XPath Extraction Package

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Files:**

- Create: `validate/extract.go`
- Create: `validate/extract_test.go`

**Step 1: Write the failing test**

```go
// validate/extract_test.go
package validate_test

import (
	"testing"

	"github.com/donovanmods/7dtd-modtools/validate"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var FS = &afero.Afero{Fs: afero.NewMemMapFs()}

func TestExtractXPaths(t *testing.T) {
	validate.FS = FS

	configXML := `<?xml version="1.0" encoding="UTF-8"?>
<configs>
	<set xpath="//block[@name='terrStone']/drop/@count">82</set>
	<append xpath="//block[@name='terrDirt']"><property name="New" value="1"/></append>
	<remove xpath="//block[@name='oldBlock']"/>
</configs>`

	path := "/test/Config/blocks.xml"
	require.NoError(t, FS.MkdirAll("/test/Config", 0755))
	require.NoError(t, FS.WriteFile(path, []byte(configXML), 0644))

	instructions, err := validate.ExtractXPaths(path)
	require.NoError(t, err)

	assert.Len(t, instructions, 3)

	assert.Equal(t, "//block[@name='terrStone']/drop/@count", instructions[0].XPath)
	assert.Equal(t, "set", instructions[0].Instruction)

	assert.Equal(t, "//block[@name='terrDirt']", instructions[1].XPath)
	assert.Equal(t, "append", instructions[1].Instruction)

	assert.Equal(t, "//block[@name='oldBlock']", instructions[2].XPath)
	assert.Equal(t, "remove", instructions[2].Instruction)
}

func TestExtractXPaths_Empty(t *testing.T) {
	validate.FS = FS

	configXML := `<?xml version="1.0" encoding="UTF-8"?>
<configs>
</configs>`

	path := "/test/empty.xml"
	require.NoError(t, FS.WriteFile(path, []byte(configXML), 0644))

	instructions, err := validate.ExtractXPaths(path)
	require.NoError(t, err)
	assert.Len(t, instructions, 0)
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run "TestExtractXPaths" ./validate/...
```

Expected: FAIL with package not found

**Step 3: Write minimal implementation**

```go
// validate/extract.go
package validate

import (
	"encoding/xml"
	"fmt"

	"github.com/spf13/afero"
)

// FS is the filesystem to use (can be swapped for testing)
var FS = &afero.Afero{Fs: afero.NewOsFs()}

// XPathInstruction represents an extracted XPath instruction from a modlet config
type XPathInstruction struct {
	XPath       string
	Instruction string // "set", "append", "remove", etc.
	Line        int
	File        string
}

// instructionElement represents an XML element with an xpath attribute
type instructionElement struct {
	XMLName xml.Name
	XPath   string `xml:"xpath,attr"`
}

// configsWrapper wraps the configs element to extract instructions
type configsWrapper struct {
	XMLName  xml.Name `xml:"configs"`
	Children []instructionElement `xml:",any"`
}

// ExtractXPaths parses a modlet config file and extracts all XPath instructions
func ExtractXPaths(configPath string) ([]XPathInstruction, error) {
	data, err := FS.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var configs configsWrapper
	if err := xml.Unmarshal(data, &configs); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	var instructions []XPathInstruction
	for _, child := range configs.Children {
		if child.XPath == "" {
			continue
		}
		instructions = append(instructions, XPathInstruction{
			XPath:       child.XPath,
			Instruction: child.XMLName.Local,
			File:        configPath,
		})
	}

	return instructions, nil
}
```

**Step 4: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run "TestExtractXPaths" ./validate/...
```

Expected: PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
make check
git add validate/
git commit -m "feat(validate): add XPath extraction from modlet config files"
```

---

## Task 4: Add Line Number Tracking to Extraction

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Files:**

- Modify: `validate/extract.go`
- Modify: `validate/extract_test.go`

**Step 1: Write the failing test**

Add to `validate/extract_test.go`:

```go
func TestExtractXPaths_LineNumbers(t *testing.T) {
	validate.FS = FS

	configXML := `<?xml version="1.0" encoding="UTF-8"?>
<configs>
	<set xpath="//block[@name='first']">1</set>
	<set xpath="//block[@name='second']">2</set>
	<set xpath="//block[@name='third']">3</set>
</configs>`

	path := "/test/lines.xml"
	require.NoError(t, FS.WriteFile(path, []byte(configXML), 0644))

	instructions, err := validate.ExtractXPaths(path)
	require.NoError(t, err)

	require.Len(t, instructions, 3)
	assert.Equal(t, 3, instructions[0].Line)
	assert.Equal(t, 4, instructions[1].Line)
	assert.Equal(t, 5, instructions[2].Line)
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run "TestExtractXPaths_LineNumbers" ./validate/...
```

Expected: FAIL (Line is 0)

**Step 3: Update implementation with line tracking**

Replace `ExtractXPaths` in `validate/extract.go`:

```go
// ExtractXPaths parses a modlet config file and extracts all XPath instructions
func ExtractXPaths(configPath string) ([]XPathInstruction, error) {
	data, err := FS.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	var instructions []XPathInstruction
	var inConfigs bool

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parsing config: %w", err)
		}

		switch elem := token.(type) {
		case xml.StartElement:
			if elem.Name.Local == "configs" {
				inConfigs = true
				continue
			}

			if !inConfigs {
				continue
			}

			// Look for xpath attribute
			var xpath string
			for _, attr := range elem.Attr {
				if attr.Name.Local == "xpath" {
					xpath = attr.Value
					break
				}
			}

			if xpath != "" {
				line, _ := decoder.InputPos()
				instructions = append(instructions, XPathInstruction{
					XPath:       xpath,
					Instruction: elem.Name.Local,
					Line:        line,
					File:        configPath,
				})
			}

		case xml.EndElement:
			if elem.Name.Local == "configs" {
				inConfigs = false
			}
		}
	}

	return instructions, nil
}
```

Add imports to `validate/extract.go`:

```go
import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/spf13/afero"
)
```

**Step 4: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run "TestExtractXPaths" ./validate/...
```

Expected: All extraction tests PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
make check
git add validate/extract.go validate/extract_test.go
git commit -m "feat(validate): add line number tracking to XPath extraction"
```

---

## Task 5: Implement Core Validation Logic

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Files:**

- Create: `validate/validate.go`
- Modify: `validate/extract_test.go` (add validation tests)

**Step 1: Update go.mod to use latest gamedata**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go mod tidy
```

**Step 2: Write the failing test**

Add to `validate/extract_test.go` (or create `validate/validate_test.go`):

```go
func TestValidateXPath(t *testing.T) {
	validate.FS = FS

	// Create game XML
	gamedir := "/game/Config"
	require.NoError(t, FS.MkdirAll(gamedir, 0755))

	blocksXML := `<?xml version="1.0" encoding="UTF-8"?>
<blocks>
	<block name="terrStone">
		<property name="Material" value="Mite"/>
	</block>
</blocks>`
	require.NoError(t, FS.WriteFile(gamedir+"/blocks.xml", []byte(blocksXML), 0644))

	gx := gamexml.NewGameXML(gamedir)
	gamexml.FS = FS

	// Valid XPath
	result := validate.ValidateXPath(validate.XPathInstruction{
		XPath:       "//block[@name='terrStone']",
		Instruction: "set",
		Line:        1,
		File:        "Config/blocks.xml",
	}, gx, "blocks.xml")

	assert.True(t, result.Valid)
	assert.Equal(t, 1, result.MatchCount)
	assert.NoError(t, result.Error)

	// Invalid XPath (no match)
	result = validate.ValidateXPath(validate.XPathInstruction{
		XPath:       "//block[@name='nonexistent']",
		Instruction: "set",
		Line:        2,
		File:        "Config/blocks.xml",
	}, gx, "blocks.xml")

	assert.False(t, result.Valid)
	assert.Equal(t, 0, result.MatchCount)
}
```

Add import for gamexml:

```go
import (
	"github.com/donovanmods/7dtd-gamedata/gamexml"
)
```

**Step 3: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run "TestValidateXPath" ./validate/...
```

Expected: FAIL with undefined function

**Step 4: Write implementation**

```go
// validate/validate.go
package validate

import (
	"github.com/donovanmods/7dtd-gamedata/gamexml"
)

// ValidationResult holds the result of validating a single XPath
type ValidationResult struct {
	Instruction XPathInstruction
	Valid       bool
	MatchCount  int
	Error       error
	Diagnostic  *Diagnostic
}

// Diagnostic provides context about why an XPath failed
type Diagnostic struct {
	ParentExists bool
	ParentPath   string
	MissingPart  string
}

// ValidateXPath checks if an XPath matches any elements in the game XML
func ValidateXPath(inst XPathInstruction, gx *gamexml.GameXML, gameFile string) ValidationResult {
	result := ValidationResult{
		Instruction: inst,
	}

	nodes, err := gx.QueryXPath(gameFile, inst.XPath)
	if err != nil {
		result.Error = err
		return result
	}

	result.MatchCount = len(nodes)
	result.Valid = len(nodes) > 0

	return result
}
```

**Step 5: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run "TestValidateXPath" ./validate/...
```

Expected: PASS

**Step 6: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
make check
git add validate/
git commit -m "feat(validate): add core XPath validation logic"
```

---

## Task 6: Add Failure Diagnostics

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Files:**

- Create: `validate/diagnose.go`
- Create: `validate/diagnose_test.go`
- Modify: `validate/validate.go`

**Step 1: Write the failing test**

```go
// validate/diagnose_test.go
package validate_test

import (
	"testing"

	"github.com/donovanmods/7dtd-gamedata/gamexml"
	"github.com/donovanmods/7dtd-modtools/validate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiagnoseFailure_ParentExists(t *testing.T) {
	validate.FS = FS
	gamexml.FS = FS

	gamedir := "/game/Config"
	require.NoError(t, FS.MkdirAll(gamedir, 0755))

	blocksXML := `<?xml version="1.0" encoding="UTF-8"?>
<blocks>
	<block name="terrStone">
		<property name="Material" value="Mite"/>
	</block>
</blocks>`
	require.NoError(t, FS.WriteFile(gamedir+"/blocks.xml", []byte(blocksXML), 0644))

	gx := gamexml.NewGameXML(gamedir)

	// Block exists but property doesn't
	diag := validate.DiagnoseFailure(gx, "blocks.xml", "//block[@name='terrStone']/property[@name='NonExistent']/@value")

	assert.True(t, diag.ParentExists)
	assert.Contains(t, diag.ParentPath, "terrStone")
	assert.Contains(t, diag.MissingPart, "NonExistent")
}

func TestDiagnoseFailure_RootMissing(t *testing.T) {
	validate.FS = FS
	gamexml.FS = FS

	gamedir := "/game/Config"
	require.NoError(t, FS.MkdirAll(gamedir, 0755))

	blocksXML := `<?xml version="1.0" encoding="UTF-8"?>
<blocks>
	<block name="terrStone"/>
</blocks>`
	require.NoError(t, FS.WriteFile(gamedir+"/blocks.xml", []byte(blocksXML), 0644))

	gx := gamexml.NewGameXML(gamedir)

	// Block doesn't exist at all
	diag := validate.DiagnoseFailure(gx, "blocks.xml", "//block[@name='nonexistent']/drop/@count")

	assert.False(t, diag.ParentExists)
	assert.Contains(t, diag.MissingPart, "nonexistent")
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run "TestDiagnoseFailure" ./validate/...
```

Expected: FAIL with undefined function

**Step 3: Write implementation**

```go
// validate/diagnose.go
package validate

import (
	"regexp"
	"strings"

	"github.com/donovanmods/7dtd-gamedata/gamexml"
)

// DiagnoseFailure analyzes why an XPath didn't match and provides context
func DiagnoseFailure(gx *gamexml.GameXML, gameFile, xpath string) Diagnostic {
	diag := Diagnostic{}

	// Split xpath into segments
	segments := splitXPath(xpath)
	if len(segments) == 0 {
		diag.MissingPart = xpath
		return diag
	}

	// Try progressively shorter paths until one matches
	for i := len(segments) - 1; i >= 0; i-- {
		partialPath := strings.Join(segments[:i+1], "")
		nodes, err := gx.QueryXPath(gameFile, partialPath)
		if err != nil {
			continue
		}

		if len(nodes) > 0 {
			// Found where it stops matching
			diag.ParentExists = true
			diag.ParentPath = partialPath

			if i+1 < len(segments) {
				diag.MissingPart = extractNameFromSegment(segments[i+1])
			}
			return diag
		}
	}

	// Nothing matched - extract the first element name
	diag.ParentExists = false
	diag.MissingPart = extractNameFromSegment(segments[0])

	return diag
}

// splitXPath splits an XPath into logical segments
func splitXPath(xpath string) []string {
	// Split on / but keep the delimiter with the following segment
	var segments []string
	current := ""

	for i := 0; i < len(xpath); i++ {
		if xpath[i] == '/' {
			if current != "" {
				segments = append(segments, current)
			}
			current = "/"
		} else {
			current += string(xpath[i])
		}
	}

	if current != "" && current != "/" {
		segments = append(segments, current)
	}

	return segments
}

// extractNameFromSegment extracts the name attribute value from an XPath segment
func extractNameFromSegment(segment string) string {
	// Look for @name='value' pattern
	re := regexp.MustCompile(`@name=['"](.*?)['"]`)
	matches := re.FindStringSubmatch(segment)
	if len(matches) > 1 {
		return matches[1]
	}

	// Fall back to element name
	re = regexp.MustCompile(`^/?/?(\w+)`)
	matches = re.FindStringSubmatch(segment)
	if len(matches) > 1 {
		return matches[1]
	}

	return segment
}
```

**Step 4: Update ValidateXPath to include diagnostics**

In `validate/validate.go`, update `ValidateXPath`:

```go
// ValidateXPath checks if an XPath matches any elements in the game XML
func ValidateXPath(inst XPathInstruction, gx *gamexml.GameXML, gameFile string) ValidationResult {
	result := ValidationResult{
		Instruction: inst,
	}

	nodes, err := gx.QueryXPath(gameFile, inst.XPath)
	if err != nil {
		result.Error = err
		return result
	}

	result.MatchCount = len(nodes)
	result.Valid = len(nodes) > 0

	// Add diagnostics for failures
	if !result.Valid {
		diag := DiagnoseFailure(gx, gameFile, inst.XPath)
		result.Diagnostic = &diag
	}

	return result
}
```

**Step 5: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run "TestDiagnoseFailure" ./validate/...
```

Expected: PASS

**Step 6: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
make check
git add validate/
git commit -m "feat(validate): add failure diagnostics with parent path analysis"
```

---

## Task 7: Implement Report Formatting

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Files:**

- Create: `validate/report.go`
- Create: `validate/report_test.go`

**Step 1: Write the failing test**

```go
// validate/report_test.go
package validate_test

import (
	"bytes"
	"testing"

	"github.com/donovanmods/7dtd-modtools/validate"
	"github.com/stretchr/testify/assert"
)

func TestReport_Print(t *testing.T) {
	report := &validate.Report{
		Files: []validate.FileReport{
			{
				ConfigFile: "Config/blocks.xml",
				GameFile:   "blocks.xml",
				Total:      10,
				Valid:      8,
				Warnings: []validate.ValidationResult{
					{
						Instruction: validate.XPathInstruction{
							XPath: "//block[@name='missing']",
							Line:  5,
						},
						Diagnostic: &validate.Diagnostic{
							ParentExists: false,
							MissingPart:  "missing",
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	report.Print(&buf)

	output := buf.String()
	assert.Contains(t, output, "Config/blocks.xml")
	assert.Contains(t, output, "8/10")
	assert.Contains(t, output, "Line 5")
	assert.Contains(t, output, "missing")
}

func TestReport_HasWarnings(t *testing.T) {
	report := &validate.Report{
		Files: []validate.FileReport{
			{Total: 10, Valid: 10},
		},
	}
	assert.False(t, report.HasWarnings())

	report.Files[0].Warnings = []validate.ValidationResult{{}}
	assert.True(t, report.HasWarnings())
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run "TestReport" ./validate/...
```

Expected: FAIL with undefined types

**Step 3: Write implementation**

```go
// validate/report.go
package validate

import (
	"fmt"
	"io"
	"time"
)

// FileReport holds validation results for a single config file
type FileReport struct {
	ConfigFile string
	GameFile   string
	Total      int
	Valid      int
	Warnings   []ValidationResult // zero matches
	Errors     []ValidationResult // syntax/load errors
	Skipped    bool               // game file not found
}

// Report holds the complete validation report
type Report struct {
	Files    []FileReport
	Duration time.Duration
}

// Print outputs the validation report
func (r *Report) Print(w io.Writer) {
	fmt.Fprintln(w, "Validation Results:")

	totalValid := 0
	totalCount := 0
	totalWarnings := 0
	totalErrors := 0

	for _, f := range r.Files {
		totalValid += f.Valid
		totalCount += f.Total
		totalWarnings += len(f.Warnings)
		totalErrors += len(f.Errors)

		if f.Skipped {
			fmt.Fprintf(w, "  ○ %s: skipped (%s not found in gamedir)\n", f.ConfigFile, f.GameFile)
			continue
		}

		status := "✓"
		if len(f.Warnings) > 0 || len(f.Errors) > 0 {
			status = "⚠"
		}

		fmt.Fprintf(w, "  %s %s: %d/%d XPaths valid\n", status, f.ConfigFile, f.Valid, f.Total)

		for _, warn := range f.Warnings {
			fmt.Fprintf(w, "    ⚠ Line %d: %s\n", warn.Instruction.Line, warn.Instruction.XPath)
			if warn.Diagnostic != nil {
				if warn.Diagnostic.ParentExists {
					fmt.Fprintf(w, "      → parent exists, but no '%s'\n", warn.Diagnostic.MissingPart)
				} else {
					fmt.Fprintf(w, "      → '%s' not found\n", warn.Diagnostic.MissingPart)
				}
			}
		}

		for _, err := range f.Errors {
			fmt.Fprintf(w, "    ✗ Line %d: %s\n", err.Instruction.Line, err.Instruction.XPath)
			fmt.Fprintf(w, "      → error: %v\n", err.Error)
		}
	}

	fmt.Fprintf(w, "\nSummary: %d/%d valid, %d warnings, %d errors\n",
		totalValid, totalCount, totalWarnings, totalErrors)
}

// HasWarnings returns true if any file has validation warnings
func (r *Report) HasWarnings() bool {
	for _, f := range r.Files {
		if len(f.Warnings) > 0 {
			return true
		}
	}
	return false
}

// HasErrors returns true if any file has validation errors
func (r *Report) HasErrors() bool {
	for _, f := range r.Files {
		if len(f.Errors) > 0 {
			return true
		}
	}
	return false
}
```

**Step 4: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run "TestReport" ./validate/...
```

Expected: PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
make check
git add validate/
git commit -m "feat(validate): add report formatting with summary output"
```

---

## CHECKPOINT B: Verify Validation Engine

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v ./validate/...
go test -v ./...
```

**Expected:** All tests PASS

**If any test fails:** STOP and diagnose before proceeding to Task 8.

---

## Task 8: Add modlet validate CLI Command

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Files:**

- Create: `cmd/modlet/validate/validate.go`
- Modify: `cmd/modlet/modlet.go` (add subcommand)

**Step 1: Create the validate command**

```go
// cmd/modlet/validate/validate.go
package validate

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/donovanmods/7dtd-gamedata/gamexml"
	"github.com/donovanmods/7dtd-modtools/lib/logger"
	"github.com/donovanmods/7dtd-modtools/validate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var strict bool

var ValidateCmd = &cobra.Command{
	Use:   "validate <modlet-dir>",
	Short: "Validate modlet XPath expressions against game XML",
	Long: `Validate checks that all XPath expressions in a modlet's config files
match elements in the game XML files. This helps catch stale or incorrect paths
that won't have any effect when the mod is loaded.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		modletDir := args[0]
		gamedir := viper.GetString("gamedir")

		if gamedir == "" {
			logger.Fatal("gamedir is required (set via --gamedir or config)")
		}

		// Find all Config/*.xml files
		configDir := filepath.Join(modletDir, "Config")
		entries, err := os.ReadDir(configDir)
		if err != nil {
			logger.Fatal("reading Config directory: %v", err)
		}

		gx := gamexml.NewGameXML(gamedir)
		report := &validate.Report{}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".xml") {
				continue
			}

			configPath := filepath.Join(configDir, entry.Name())
			gameFile := entry.Name() // Convention: Config/blocks.xml -> blocks.xml

			fileReport := validateFile(configPath, gameFile, gx)
			report.Files = append(report.Files, fileReport)
		}

		report.Print(os.Stdout)

		if report.HasErrors() {
			os.Exit(1)
		}
		if strict && report.HasWarnings() {
			os.Exit(1)
		}
	},
}

func validateFile(configPath, gameFile string, gx *gamexml.GameXML) validate.FileReport {
	report := validate.FileReport{
		ConfigFile: configPath,
		GameFile:   gameFile,
	}

	// Check if game file exists
	_, err := gx.Load(gameFile)
	if err != nil {
		report.Skipped = true
		return report
	}

	instructions, err := validate.ExtractXPaths(configPath)
	if err != nil {
		logger.Error("extracting XPaths from %s: %v", configPath, err)
		return report
	}

	report.Total = len(instructions)

	for _, inst := range instructions {
		result := validate.ValidateXPath(inst, gx, gameFile)

		if result.Error != nil {
			report.Errors = append(report.Errors, result)
		} else if !result.Valid {
			report.Warnings = append(report.Warnings, result)
		} else {
			report.Valid++
		}
	}

	return report
}

func init() {
	ValidateCmd.Flags().BoolVar(&strict, "strict", false, "Treat warnings as errors")
}
```

**Step 2: Register command in modlet.go**

In `cmd/modlet/modlet.go`, add:

```go
import (
	// ... existing imports
	validateCmd "github.com/donovanmods/7dtd-modtools/cmd/modlet/validate"
)

func init() {
	// ... existing init code
	ModletCmd.AddCommand(validateCmd.ValidateCmd)
}
```

**Step 3: Verify it builds**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go build ./...
```

Expected: Build succeeds

**Step 4: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
make check
git add cmd/modlet/validate/ cmd/modlet/modlet.go
git commit -m "feat(cmd): add modlet validate command"
```

---

## Task 9: Integrate Validation into Unpack Command

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Files:**

- Modify: `cmd/modlet/unpack/unpack.go`
- Modify: `modlet/unpack.go` (add validation call)

**Step 1: Add flags to unpack command**

In `cmd/modlet/unpack/unpack.go`, add:

```go
var (
	noValidate bool
	strict     bool
)

func init() {
	// ... existing flag setup
	UnpackCmd.Flags().BoolVar(&noValidate, "no-validate", false, "Skip XPath validation")
	UnpackCmd.Flags().BoolVar(&strict, "strict", false, "Treat validation warnings as errors")
}
```

**Step 2: Add validation after unpack in the Run function**

Update the Run function to call validation after unpacking completes:

```go
Run: func(cmd *cobra.Command, args []string) {
	// ... existing unpack logic ...

	// After all templates unpacked, run validation
	if !noValidate && gamedir != "" {
		// Run validation on generated modlet
		// (implementation details depend on how modlet output path is tracked)
	}
}
```

**Step 3: Create validation integration function**

Add a function to run validation after unpack:

```go
// In modlet/unpack.go or validate/integrate.go

func ValidateAfterUnpack(modletDir, gamedir string, strict bool) error {
	// Find Config/*.xml files
	// Run validation
	// Print report
	// Return error if strict and warnings, or if errors
}
```

**Step 4: Run tests**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v ./...
```

Expected: All tests PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
make check
git add cmd/modlet/unpack/ modlet/
git commit -m "feat(unpack): integrate validation with --no-validate and --strict flags"
```

---

## CHECKPOINT C: Final Verification

> **Verify both repositories**

**Step 1: Run all tests in 7dtd-gamedata**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go test -v ./...
```

**Step 2: Run all tests in 7dtd-modtools**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v ./...
```

**Step 3: Run lint checks**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
make check
```

**Step 4: Manual test**

```bash
# Build the tool
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go build -o bin/7dmt ./main.go

# Test validate command (if you have a modlet and game data available)
./bin/7dmt modlet validate /path/to/modlet --gamedir /path/to/game/Config
```

**Step 5: Final cleanup commit if needed**

```bash
git status
# If changes:
git add -A
git commit -m "chore: Phase 2 cleanup"
```

---

## Summary

Phase 2 establishes XPath validation:

1. **7dtd-gamedata enhancements** (Tasks 1-2):
   - `XMLDocument` type for raw XML access
   - `GameXML` loader with caching
   - `QueryXPath` method for XPath execution

2. **Validation engine** (Tasks 3-7):
   - XPath extraction from modlet configs with line numbers
   - Core validation logic using gamexml
   - Failure diagnostics (parent path analysis)
   - Report formatting with summary

3. **CLI integration** (Tasks 8-9):
   - `modlet validate <dir>` command
   - Validation integrated into `modlet unpack`
   - `--no-validate` and `--strict` flags

**Modlets can now be validated to catch:**

- Stale XPaths referencing removed game elements
- Typos in element or attribute names
- Incorrect path structures
