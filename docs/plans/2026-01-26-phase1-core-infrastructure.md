# Phase 1: Core Infrastructure Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Complete the template engine foundation with game data support, enabling templates to iterate over game XML and generate modlets.

**Architecture:** Extend 7dtd-gamedata library with parsing for game XML files (blocks, entityclasses, items). Create a GameData struct in modtools that loads all game data and passes it to template execution. Add helper functions for string matching and value manipulation.

**Tech Stack:** Go 1.25, text/template, encoding/xml, 7dtd-gamedata library

---

## Critical Context

### Repository Locations

| Repository        | Path                                            | Purpose                      |
| ----------------- | ----------------------------------------------- | ---------------------------- |
| **7dtd-gamedata** | `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata` | Library for parsing game XML |
| **7dtd-modtools** | `/home/dyoung/Projects/mods/7dtd/7dtd-modtools` | CLI tool (this project)      |

**IMPORTANT:** Tasks 1-3 work in `7dtd-gamedata`. Tasks 4+ work in `7dtd-modtools`. Always verify your working directory before running commands.

### Files to Read Before Starting

Before implementing ANY task, read these files to understand existing patterns:

**In 7dtd-gamedata:**

- `gamexml/blocks.go` - Existing Block struct
- `gamexml/entityclasses.go` - Existing EntityClass struct
- `gamexml/gamexml.go` - Common types (Property, Tags)
- `gamexml/tests/blocks_test.go` - Test patterns

**In 7dtd-modtools:**

- `modlet/functions.go` - Existing template functions and FuncArgs struct
- `modlet/unpack.go` - Template execution flow
- `modlet/tests/modlet_test.go` - Test setup pattern (afero MemMapFs, logger.Testing)
- `modlet/tests/functions_test.go` - Function test patterns

### Testing Conventions

Both projects use these patterns:

```go
// In 7dtd-modtools tests:
func setup(t *testing.T) *assert.Assertions {
    t.Helper()
    logger.Testing = true  // CRITICAL: prevents os.Exit()
    modlet.FS = FS         // Use in-memory filesystem
    mkTempDir(t)
    return assert.New(t)
}
```

**Always set `logger.Testing = true`** in tests to prevent panics from becoming `os.Exit()` calls.

### Test Helper Dependencies (7dtd-modtools)

The integration tests in Tasks 12-13 use helper functions defined in `modlet/tests/modlet_test.go`. **Do not redefine these** - they are already available in the test package.

| Helper         | Defined In                    | Purpose                                                                      |
| -------------- | ----------------------------- | ---------------------------------------------------------------------------- |
| `setup(t)`     | `modlet/tests/modlet_test.go` | Sets `logger.Testing=true`, `modlet.FS=FS`, creates temp dir, returns assert |
| `cleanup(t)`   | `modlet/tests/modlet_test.go` | Removes temp directory after test                                            |
| `FS`           | `modlet/tests/modlet_test.go` | Package-level `*afero.Afero` pointing to `MemMapFs`                          |
| `testTMP`      | `modlet/tests/modlet_test.go` | Package-level string with temp directory path                                |
| `mkTempDir(t)` | `modlet/tests/modlet_test.go` | Creates the temp directory in MemMapFs                                       |

**When writing new test files in `modlet/tests/`:**

1. The file must be in package `modlet_test` (note the `_test` suffix)
2. Import the helpers implicitly - they're in the same test package
3. Always call `setup(t)` at the start and `defer cleanup(t)`
4. Use `FS` for all filesystem operations (not `os` or `afero.NewOsFs()`)
5. Use `testTMP` as the base path for test files

**Example test structure:**

```go
package modlet_test  // Same package as modlet_test.go

import (
    "path/filepath"
    "testing"

    "github.com/donovanmods/7dtd-modtools/modlet"
    "github.com/stretchr/testify/require"
)

func TestMyFeature(t *testing.T) {
    assert := setup(t)      // Uses existing helper
    defer cleanup(t)        // Uses existing helper

    // Use FS and testTMP from modlet_test.go
    path := filepath.Join(testTMP, "myfile.txt")
    require.NoError(t, FS.WriteFile(path, []byte("content"), 0644))

    // ... test logic ...
}
```

---

## Task 1: Add DropEvents Method to Block Type

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata`

**Read First:**

- `gamexml/blocks.go` - Understand Block struct and existing Events field

**Files:**

- Modify: `gamexml/blocks.go`
- Create: `gamexml/tests/blocks_methods_test.go`

**Step 1: Write the failing test**

```go
// blocks_methods_test.go
package gamexml_test

import (
	"testing"

	"github.com/donovanmods/7dtd-gamedata/gamexml"
	"github.com/stretchr/testify/assert"
)

func TestBlockDropEvents(t *testing.T) {
	block := gamexml.Block{
		Name: "terrStone",
		Events: []gamexml.DropEvent{
			{Event: "Harvest", Name: "resourceRockSmall", Count: "55"},
			{Event: "Harvest", Name: "resourceite", Count: "10"},
			{Event: "Destroy", Name: "resourceRockSmall", Count: "5"},
		},
	}

	harvest := block.DropEvents("Harvest")
	assert.Len(t, harvest, 2)
	assert.Equal(t, "resourceRockSmall", harvest[0].Name)
	assert.Equal(t, "resourceite", harvest[1].Name)

	destroy := block.DropEvents("Destroy")
	assert.Len(t, destroy, 1)
	assert.Equal(t, "resourceRockSmall", destroy[0].Name)

	fall := block.DropEvents("Fall")
	assert.Len(t, fall, 0)
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go test -v -run TestBlockDropEvents ./gamexml/tests/
```

Expected: FAIL with "block.DropEvents undefined"

**Step 3: Write minimal implementation**

Add to `gamexml/blocks.go`:

```go
// DropEvents returns all drop events matching the given event type (e.g., "Harvest", "Destroy")
func (b Block) DropEvents(event string) []DropEvent {
	var events []DropEvent
	for _, e := range b.Events {
		if e.Event == event {
			events = append(events, e)
		}
	}
	return events
}
```

**Step 4: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go test -v -run TestBlockDropEvents ./gamexml/tests/
```

Expected: PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
git add gamexml/blocks.go gamexml/tests/blocks_methods_test.go
git commit -m "feat(gamexml): add DropEvents method to Block type"
```

---

## Task 2: Add GetProperty and HasProperty Methods to Block Type

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata`

**Read First:**

- `gamexml/blocks.go` - Note that `Property` is a field name, so method must be named differently

**Files:**

- Modify: `gamexml/blocks.go`
- Modify: `gamexml/tests/blocks_methods_test.go`

**Step 1: Write the failing test**

Add to `gamexml/tests/blocks_methods_test.go`:

```go
func TestBlockGetProperty(t *testing.T) {
	block := gamexml.Block{
		Name: "terrStone",
		Property: []gamexml.Property{
			{Name: "Material", Value: "Mite"},
			{Name: "TerrainIndex", Value: "1"},
			{Name: "Shape", Value: "Terrain"},
		},
	}

	assert.Equal(t, "Mite", block.GetProperty("Material"))
	assert.Equal(t, "1", block.GetProperty("TerrainIndex"))
	assert.Equal(t, "", block.GetProperty("NonExistent"))
}

func TestBlockHasProperty(t *testing.T) {
	block := gamexml.Block{
		Name: "terrStone",
		Property: []gamexml.Property{
			{Name: "Material", Value: "Mite"},
		},
	}

	assert.True(t, block.HasProperty("Material"))
	assert.False(t, block.HasProperty("NonExistent"))
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go test -v -run "TestBlock(GetProperty|HasProperty)" ./gamexml/tests/
```

Expected: FAIL with "block.GetProperty undefined"

**Step 3: Write minimal implementation**

Add to `gamexml/blocks.go`:

```go
// GetProperty returns the value of a property by name, or empty string if not found
func (b Block) GetProperty(name string) string {
	for _, p := range b.Property {
		if p.Name == name {
			return p.Value
		}
	}
	return ""
}

// HasProperty returns true if the block has a property with the given name
func (b Block) HasProperty(name string) bool {
	for _, p := range b.Property {
		if p.Name == name {
			return true
		}
	}
	return false
}
```

**Step 4: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go test -v -run "TestBlock(GetProperty|HasProperty)" ./gamexml/tests/
```

Expected: PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
git add gamexml/blocks.go gamexml/tests/blocks_methods_test.go
git commit -m "feat(gamexml): add GetProperty and HasProperty methods to Block type"
```

---

## Task 3: Add DropEvents Field and Method to EntityClass Type

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata`

**Read First:**

- `gamexml/entityclasses.go` - Note EntityClass currently has NO DropEvents field
- `gamexml/blocks.go` - Reference the DropEvent type definition

**Files:**

- Modify: `gamexml/entityclasses.go`
- Create: `gamexml/tests/entityclasses_methods_test.go`

**Step 1: Write the failing test**

```go
// entityclasses_methods_test.go
package gamexml_test

import (
	"testing"

	"github.com/donovanmods/7dtd-gamedata/gamexml"
	"github.com/stretchr/testify/assert"
)

func TestEntityClassGetDropEvents(t *testing.T) {
	entity := gamexml.EntityClass{
		Name: "animalChicken",
		DropEvents: []gamexml.DropEvent{
			{Event: "Harvest", Name: "foodRawMeat", Count: "5", Tag: "butcherHarvest"},
			{Event: "Harvest", Name: "resourceFeather", Count: "10", Tag: "allToolsHarvest"},
			{Event: "Destroy", Name: "resourceBone", Count: "1"},
		},
	}

	harvest := entity.GetDropEvents("Harvest")
	assert.Len(t, harvest, 2)
	assert.Equal(t, "foodRawMeat", harvest[0].Name)

	destroy := entity.GetDropEvents("Destroy")
	assert.Len(t, destroy, 1)
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go test -v -run TestEntityClassGetDropEvents ./gamexml/tests/
```

Expected: FAIL with "unknown field 'DropEvents' in struct literal"

**Step 3: Write minimal implementation**

Modify `gamexml/entityclasses.go` - add DropEvents field and method:

```go
type EntityClass struct {
	XMLName     string        `xml:"entity_class"`
	Name        string        `xml:"name,attr"`
	Property    []Property    `xml:"property,omitempty"`
	EffectGroup []EffectGroup `xml:"effect_group,omitempty"`
	DropEvents  []DropEvent   `xml:"drop,omitempty"` // ADD THIS LINE
}

// GetDropEvents returns all drop events matching the given event type
func (e EntityClass) GetDropEvents(event string) []DropEvent {
	var events []DropEvent
	for _, ev := range e.DropEvents {
		if ev.Event == event {
			events = append(events, ev)
		}
	}
	return events
}
```

**Step 4: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go test -v -run TestEntityClassGetDropEvents ./gamexml/tests/
```

Expected: PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
git add gamexml/entityclasses.go gamexml/tests/entityclasses_methods_test.go
git commit -m "feat(gamexml): add DropEvents field and GetDropEvents method to EntityClass"
```

---

## CHECKPOINT A: Verify 7dtd-gamedata

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata`

**Run all tests and verify everything passes before continuing:**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go test -v ./...
```

**Expected:** All tests PASS (existing tests + 3 new test functions)

**If any test fails:** STOP and diagnose before proceeding to Task 4.

### Recovery Guide for Checkpoint A

**Common failures and fixes:**

| Symptom                                          | Likely Cause                                    | Fix                                                                                       |
| ------------------------------------------------ | ----------------------------------------------- | ----------------------------------------------------------------------------------------- |
| `undefined: DropEvent`                           | DropEvent type not exported or in wrong package | Check `gamexml/blocks.go` - ensure `DropEvent` struct is defined and exported (capital D) |
| `unknown field 'DropEvents'` in EntityClass test | Field not added to struct                       | Add `DropEvents []DropEvent` field to EntityClass in `gamexml/entityclasses.go`           |
| Test compiles but wrong results                  | Method logic error                              | Compare method implementation against test expectations; add debug prints                 |
| Import cycle error                               | Circular import between packages                | Move shared types to a common package or restructure                                      |

**Diagnostic commands:**

```bash
# Check what's exported from gamexml
go doc github.com/donovanmods/7dtd-gamedata/gamexml

# Run single failing test with verbose output
go test -v -run TestBlockDropEvents ./gamexml/tests/ 2>&1

# Check for compilation errors only
go build ./...
```

**If stuck:** Review the original Block struct in `gamexml/blocks.go` to ensure your additions match the existing patterns (XML tags, field naming conventions).

---

## Task 4: Create GameData Loader Package

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Read First:**

- `go.mod` - Note the `replace` directive for 7dtd-gamedata
- `modlet/modlet.go` - See how afero.Afero is used (FS variable pattern)

**Files:**

- Create: `gamedata/gamedata.go`
- Create: `gamedata/gamedata_test.go`

**Step 1: Write the failing test**

```go
// gamedata/gamedata_test.go
package gamedata_test

import (
	"path/filepath"
	"testing"

	"github.com/donovanmods/7dtd-modtools/gamedata"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadGameData(t *testing.T) {
	// Use in-memory filesystem
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}
	gamedata.FS = fs

	tmpDir := "/tmp/gamedata"
	require.NoError(t, fs.MkdirAll(tmpDir, 0755))

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

	entityClassesXML := `<?xml version="1.0" encoding="UTF-8"?>
<entity_classes>
	<entity_class name="animalChicken">
		<property name="Class" value="EntityAnimal"/>
		<drop event="Harvest" name="foodRawMeat" count="5" tag="butcherHarvest"/>
	</entity_class>
</entity_classes>`

	require.NoError(t, fs.WriteFile(filepath.Join(tmpDir, "blocks.xml"), []byte(blocksXML), 0644))
	require.NoError(t, fs.WriteFile(filepath.Join(tmpDir, "entityclasses.xml"), []byte(entityClassesXML), 0644))

	gd, err := gamedata.Load(tmpDir)
	require.NoError(t, err)

	assert.Len(t, gd.Blocks, 2)
	assert.Equal(t, "terrStone", gd.Blocks[0].Name)
	assert.Equal(t, "Mite", gd.Blocks[0].GetProperty("Material"))

	assert.Len(t, gd.EntityClasses, 1)
	assert.Equal(t, "animalChicken", gd.EntityClasses[0].Name)
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run TestLoadGameData ./gamedata/
```

Expected: FAIL with "package gamedata is not in GOROOT" or similar

**Step 3: Write minimal implementation**

```go
// gamedata/gamedata.go
package gamedata

import (
	"encoding/xml"
	"fmt"
	"path/filepath"

	"github.com/donovanmods/7dtd-gamedata/gamexml"
	"github.com/spf13/afero"
)

// FS is the filesystem to use (can be swapped for testing)
var FS = &afero.Afero{Fs: afero.NewOsFs()}

// GameData holds parsed game XML data
type GameData struct {
	Blocks        []gamexml.Block
	EntityClasses []gamexml.EntityClass
}

// Load parses game XML files from the given directory
func Load(gamedir string) (*GameData, error) {
	gd := &GameData{}

	// Load blocks.xml
	blocksPath := filepath.Join(gamedir, "blocks.xml")
	if err := gd.loadBlocks(blocksPath); err != nil {
		return nil, fmt.Errorf("loading blocks.xml: %w", err)
	}

	// Load entityclasses.xml
	entityPath := filepath.Join(gamedir, "entityclasses.xml")
	if err := gd.loadEntityClasses(entityPath); err != nil {
		return nil, fmt.Errorf("loading entityclasses.xml: %w", err)
	}

	return gd, nil
}

func (gd *GameData) loadBlocks(path string) error {
	data, err := FS.ReadFile(path)
	if err != nil {
		return err
	}

	var blocks gamexml.Blocks
	if err := xml.Unmarshal(data, &blocks); err != nil {
		return err
	}

	gd.Blocks = blocks.Block
	return nil
}

func (gd *GameData) loadEntityClasses(path string) error {
	data, err := FS.ReadFile(path)
	if err != nil {
		return err
	}

	var entities gamexml.EntityClasses
	if err := xml.Unmarshal(data, &entities); err != nil {
		return err
	}

	gd.EntityClasses = entities.EntityClass
	return nil
}
```

**Step 4: Run go mod tidy**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go mod tidy
```

**Step 5: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run TestLoadGameData ./gamedata/
```

Expected: PASS

**Step 6: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add gamedata/
git commit -m "feat(gamedata): add GameData loader for blocks and entityclasses"
```

---

## Task 5: Add String Helper Functions

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Read First:**

- `modlet/functions.go` - See existing function patterns
- `lib/logger/logger.go` - Understand logger.Error usage

**Files:**

- Create: `modlet/helpers.go`
- Create: `modlet/tests/helpers_test.go`

**Step 1: Write the failing test**

```go
// modlet/tests/helpers_test.go
package modlet_test

import (
	"testing"

	"github.com/donovanmods/7dtd-modtools/modlet"
	"github.com/stretchr/testify/assert"
)

func TestHasPrefix(t *testing.T) {
	assert.True(t, modlet.HasPrefix("terrStone", "terr"))
	assert.True(t, modlet.HasPrefix("terrDirt", "terr"))
	assert.False(t, modlet.HasPrefix("rockStone", "terr"))
	assert.False(t, modlet.HasPrefix("", "terr"))
}

func TestHasSuffix(t *testing.T) {
	assert.True(t, modlet.HasSuffix("terrStone", "Stone"))
	assert.True(t, modlet.HasSuffix("rockStone", "Stone"))
	assert.False(t, modlet.HasSuffix("terrDirt", "Stone"))
}

func TestMatch(t *testing.T) {
	assert.True(t, modlet.Match("terrStone", `^terr`))
	assert.True(t, modlet.Match("plantedCorn", `^planted`))
	assert.True(t, modlet.Match("blockShapes", `Shapes$`))
	assert.False(t, modlet.Match("terrStone", `^planted`))
}

func TestNotMatch(t *testing.T) {
	assert.True(t, modlet.NotMatch("terrStone", `^planted`))
	assert.False(t, modlet.NotMatch("plantedCorn", `^planted`))
}

func TestMultValue(t *testing.T) {
	assert.Equal(t, "82", modlet.MultValue("55", 1.5))
	assert.Equal(t, "8", modlet.MultValue("5", 1.5))
	assert.Equal(t, "15,30,45", modlet.MultValue("10,20,30", 1.5))
	assert.Equal(t, "0", modlet.MultValue("", 1.5))
}

func TestProbMult(t *testing.T) {
	assert.Equal(t, "0.45", modlet.ProbMult("0.3", 1.5))
	assert.Equal(t, "1", modlet.ProbMult("0.8", 1.5))
	assert.Equal(t, "1", modlet.ProbMult("0.9", 2.0))
	assert.Equal(t, "1", modlet.ProbMult("1.0", 1.5))
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run "Test(HasPrefix|HasSuffix|Match|NotMatch|MultValue|ProbMult)" ./modlet/tests/
```

Expected: FAIL with "undefined: modlet.HasPrefix"

**Step 3: Write minimal implementation**

```go
// modlet/helpers.go
package modlet

import (
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/donovanmods/7dtd-modtools/lib/logger"
)

// HasPrefix returns true if s starts with prefix
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// HasSuffix returns true if s ends with suffix
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// Match returns true if s matches the regex pattern
func Match(s, pattern string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		logger.Error("invalid regex pattern %q: %v", pattern, err)
		return false
	}
	return re.MatchString(s)
}

// NotMatch returns true if s does NOT match the regex pattern
func NotMatch(s, pattern string) bool {
	return !Match(s, pattern)
}

// MultValue multiplies a numeric value (or CSV of values) by factor, rounding to int
func MultValue(value string, factor float64) string {
	if value == "" {
		return "0"
	}

	parts := strings.Split(value, ",")
	results := make([]string, len(parts))

	for i, part := range parts {
		part = strings.TrimSpace(part)
		num, err := strconv.ParseFloat(part, 64)
		if err != nil {
			logger.Error("invalid number %q: %v", part, err)
			results[i] = part
			continue
		}
		result := math.Round(num * factor)
		results[i] = strconv.FormatFloat(result, 'f', 0, 64)
	}

	return strings.Join(results, ",")
}

// ProbMult multiplies a probability value by factor, capping at 1.0
func ProbMult(value string, factor float64) string {
	if value == "" {
		return "0"
	}

	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		logger.Error("invalid probability %q: %v", value, err)
		return value
	}

	result := num * factor
	if result > 1.0 {
		result = 1.0
	}

	formatted := strconv.FormatFloat(result, 'f', -1, 64)
	return formatted
}
```

**Step 4: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run "Test(HasPrefix|HasSuffix|Match|NotMatch|MultValue|ProbMult)" ./modlet/tests/
```

Expected: PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add modlet/helpers.go modlet/tests/helpers_test.go
git commit -m "feat(modlet): add string helper functions for templates"
```

---

## Task 6: Complete mult Function with min/max Bounds

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Read First:**

- `modlet/functions.go:124-155` - Current MultFunc implementation with TODO comment
- `modlet/tests/functions_test.go` - Existing test patterns

**Files:**

- Modify: `modlet/functions.go`
- Modify: `modlet/tests/functions_test.go`

**Step 1: Write the failing test**

Add to `modlet/tests/functions_test.go`:

```go
func TestMultFuncWithBounds(t *testing.T) {
	assert := setup(t)
	defer cleanup(t)

	funcArgs := modlet.FuncArgs{
		ModInfo: &modinfo.ModInfo{},
	}

	fn := modlet.MultFunc(funcArgs)

	// Basic multiplication (no bounds)
	result := fn("//test/@value", "by=2.0")
	assert.Contains(result, "<set")
	assert.Contains(result, "2")

	// With min bound - result should be clamped to min
	result = fn("//test/@value", "by=0.5", "min=1")
	assert.Contains(result, ">1<")

	// With max bound - result should be clamped to max
	result = fn("//test/@value", "by=10.0", "max=5")
	assert.Contains(result, ">5<")

	// Both bounds
	result = fn("//test/@value", "by=0.1", "min=2", "max=10")
	assert.Contains(result, ">2<")
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run TestMultFuncWithBounds ./modlet/tests/
```

Expected: FAIL (min/max not applied, test assertions fail)

**Step 3: Replace MultFunc implementation**

In `modlet/functions.go`, replace the entire `MultFunc` function (lines ~124-155):

```go
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
```

**Step 4: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run TestMultFuncWithBounds ./modlet/tests/
```

Expected: PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add modlet/functions.go modlet/tests/functions_test.go
git commit -m "feat(modlet): complete mult function with min/max bounds"
```

---

## Task 7: Add comment Template Function

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Read First:**

- `modlet/functions.go` - See existing function patterns, note `fmt` may need to be added to imports
- `modlet/unpack.go:167-176` - Where FuncMap is defined

**Files:**

- Modify: `modlet/functions.go`
- Modify: `modlet/unpack.go`
- Modify: `modlet/tests/functions_test.go`

**Step 1: Write the failing test**

Add to `modlet/tests/functions_test.go`:

```go
func TestCommentFunc(t *testing.T) {
	_ = setup(t)

	fn := modlet.CommentFunc()

	result := fn("This is a comment")
	assert.New(t).Equal("<!-- This is a comment -->", result)

	result = fn("Multiple\nlines")
	assert.New(t).Equal("<!-- Multiple\nlines -->", result)
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run TestCommentFunc ./modlet/tests/
```

Expected: FAIL with "undefined: modlet.CommentFunc"

**Step 3: Add CommentFunc to functions.go**

Add to `modlet/functions.go` (ensure `"fmt"` is in imports):

```go
// CommentFunc returns an XML comment string
func CommentFunc() func(string) string {
	return func(text string) string {
		return fmt.Sprintf("<!-- %s -->", text)
	}
}
```

**Step 4: Register in unpack.go FuncMap**

In `modlet/unpack.go`, update the `Funcs(template.FuncMap{...})` block to add:

```go
"comment": CommentFunc(),
```

**Step 5: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run TestCommentFunc ./modlet/tests/
```

Expected: PASS

**Step 6: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add modlet/functions.go modlet/unpack.go modlet/tests/functions_test.go
git commit -m "feat(modlet): add comment template function"
```

---

## Task 8: Add prob Template Function

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Read First:**

- `modlet/helpers.go` - ProbMult function we'll use
- `modlet/functions.go` - Pattern from SetFunc and MultFunc

**Files:**

- Modify: `modlet/functions.go`
- Modify: `modlet/unpack.go`
- Modify: `modlet/tests/functions_test.go`

**Step 1: Write the failing test**

Add to `modlet/tests/functions_test.go`:

```go
func TestProbFunc(t *testing.T) {
	assert := setup(t)

	funcArgs := modlet.FuncArgs{
		ModInfo: &modinfo.ModInfo{},
	}

	fn := modlet.ProbFunc(funcArgs)

	// Normal multiplication
	result := fn("//block/@prob", "0.3", "by=1.5")
	assert.Contains(result, "<set")
	assert.Contains(result, "0.45")

	// Capped at 1.0
	result = fn("//block/@prob", "0.8", "by=1.5")
	assert.Contains(result, ">1<")
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run TestProbFunc ./modlet/tests/
```

Expected: FAIL with "undefined: modlet.ProbFunc"

**Step 3: Add ProbFunc to functions.go**

```go
// ProbFunc creates a probability instruction with value capped at 1.0
func ProbFunc(fargs FuncArgs) func(string, string, ...string) string {
	return func(xpath string, value string, args ...string) string {
		xpath = strings.TrimSpace(xpath)
		if xpath == "" {
			logger.Fatal("xpath must be provided to the prob command")
		}

		pargs := ParseArgs(args)

		var multiplier float64 = 1.0
		var err error

		if by, ok := pargs["by"]; ok && by != "" {
			if multiplier, err = strconv.ParseFloat(by, 64); err != nil {
				logger.Fatal("error parsing multiplier %q: %w", by, err)
			}
		}

		result := ProbMult(value, multiplier)
		return must(mkSet(xpath, result))
	}
}
```

**Step 4: Register in unpack.go FuncMap**

Add to the FuncMap in `modlet/unpack.go`:

```go
"prob": ProbFunc(fargs),
```

**Step 5: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run TestProbFunc ./modlet/tests/
```

Expected: PASS

**Step 6: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add modlet/functions.go modlet/unpack.go modlet/tests/functions_test.go
git commit -m "feat(modlet): add prob template function with 1.0 cap"
```

---

## CHECKPOINT B: Verify modtools Tests Pass

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Run all existing tests:**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v ./...
```

**Expected:** All tests PASS

**If any test fails:** STOP and diagnose before proceeding to Task 9.

### Recovery Guide for Checkpoint B

**Common failures and fixes:**

| Symptom                          | Likely Cause                             | Fix                                                                              |
| -------------------------------- | ---------------------------------------- | -------------------------------------------------------------------------------- |
| `cannot find package "gamedata"` | Package not created or wrong import path | Verify `gamedata/gamedata.go` exists with correct `package gamedata` declaration |
| `undefined: modlet.HasPrefix`    | Function not exported                    | Ensure function names start with capital letter in `modlet/helpers.go`           |
| `logger.Fatal` causes test panic | `logger.Testing` not set                 | Ensure test calls `setup(t)` which sets `logger.Testing = true`                  |
| `MultFunc` test fails on bounds  | Bounds logic incorrect                   | Review the min/max clamping logic in `MultFunc`                                  |
| `FuncMap` registration error     | Function signature mismatch              | Check that function signatures match what template expects                       |

**Diagnostic commands:**

```bash
# Check what's exported from modlet package
go doc github.com/donovanmods/7dtd-modtools/modlet | head -50

# Run single failing test
go test -v -run TestMultFuncWithBounds ./modlet/tests/

# Check imports are resolved
go mod tidy && go build ./...

# See all test functions
go test -list '.*' ./modlet/tests/
```

**If stuck:**

1. Check that `go.mod` has the `replace` directive for 7dtd-gamedata pointing to the local path
2. Run `go mod tidy` to fix import issues
3. Compare your function signatures against existing functions in `modlet/functions.go`

---

## Task 9: Add GameData Field to FuncArgs

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Read First:**

- `modlet/functions.go:58-66` - Current FuncArgs struct definition

**Files:**

- Modify: `modlet/functions.go`

**Step 1: Add import for gamedata package**

At the top of `modlet/functions.go`, add to imports:

```go
"github.com/donovanmods/7dtd-modtools/gamedata"
```

**Step 2: Add GameData field to FuncArgs**

Update the FuncArgs struct in `modlet/functions.go`:

```go
type FuncArgs struct {
	Output   string
	Gamedir  string
	ModInfo  *modinfo.ModInfo
	FBuffer  *FileBuffer
	GBuffer  *bytes.Buffer
	FBufMap  FileBufferMap
	Options  map[string]string
	GameData *gamedata.GameData // ADD THIS LINE
}
```

**Step 3: Verify compilation**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go build ./...
```

Expected: Build succeeds (no errors)

**Step 4: Run existing tests to ensure no regression**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v ./modlet/tests/
```

Expected: All existing tests still PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add modlet/functions.go
git commit -m "refactor(modlet): add GameData field to FuncArgs struct"
```

---

## Task 10: Load GameData in Unpack Function

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Read First:**

- `modlet/unpack.go` - Full Unpack function, note where fargs is created
- `gamedata/gamedata.go` - Load function signature

**Files:**

- Modify: `modlet/unpack.go`

**Step 1: Add import for gamedata package**

At the top of `modlet/unpack.go`, add to imports:

```go
"github.com/donovanmods/7dtd-modtools/gamedata"
```

**Step 2: Add game data loading after path cleaning**

In the `Unpack` function, after the line `output = filepath.Clean(output)`, add:

```go
// Load game data
var gd *gamedata.GameData
if gamedir != "" && gamedir != "." {
	var err error
	gd, err = gamedata.Load(gamedir)
	if err != nil {
		logger.Warn("could not load game data from %s: %v", gamedir, err)
		gd = &gamedata.GameData{} // Empty game data
	}
} else {
	gd = &gamedata.GameData{}
}
```

**Step 3: Add GameData to fargs**

Update the fargs initialization to include GameData:

```go
fargs := FuncArgs{
	Output:   output,
	Gamedir:  gamedir,
	ModInfo:  &modInfo,
	FBuffer:  fBuffer,
	GBuffer:  gBuffer,
	FBufMap:  fBufMap,
	GameData: gd, // ADD THIS LINE
	Options: map[string]string{
		"force": strconv.FormatBool(force),
	},
}
```

**Step 4: Pass GameData to template execution**

Replace the `t.ExecuteTemplate` call with:

```go
// Create template data with GameData
templateData := struct {
	GameData *gamedata.GameData
}{
	GameData: gd,
}

if err := t.ExecuteTemplate(gBuffer, templateName, templateData); err != nil {
	return fmt.Errorf("error executing template %s: %w", templateName, err)
}
```

**Step 5: Verify compilation**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go build ./...
```

Expected: Build succeeds

**Step 6: Run existing tests**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v ./modlet/tests/
```

Expected: All existing tests still PASS (they don't use GameData yet)

**Step 7: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add modlet/unpack.go
git commit -m "feat(modlet): load GameData in Unpack and pass to template"
```

---

## Task 11: Update Test Setup for GameData Filesystem

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Read First:**

- `modlet/tests/modlet_test.go` - Current setup() function
- `gamedata/gamedata.go` - Note the FS variable

**Files:**

- Modify: `modlet/tests/modlet_test.go`

**Step 1: Add import for gamedata package**

Add to imports in `modlet/tests/modlet_test.go`:

```go
"github.com/donovanmods/7dtd-modtools/gamedata"
```

**Step 2: Update setup() to set gamedata.FS**

Update the `setup` function:

```go
func setup(t *testing.T) *assert.Assertions {
	t.Helper()

	logger.Testing = true

	// Use MemMapFs for testing
	modlet.FS = FS
	gamedata.FS = FS // ADD THIS LINE

	mkTempDir(t)

	return assert.New(t)
}
```

**Step 3: Run tests to verify setup works**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v ./modlet/tests/
```

Expected: All tests PASS

**Step 4: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add modlet/tests/modlet_test.go
git commit -m "test(modlet): update setup to share filesystem with gamedata"
```

---

## Task 12: Write GameData Integration Test

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Read First:**

- `modlet/tests/unpack_test.go` - If exists, see existing patterns
- `modlet/tests/modlet_test.go` - setup/cleanup patterns

**Files:**

- Create: `modlet/tests/gamedata_integration_test.go`

**Step 1: Write the integration test**

```go
// modlet/tests/gamedata_integration_test.go
package modlet_test

import (
	"path/filepath"
	"testing"

	"github.com/donovanmods/7dtd-modtools/modlet"
	"github.com/stretchr/testify/require"
)

func TestUnpackWithGameData(t *testing.T) {
	assert := setup(t)
	defer cleanup(t)

	// Create game data directory with test XML
	gameDir := filepath.Join(testTMP, "gamedata")
	require.NoError(t, FS.MkdirAll(gameDir, 0755))

	blocksXML := `<?xml version="1.0" encoding="UTF-8"?>
<blocks>
	<block name="terrStone">
		<property name="Material" value="Mite"/>
		<drop event="Harvest" name="resourceRockSmall" count="55"/>
	</block>
</blocks>`

	entityXML := `<?xml version="1.0" encoding="UTF-8"?>
<entity_classes>
	<entity_class name="animalChicken">
		<drop event="Harvest" name="foodRawMeat" count="5" tag="butcherHarvest"/>
	</entity_class>
</entity_classes>`

	require.NoError(t, FS.WriteFile(filepath.Join(gameDir, "blocks.xml"), []byte(blocksXML), 0644))
	require.NoError(t, FS.WriteFile(filepath.Join(gameDir, "entityclasses.xml"), []byte(entityXML), 0644))

	// Create template that uses game data
	tmplContent := `{{- modlet "test-gamedata" -}}
{{- output "Config/blocks.xml" -}}
{{ xmlHeader }}
<configs>
{{- range .GameData.Blocks }}
{{ comment .Name }}
{{- end }}
</configs>
{{- write -}}`

	tmplPath := filepath.Join(testTMP, "test.tmpl")
	require.NoError(t, FS.WriteFile(tmplPath, []byte(tmplContent), 0644))

	// Unpack with game data
	args := modlet.CmdArgs{
		Input:   []string{tmplPath},
		Output:  testTMP,
		Gamedir: gameDir,
	}

	err := modlet.Unpack(tmplPath, args)
	require.NoError(t, err)

	// Verify output contains block name in comment
	outputPath := filepath.Join(testTMP, "test-gamedata", "Config", "blocks.xml")
	content, err := FS.ReadFile(outputPath)
	require.NoError(t, err)

	assert.Contains(string(content), "<!-- terrStone -->")
}
```

**Step 2: Run test**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run TestUnpackWithGameData ./modlet/tests/
```

Expected: PASS

**Step 3: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add modlet/tests/gamedata_integration_test.go
git commit -m "test(modlet): add integration test for GameData in templates"
```

---

## Task 13: Add Helper Functions to Template FuncMap

> **Working Directory:** `/home/dyoung/Projects/mods/7dtd/7dtd-modtools`

**Read First:**

- `modlet/unpack.go` - Current FuncMap in NewTemplate
- `modlet/helpers.go` - Functions to add

**Files:**

- Modify: `modlet/unpack.go`
- Create: `modlet/tests/template_helpers_test.go`

**Step 1: Write the test**

```go
// modlet/tests/template_helpers_test.go
package modlet_test

import (
	"path/filepath"
	"testing"

	"github.com/donovanmods/7dtd-modtools/modlet"
	"github.com/stretchr/testify/require"
)

func TestTemplateHelperFunctions(t *testing.T) {
	assert := setup(t)
	defer cleanup(t)

	// Create minimal game data
	gameDir := filepath.Join(testTMP, "gamedata")
	require.NoError(t, FS.MkdirAll(gameDir, 0755))
	require.NoError(t, FS.WriteFile(filepath.Join(gameDir, "blocks.xml"), []byte(`<blocks></blocks>`), 0644))
	require.NoError(t, FS.WriteFile(filepath.Join(gameDir, "entityclasses.xml"), []byte(`<entity_classes></entity_classes>`), 0644))

	// Template using helper functions
	tmplContent := `{{- modlet "test-helpers" -}}
{{- output "test.txt" -}}
hasPrefix: {{ hasPrefix "terrStone" "terr" }}
hasSuffix: {{ hasSuffix "terrStone" "Stone" }}
match: {{ match "plantedCorn" "^planted" }}
notMatch: {{ notMatch "terrStone" "^planted" }}
multValue: {{ multValue "55" 1.5 }}
probMult: {{ probMult "0.8" 1.5 }}
{{- write -}}`

	tmplPath := filepath.Join(testTMP, "helpers.tmpl")
	require.NoError(t, FS.WriteFile(tmplPath, []byte(tmplContent), 0644))

	args := modlet.CmdArgs{
		Input:   []string{tmplPath},
		Output:  testTMP,
		Gamedir: gameDir,
	}

	err := modlet.Unpack(tmplPath, args)
	require.NoError(t, err)

	content, err := FS.ReadFile(filepath.Join(testTMP, "test-helpers", "test.txt"))
	require.NoError(t, err)

	output := string(content)
	assert.Contains(output, "hasPrefix: true")
	assert.Contains(output, "hasSuffix: true")
	assert.Contains(output, "match: true")
	assert.Contains(output, "notMatch: true")
	assert.Contains(output, "multValue: 82")
	assert.Contains(output, "probMult: 1")
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run TestTemplateHelperFunctions ./modlet/tests/
```

Expected: FAIL with "function 'hasPrefix' not defined"

**Step 3: Add helper functions to FuncMap**

Update `NewTemplate` in `modlet/unpack.go`:

```go
return template.New(name).
	Funcs(template.FuncMap{
		// Core functions
		"modlet":    ModletFunc(fargs),
		"mult":      MultFunc(fargs),
		"output":    OutputFunc(fargs),
		"set":       SetFunc(fargs),
		"write":     WriteFunc(fargs),
		"xmlHeader": func() string { return xml.Header },
		"comment":   CommentFunc(),
		"prob":      ProbFunc(fargs),
		// Helper functions
		"hasPrefix": HasPrefix,
		"hasSuffix": HasSuffix,
		"match":     Match,
		"notMatch":  NotMatch,
		"multValue": MultValue,
		"probMult":  ProbMult,
	}).
	Parse(string(data))
```

**Step 4: Run test to verify it passes**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run TestTemplateHelperFunctions ./modlet/tests/
```

Expected: PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add modlet/unpack.go modlet/tests/template_helpers_test.go
git commit -m "feat(modlet): add helper functions to template FuncMap"
```

---

## CHECKPOINT C: Final Verification

> **Verify both repositories before completing Phase 1**

**Step 1: Run all tests in 7dtd-gamedata**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
go test -v ./...
```

Expected: All tests PASS

**Step 2: Run all tests in 7dtd-modtools**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v ./...
```

Expected: All tests PASS

**Step 3: Run lint checks**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
task check
```

Expected: No errors (warnings OK)

**Step 4: Run go mod tidy in both repos**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata && go mod tidy
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go mod tidy
```

**Step 5: Final commit if any cleanup needed**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git status
# If there are changes:
git add -A
git commit -m "chore: Phase 1 cleanup"
```

### Recovery Guide for Checkpoint C

**Common failures and fixes:**

| Symptom                                      | Likely Cause                               | Fix                                                                                               |
| -------------------------------------------- | ------------------------------------------ | ------------------------------------------------------------------------------------------------- |
| Integration test fails on template execution | `templateData` struct not passed correctly | Check Task 10 implementation - ensure `ExecuteTemplate` receives the struct with `GameData` field |
| `function 'hasPrefix' not defined`           | Helper not added to FuncMap                | Check `modlet/unpack.go` FuncMap registration (Task 13)                                           |
| Lint errors about unused imports             | Import added but function not used         | Remove unused import or add the function that uses it                                             |
| `gamedata.FS` nil pointer                    | Test setup doesn't set `gamedata.FS`       | Ensure `setup(t)` includes `gamedata.FS = FS` (Task 11)                                           |
| Template output missing expected content     | Template data not accessible               | Verify template uses `.GameData.Blocks` (with dot prefix)                                         |

**Diagnostic commands:**

```bash
# Run integration test with verbose output
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
go test -v -run TestUnpackWithGameData ./modlet/tests/

# Check template FuncMap has all functions
grep -A 30 'template.FuncMap' modlet/unpack.go

# Verify gamedata package exports
go doc github.com/donovanmods/7dtd-modtools/gamedata

# Run lint with auto-fix for formatting issues
task format
```

**If integration tests fail:**

1. First verify the unit tests pass: `go test -v -run 'Test(HasPrefix|MultValue|CommentFunc)' ./modlet/tests/`
2. Check that `gamedata.FS` is set in test setup (Task 11)
3. Verify the template syntax - use `{{ .GameData.Blocks }}` not `{{ GameData.Blocks }}`
4. Add debug output to template: `{{ printf "%+v" .GameData }}`

**If lint fails:**

1. Run `task format` first to auto-fix formatting
2. For unused variable warnings, prefix with `_` or remove
3. For import order issues, use `goimports -w <file>`

---

## Summary

Phase 1 establishes the foundation for template-based modlet generation:

1. **7dtd-gamedata enhancements** (Tasks 1-3):
   - `Block.DropEvents(event)` - Filter drop events by type
   - `Block.GetProperty(name)` - Get property value
   - `Block.HasProperty(name)` - Check property exists
   - `EntityClass.DropEvents` field + `GetDropEvents(event)` method

2. **GameData loader** (Task 4):
   - `gamedata.Load(dir)` - Parse blocks.xml and entityclasses.xml
   - Uses afero for filesystem abstraction

3. **Helper functions** (Task 5):
   - `hasPrefix`, `hasSuffix` - String matching
   - `match`, `notMatch` - Regex matching
   - `multValue` - Multiply CSV values
   - `probMult` - Multiply probability with 1.0 cap

4. **Template functions** (Tasks 6-8):
   - `mult` - Completed with min/max bounds
   - `comment` - XML comment generation
   - `prob` - Probability with cap

5. **Template integration** (Tasks 9-13):
   - GameData passed to templates as `.GameData`
   - All helper functions available in templates

**Templates can now:**

- Access game blocks via `{{ range .GameData.Blocks }}`
- Filter by name patterns via `{{ if hasPrefix .Name "terr" }}`
- Get drop events via `{{ range .DropEvents "Harvest" }}`
- Generate XPath instructions via `{{ set }}`, `{{ mult }}`, `{{ prob }}`
- Add comments via `{{ comment "text" }}`
