# Phase 1: Core Infrastructure Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Complete the template engine foundation with game data support, enabling templates to iterate over game XML and generate modlets.

**Architecture:** Extend 7dtd-gamedata library with parsing for game XML files (blocks, entityclasses, items). Create a GameData struct in modtools that loads all game data and passes it to template execution. Add helper functions for string matching and value manipulation.

**Tech Stack:** Go 1.25, text/template, encoding/xml, 7dtd-gamedata library

---

## Task 1: Add DropEvents Method to Block Type (7dtd-gamedata)

**Files:**

- Modify: `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata/gamexml/blocks.go`
- Create: `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata/gamexml/tests/blocks_methods_test.go`

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

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata && go test -v -run TestBlockDropEvents ./gamexml/tests/`
Expected: FAIL with "block.DropEvents undefined"

**Step 3: Write minimal implementation**

Add to `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata/gamexml/blocks.go`:

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

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata && go test -v -run TestBlockDropEvents ./gamexml/tests/`
Expected: PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
git add gamexml/blocks.go gamexml/tests/blocks_methods_test.go
git commit -m "feat(gamexml): add DropEvents method to Block type"
```

---

## Task 2: Add Property and HasProperty Methods to Block Type (7dtd-gamedata)

**Files:**

- Modify: `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata/gamexml/blocks.go`
- Modify: `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata/gamexml/tests/blocks_methods_test.go`

**Step 1: Write the failing test**

Add to `blocks_methods_test.go`:

```go
func TestBlockProperty(t *testing.T) {
	block := gamexml.Block{
		Name: "terrStone",
		Property: []gamexml.Property{
			{Name: "Material", Value: "Mite"},
			{Name: "TerrainIndex", Value: "1"},
			{Name: "Shape", Value: "Terrain"},
		},
	}

	assert.Equal(t, "Mite", block.Property("Material"))
	assert.Equal(t, "1", block.Property("TerrainIndex"))
	assert.Equal(t, "", block.Property("NonExistent"))
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

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata && go test -v -run TestBlockProperty ./gamexml/tests/`
Expected: FAIL with "block.Property undefined" (method conflicts with field)

**Step 3: Write minimal implementation**

Add to `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata/gamexml/blocks.go`:

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

**Step 4: Update test to use GetProperty instead of Property**

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
```

**Step 5: Run test to verify it passes**

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata && go test -v -run "TestBlock(GetProperty|HasProperty)" ./gamexml/tests/`
Expected: PASS

**Step 6: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
git add gamexml/blocks.go gamexml/tests/blocks_methods_test.go
git commit -m "feat(gamexml): add GetProperty and HasProperty methods to Block type"
```

---

## Task 3: Add DropEvents Method to EntityClass Type (7dtd-gamedata)

**Files:**

- Modify: `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata/gamexml/entityclasses.go`
- Create: `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata/gamexml/tests/entityclasses_methods_test.go`

**Step 1: Write the failing test**

```go
// entityclasses_methods_test.go
package gamexml_test

import (
	"testing"

	"github.com/donovanmods/7dtd-gamedata/gamexml"
	"github.com/stretchr/testify/assert"
)

func TestEntityClassDropEvents(t *testing.T) {
	// Note: EntityClass doesn't have Events field yet - we need to add it
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

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata && go test -v -run TestEntityClassDropEvents ./gamexml/tests/`
Expected: FAIL with "unknown field 'DropEvents' in struct literal"

**Step 3: Write minimal implementation**

First, add DropEvents field to EntityClass in `/home/dyoung/Projects/mods/7dtd/7dtd-gamedata/gamexml/entityclasses.go`:

```go
type EntityClass struct {
	XMLName     string        `xml:"entity_class"`
	Name        string        `xml:"name,attr"`
	Property    []Property    `xml:"property,omitempty"`
	EffectGroup []EffectGroup `xml:"effect_group,omitempty"`
	DropEvents  []DropEvent   `xml:"drop,omitempty"` // Add this line
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

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata && go test -v -run TestEntityClassDropEvents ./gamexml/tests/`
Expected: PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata
git add gamexml/entityclasses.go gamexml/tests/entityclasses_methods_test.go
git commit -m "feat(gamexml): add DropEvents field and GetDropEvents method to EntityClass"
```

---

## Task 4: Create GameData Loader in modtools

**Files:**

- Create: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/gamedata/gamedata.go`
- Create: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/gamedata/gamedata_test.go`

**Step 1: Write the failing test**

```go
// gamedata_test.go
package gamedata_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/donovanmods/7dtd-modtools/gamedata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadGameData(t *testing.T) {
	// Create temp directory with test XML files
	tmpDir := t.TempDir()

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

	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "blocks.xml"), []byte(blocksXML), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "entityclasses.xml"), []byte(entityClassesXML), 0644))

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

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go test -v -run TestLoadGameData ./gamedata/`
Expected: FAIL with "package gamedata is not in GOROOT"

**Step 3: Write minimal implementation**

```go
// gamedata.go
package gamedata

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/donovanmods/7dtd-gamedata/gamexml"
)

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
	data, err := os.ReadFile(path)
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
	data, err := os.ReadFile(path)
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

**Step 4: Run test to verify it passes**

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go test -v -run TestLoadGameData ./gamedata/`
Expected: PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add gamedata/
git commit -m "feat(gamedata): add GameData loader for blocks and entityclasses"
```

---

## Task 5: Add String Helper Functions to Template System

**Files:**

- Create: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/helpers.go`
- Create: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/tests/helpers_test.go`

**Step 1: Write the failing test**

```go
// helpers_test.go
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
	// Simple integer
	assert.Equal(t, "82", modlet.MultValue("55", 1.5))

	// Float result rounds
	assert.Equal(t, "8", modlet.MultValue("5", 1.5))

	// CSV values
	assert.Equal(t, "15,30,45", modlet.MultValue("10,20,30", 1.5))

	// Empty returns "0"
	assert.Equal(t, "0", modlet.MultValue("", 1.5))
}

func TestProbMult(t *testing.T) {
	// Normal multiplication
	assert.Equal(t, "0.45", modlet.ProbMult("0.3", 1.5))

	// Cap at 1.0
	assert.Equal(t, "1", modlet.ProbMult("0.8", 1.5))
	assert.Equal(t, "1", modlet.ProbMult("0.9", 2.0))

	// Already 1.0
	assert.Equal(t, "1", modlet.ProbMult("1.0", 1.5))
}
```

**Step 2: Run test to verify it fails**

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go test -v -run "Test(HasPrefix|HasSuffix|Match|NotMatch|MultValue|ProbMult)" ./modlet/tests/`
Expected: FAIL with "undefined: modlet.HasPrefix"

**Step 3: Write minimal implementation**

```go
// helpers.go
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

	// Format without trailing zeros
	formatted := strconv.FormatFloat(result, 'f', -1, 64)
	return formatted
}
```

**Step 4: Run test to verify it passes**

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go test -v -run "Test(HasPrefix|HasSuffix|Match|NotMatch|MultValue|ProbMult)" ./modlet/tests/`
Expected: PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add modlet/helpers.go modlet/tests/helpers_test.go
git commit -m "feat(modlet): add string helper functions for templates"
```

---

## Task 6: Complete mult Function with min/max Bounds

**Files:**

- Modify: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/functions.go`
- Modify: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/tests/functions_test.go`

**Step 1: Write the failing test**

Add to `functions_test.go`:

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
	assert.Contains(result, ">1<") // Result clamped to min=1

	// With max bound - result should be clamped to max
	result = fn("//test/@value", "by=10.0", "max=5")
	assert.Contains(result, ">5<") // Result clamped to max=5

	// Both bounds
	result = fn("//test/@value", "by=0.1", "min=2", "max=10")
	assert.Contains(result, ">2<") // Result clamped to min=2
}
```

**Step 2: Run test to verify it fails**

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go test -v -run TestMultFuncWithBounds ./modlet/tests/`
Expected: FAIL (min/max not applied)

**Step 3: Update implementation**

Replace `MultFunc` in `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/functions.go`:

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

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go test -v -run TestMultFuncWithBounds ./modlet/tests/`
Expected: PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add modlet/functions.go modlet/tests/functions_test.go
git commit -m "feat(modlet): complete mult function with min/max bounds"
```

---

## Task 7: Add comment Template Function

**Files:**

- Modify: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/functions.go`
- Modify: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/unpack.go`
- Modify: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/tests/functions_test.go`

**Step 1: Write the failing test**

Add to `functions_test.go`:

```go
func TestCommentFunc(t *testing.T) {
	assert := setup(t)

	fn := modlet.CommentFunc()

	result := fn("This is a comment")
	assert.Equal("<!-- This is a comment -->", result)

	result = fn("Multiple\nlines")
	assert.Equal("<!-- Multiple\nlines -->", result)
}
```

**Step 2: Run test to verify it fails**

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go test -v -run TestCommentFunc ./modlet/tests/`
Expected: FAIL with "undefined: modlet.CommentFunc"

**Step 3: Write minimal implementation**

Add to `functions.go`:

```go
// CommentFunc returns an XML comment string
func CommentFunc() func(string) string {
	return func(text string) string {
		return fmt.Sprintf("<!-- %s -->", text)
	}
}
```

Add to imports in `functions.go` if not present: `"fmt"`

**Step 4: Register function in unpack.go NewTemplate**

In `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/unpack.go`, update `NewTemplate`:

```go
return template.New(name).
	Funcs(template.FuncMap{
		"modlet":    ModletFunc(fargs),
		"mult":      MultFunc(fargs),
		"output":    OutputFunc(fargs),
		"set":       SetFunc(fargs),
		"write":     WriteFunc(fargs),
		"xmlHeader": func() string { return xml.Header },
		"comment":   CommentFunc(), // Add this line
	}).
	Parse(string(data))
```

**Step 5: Run test to verify it passes**

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go test -v -run TestCommentFunc ./modlet/tests/`
Expected: PASS

**Step 6: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add modlet/functions.go modlet/unpack.go modlet/tests/functions_test.go
git commit -m "feat(modlet): add comment template function"
```

---

## Task 8: Add prob Template Function

**Files:**

- Modify: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/functions.go`
- Modify: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/unpack.go`
- Modify: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/tests/functions_test.go`

**Step 1: Write the failing test**

Add to `functions_test.go`:

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

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go test -v -run TestProbFunc ./modlet/tests/`
Expected: FAIL with "undefined: modlet.ProbFunc"

**Step 3: Write minimal implementation**

Add to `functions.go`:

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

**Step 4: Register function in unpack.go NewTemplate**

Add to the FuncMap in `NewTemplate`:

```go
"prob": ProbFunc(fargs),
```

**Step 5: Run test to verify it passes**

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go test -v -run TestProbFunc ./modlet/tests/`
Expected: PASS

**Step 6: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add modlet/functions.go modlet/unpack.go modlet/tests/functions_test.go
git commit -m "feat(modlet): add prob template function with 1.0 cap"
```

---

## Task 9: Integrate GameData into Template Execution

**Files:**

- Modify: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/unpack.go`
- Modify: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/functions.go`
- Create: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/tests/gamedata_integration_test.go`

**Step 1: Write the failing test**

```go
// gamedata_integration_test.go
package modlet_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/donovanmods/7dtd-modtools/modlet"
	"github.com/stretchr/testify/assert"
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

**Step 2: Run test to verify it fails**

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go test -v -run TestUnpackWithGameData ./modlet/tests/`
Expected: FAIL (GameData not available in template context)

**Step 3: Update FuncArgs and Unpack**

Add GameData to FuncArgs in `functions.go`:

```go
type FuncArgs struct {
	Output   string
	Gamedir  string
	ModInfo  *modinfo.ModInfo
	FBuffer  *FileBuffer
	GBuffer  *bytes.Buffer
	FBufMap  FileBufferMap
	Options  map[string]string
	GameData *gamedata.GameData // Add this
}
```

Add import: `"github.com/donovanmods/7dtd-modtools/gamedata"`

Update `Unpack` in `unpack.go` to load and pass game data:

```go
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

	fBufMap := make(FileBufferMap)
	gBuffer := bytes.NewBuffer(nil)
	fBuffer := &FileBuffer{
		Buffer: gBuffer,
		Writer: nil,
	}

	logger.Debug("processing template: %s", templateName)

	fargs := FuncArgs{
		Output:   output,
		Gamedir:  gamedir,
		ModInfo:  &modInfo,
		FBuffer:  fBuffer,
		GBuffer:  gBuffer,
		FBufMap:  fBufMap,
		GameData: gd, // Add this
		Options: map[string]string{
			"force": strconv.FormatBool(force),
		},
	}

	t, err := NewTemplate(tmpl, templateName, fargs)
	if err != nil {
		return fmt.Errorf("error parsing template %s: %w", templateName, err)
	}

	// Create template data with GameData
	templateData := struct {
		GameData *gamedata.GameData
	}{
		GameData: gd,
	}

	if err := t.ExecuteTemplate(gBuffer, templateName, templateData); err != nil {
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
```

Add import to `unpack.go`: `"github.com/donovanmods/7dtd-modtools/gamedata"`

**Step 4: Update gamedata.Load to use afero filesystem**

Update `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/gamedata/gamedata.go` to accept a filesystem:

```go
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

**Step 5: Update test setup to share filesystem**

Update `modlet_test.go` setup to also set gamedata.FS:

```go
func setup(t *testing.T) *assert.Assertions {
	t.Helper()

	logger.Testing = true

	// Use MemMapFs for testing
	modlet.FS = FS
	gamedata.FS = FS  // Add this

	mkTempDir(t)

	return assert.New(t)
}
```

Add import: `"github.com/donovanmods/7dtd-modtools/gamedata"`

**Step 6: Run test to verify it passes**

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go test -v -run TestUnpackWithGameData ./modlet/tests/`
Expected: PASS

**Step 7: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add modlet/ gamedata/
git commit -m "feat(modlet): integrate GameData into template execution"
```

---

## Task 10: Add Helper Functions to Template FuncMap

**Files:**

- Modify: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/unpack.go`
- Create: `/home/dyoung/Projects/mods/7dtd/7dtd-modtools/modlet/tests/template_helpers_test.go`

**Step 1: Write the failing test**

```go
// template_helpers_test.go
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

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go test -v -run TestTemplateHelperFunctions ./modlet/tests/`
Expected: FAIL with "function 'hasPrefix' not defined"

**Step 3: Add helper functions to NewTemplate FuncMap**

Update `NewTemplate` in `unpack.go`:

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

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go test -v -run TestTemplateHelperFunctions ./modlet/tests/`
Expected: PASS

**Step 5: Commit**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add modlet/unpack.go modlet/tests/template_helpers_test.go
git commit -m "feat(modlet): add helper functions to template FuncMap"
```

---

## Task 11: Run All Tests and Verify

**Step 1: Run all tests in 7dtd-gamedata**

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-gamedata && go test -v ./...`
Expected: All tests PASS

**Step 2: Run all tests in 7dtd-modtools**

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && go test -v ./...`
Expected: All tests PASS

**Step 3: Run lint checks**

Run: `cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools && task check`
Expected: No errors

**Step 4: Final commit for Phase 1**

```bash
cd /home/dyoung/Projects/mods/7dtd/7dtd-modtools
git add -A
git commit -m "chore: complete Phase 1 core infrastructure" --allow-empty
```

---

## Summary

Phase 1 establishes the foundation for template-based modlet generation:

1. **7dtd-gamedata enhancements**: Added `DropEvents()`, `GetProperty()`, `HasProperty()` methods to Block and EntityClass types
2. **GameData loader**: Created `gamedata` package to load and parse game XML files
3. **Helper functions**: Added `hasPrefix`, `hasSuffix`, `match`, `notMatch`, `multValue`, `probMult`
4. **Template functions**: Added `comment`, `prob`; completed `mult` with min/max bounds
5. **Template integration**: GameData now available as `.GameData` in templates

Templates can now:

- Access game blocks via `{{ range .GameData.Blocks }}`
- Filter by name patterns via `{{ if hasPrefix .Name "terr" }}`
- Get drop events via `{{ range .DropEvents "Harvest" }}`
- Generate XPath instructions via `{{ set }}`, `{{ mult }}`, `{{ prob }}`
- Add comments via `{{ comment "text" }}`
