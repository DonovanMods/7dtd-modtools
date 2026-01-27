# 7dtd-modtools v1.0 Product Requirements Document

**Date:** 2026-01-26
**Status:** Draft
**Author:** Donovan C. Young (with Claude)

---

## 1. Overview & Goals

**7dtd-modtools** is a Go CLI application for creating, templating, and managing modlets for the game "7 Days to Die". It replaces the Ruby/Python scripts currently in `donovan-7d2d-modlets/scripts/` with a unified, template-based approach.

### Primary Goals

1. **Full Script Replacement**: Reproduce all functionality of existing scripts (`mk_lessgrind.rb`, `mk_parts_recipes.rb`, `mk_mod_schematic_recipes.rb`, `mkbundle.py`) through Go templates with embedded logic.

2. **Template-Based Modlet Generation**: Allow modders to create `.tmpl` files that combine static XPath patches with dynamic logic (iteration over game data, conditional multipliers, pattern matching).

3. **Game Data Integration**: Parse and expose game XML files (blocks, items, entityclasses, recipes) as template context, enabling templates to generate modlets that adapt to game updates.

4. **Bidirectional Conversion**: Support both `pack` (existing mod → template) and `unpack` (template → modlet) workflows.

5. **Built-in Validation**: Validate generated XPath expressions against actual game XML during unpack, warning when expressions won't match.

### Target Audience

Experienced 7D2D modders familiar with XPath patching who want to:

- Create reusable, parameterized mod templates
- Auto-generate modlets from game data (like lessgrind)
- Maintain modlets across game version updates more easily
- Share template logic with other modders

### Non-Goals (v1.0)

- GUI interface
- Mod.io/Nexus integration
- ZIP file generation (users can use existing tools)
- Template inheritance or multi-file template packages

---

## 2. Core Features

### Commands

**`modlet new <name>`** (exists)

- Creates skeleton modlet directory with ModInfo.xml, Config/, README.md
- Flags: `--output`, `--force`

**`modlet pack <mod-dir>`** (exists, needs enhancement)

- Converts existing modlet directory into `.tmpl` file
- Wraps each file with `{{ output "path" }}` ... `{{ write }}` directives
- Flags: `--output`, `--compress` (gzip), `--force`

**`modlet unpack <templates...>`** (exists, needs enhancement)

- Executes templates against game data context to generate modlet directories
- Validates XPath expressions against game XML, warns on mismatches
- Flags: `--output`, `--gamedir` (required), `--force`, `--no-validate`, `--offline`

**`modlet validate <modlet-dir>`** (new)

- Standalone validation of existing modlet against game XML
- Reports XPath expressions that don't match any game elements
- Flags: `--gamedir` (required)

### Template System

Templates use Go's `text/template` with custom functions:

| Function                            | Purpose                             |
| ----------------------------------- | ----------------------------------- |
| `modlet "name"`                     | Initialize modlet, create directory |
| `output "path"`                     | Set current output file             |
| `write`                             | Flush buffer to file                |
| `set "xpath" "value"`               | Generate `<set>` instruction        |
| `mult "xpath" by=N [min=X] [max=Y]` | Multiply value with bounds          |
| `xmlHeader`                         | Output XML declaration              |

**New functions needed:**

| Function                          | Purpose                                 |
| --------------------------------- | --------------------------------------- |
| `prob "xpath" by=N`               | Multiply probability, cap at 1.0        |
| `skip "pattern"`                  | Mark patterns to exclude from iteration |
| `comment "text"`                  | Insert XML comment                      |
| `download url path [sha256=hash]` | Fetch external file during unpack       |

---

## 3. Game Data Context

### Pre-loaded Game Data

When `unpack` runs, it parses game XML files from `--gamedir` and exposes them as template variables:

```go
// Available in all templates as .GameData
type GameData struct {
    Blocks        []Block        // From blocks.xml
    Items         []Item         // From items.xml
    EntityClasses []EntityClass  // From entityclasses.xml
    Recipes       []Recipe       // From recipes.xml
    Loot          []LootGroup    // From loot.xml
}
```

### Template Access Patterns

**Iterating over game data:**

```go
{{- range .GameData.Blocks }}
  {{- if hasPrefix .Name "terr" }}
    {{- range .DropEvents "Harvest" }}
      {{ set (printf "//block[@name='%s']/drop[@event='Harvest']/@count" $.Name) (mult .Count 1.5) }}
    {{- end }}
  {{- end }}
{{- end }}
```

**Helper functions for common patterns:**

| Function                 | Purpose                          |
| ------------------------ | -------------------------------- |
| `hasPrefix s prefix`     | String prefix check              |
| `hasSuffix s suffix`     | String suffix check              |
| `match s pattern`        | Regex match                      |
| `notMatch s pattern`     | Regex negative match             |
| `printf format args...`  | String formatting                |
| `mult value factor`      | Multiply (returns string)        |
| `multRound value factor` | Multiply and round to int        |
| `probMult value factor`  | Multiply probability, cap at 1.0 |

### Block/Entity Accessor Methods

```go
// On Block type
func (b Block) DropEvents(event string) []DropEvent  // Filter by Harvest/Destroy
func (b Block) Property(name string) string          // Get property value
func (b Block) HasProperty(name string) bool         // Check property exists

// On EntityClass type
func (e EntityClass) DropEvents(event string) []DropEvent
func (e EntityClass) PassiveEffects() []PassiveEffect
```

---

## 4. Validation

### Validation During Unpack

After template execution generates XPath instructions, modtools validates each expression against the loaded game XML:

**Validation Process:**

1. Parse generated modlet XML files
2. Extract all XPath expressions from `<set>`, `<append>`, `<remove>`, etc.
3. Execute each XPath against the corresponding game XML file
4. Report expressions that match zero elements

**Output Example:**

```
Unpacking lessgrind.tmpl...
  Writing ModInfo.xml
  Writing Config/blocks.xml (847 instructions)
  Writing Config/entityclasses.xml (42 instructions)

Validation Results:
  ✓ Config/blocks.xml: 845/847 XPaths valid
  ⚠ Config/blocks.xml:234 - no match: //block[@name='oldRemovedBlock']/drop[@event='Harvest']/@count
  ⚠ Config/blocks.xml:512 - no match: //block[@name='renamedBlock']/property[@name='OldProp']/@value
  ✓ Config/entityclasses.xml: 42/42 XPaths valid

Generated: modlets/donovan-lessgrind/ (2 warnings)
```

**Validation Behavior:**

- Warnings don't stop generation (game tolerates non-matching XPaths)
- `--no-validate` flag skips validation entirely
- `--strict` flag (optional) treats warnings as errors

### Standalone Validate Command

```bash
7dmt modlet validate modlets/donovan-lessgrind/ --gamedir ~/.local/share/Steam/.../Data/Config
```

Useful for:

- Checking existing modlets after game updates
- CI/CD pipelines for mod repositories
- Debugging XPath expressions

---

## 5. AiO Bundling via Templates

### Template Composition Approach

Instead of a dedicated `bundle` command, the All-in-One modlet is itself a template that includes other templates' output. This uses Go's `text/template` built-in `template` action.

**aio.tmpl structure:**

```go
{{- modlet "donovan-aio" -}}

{{/* Define which modlets to include */}}
{{- $modlets := list
    "betterbatons" "betterblades" "betterbrawler" "betterbridges"
    "betterbuffs" "bettercement" "betterdyes" "betterpowertools"
    "bettertraps" "lessgrind" "longerlootbags" "megastacks"
    "morebooks" "morelootbags" "moreperks" "pickmeup"
-}}

{{- output "Config/blocks.xml" -}}
{{ xmlHeader }}
<configs>
{{- range $modlets }}
  <!-- Included from donovan-{{ . }} -->
  {{- includeConfig . "blocks.xml" }}
{{- end }}
</configs>
{{- write -}}

{{/* Repeat for items.xml, recipes.xml, etc. */}}
```

### New Template Functions for Bundling

| Function                    | Purpose                                                                                |
| --------------------------- | -------------------------------------------------------------------------------------- |
| `list items...`             | Create a slice for iteration                                                           |
| `includeConfig modlet file` | Include another modlet's config file content (strips XML declaration and root element) |
| `includeTemplate name`      | Execute another .tmpl and capture output                                               |

### Workflow

1. Individual modlet templates (lessgrind.tmpl, betterbatons.tmpl, etc.) generate standalone modlets
2. The aio.tmpl references these and merges their configs
3. Running `unpack aio.tmpl` produces the bundled donovan-aio modlet

This keeps bundling logic in templates (transparent, versionable) rather than hidden in application code.

---

## 6. External File Downloads

### Use Case

Some modlets require compiled files (DLLs, assemblies) or binary assets that cannot be generated from templates. These files are typically:

- Harmony patches (C# DLLs)
- Custom textures or icons
- Pre-compiled assemblies from other mod authors

### Template Syntax

New `download` function fetches external files during unpack:

```go
{{- modlet "donovan-customfeature" -}}

{{/* Download a DLL from GitHub releases */}}
{{ download "https://github.com/user/repo/releases/download/v1.0/CustomMod.dll" "CustomMod.dll" }}

{{/* Download to a subdirectory */}}
{{ download "https://example.com/assets/icon.png" "Resources/icon.png" }}

{{/* Download with checksum verification */}}
{{ download "https://example.com/SomeAssembly.dll" "SomeAssembly.dll" sha256="a1b2c3d4..." }}
```

### Function Signature

| Function                          | Purpose                                         |
| --------------------------------- | ----------------------------------------------- |
| `download url path [sha256=hash]` | Fetch URL, save to path relative to modlet root |

### Behavior

- Downloads occur during `unpack` after template execution
- Files are placed relative to the modlet output directory
- Parent directories created automatically
- **Checksum verification** (optional): If `sha256=` provided, verify hash matches; fail if mismatch
- **Caching**: Downloaded files cached in `~/.cache/7dtd-modtools/` to avoid re-downloading
- **Offline mode**: `--offline` flag uses only cached files, fails if missing
- **Failure handling**: Download failure stops unpack with clear error message

### Security Considerations

- Only HTTPS URLs allowed (reject HTTP)
- Warn user before downloading from non-GitHub/non-trusted domains
- Checksum strongly recommended for DLLs in documentation

---

## 7. Implementation Roadmap

### Phase 1: Core Infrastructure (Foundation)

**Goal:** Complete the template engine with game data support

- [ ] Finish `mult` function (apply min/max bounds)
- [ ] Add `prob` function for probability multiplication with cap
- [ ] Add `comment` function for XML comments
- [ ] Implement game XML parsing in 7dtd-gamedata (Items, Recipes, Loot)
- [ ] Create GameData context struct that loads all game files
- [ ] Pass GameData to template execution in unpack
- [ ] Add string helper functions (hasPrefix, hasSuffix, match, notMatch)

### Phase 2: Validation

**Goal:** Validate generated XPaths against game data

- [ ] Extract XPath expressions from generated XML
- [ ] Execute XPaths against loaded game XML
- [ ] Report non-matching expressions with line numbers
- [ ] Add `--no-validate` and `--strict` flags
- [ ] Implement standalone `modlet validate` command

### Phase 3: External Downloads

**Goal:** Support external file inclusion

- [ ] Implement `download` template function
- [ ] Add download caching in `~/.cache/7dtd-modtools/`
- [ ] Implement SHA256 checksum verification
- [ ] Add `--offline` flag support
- [ ] HTTPS-only enforcement with domain warnings

### Phase 4: Template Migration

**Goal:** Recreate existing modlets as templates

- [ ] Create lessgrind.tmpl (most complex - blocks + entityclasses)
- [ ] Create craftableparts.tmpl
- [ ] Create modschematics.tmpl
- [ ] Create simple modlet templates (betterbatons, megastacks, etc.)
- [ ] Create aio.tmpl with bundling logic
- [ ] Implement `list`, `includeConfig`, `includeTemplate` functions
- [ ] Verify generated output matches expected modlet structure

### Phase 5: Polish

**Goal:** Production-ready for other modders

- [ ] Comprehensive error messages with template line context
- [ ] Documentation: README, template authoring guide, examples
- [ ] `--verbose` output improvements
- [ ] Performance optimization for large game XML files

---

## 8. Success Criteria

### v1.0 Release Criteria

1. **Template Completeness**: All existing auto-generated modlets (lessgrind, craftableparts, modschematics) can be reproduced via templates

2. **Game Data Access**: Templates can iterate over blocks, items, entityclasses, and recipes with filtering and pattern matching

3. **Validation Works**: Running `unpack` reports XPath mismatches; `validate` command works on any modlet

4. **AiO Generation**: The aio.tmpl produces a bundled modlet equivalent to current donovan-aio

5. **Bidirectional**: `pack` creates valid templates from existing modlets; `unpack` generates valid modlets from templates

6. **External Downloads**: Templates can include external files via `download` function with caching and checksum support

7. **Documentation**: README covers installation, all commands, and template authoring basics; at least 3 example templates included

### Out of Scope for v1.0

- Template inheritance/packages (single-file templates only)
- ZIP/archive generation
- GUI or web interface
- Mod hosting platform integration
- Windows installer

---

## 9. Dependencies

### Internal

- **7dtd-gamedata library**: Needs Items, Recipes, Loot structs added (Blocks and EntityClasses exist)

### External

- **Game installation**: Required for `--gamedir` to parse game XML
- **Network access**: Required for `download` function (optional, can use `--offline`)

---

## 10. Risks & Mitigations

| Risk                             | Mitigation                                                                 |
| -------------------------------- | -------------------------------------------------------------------------- |
| Game XML structure changes       | Validation catches mismatches; templates can be updated                    |
| Complex Ruby logic hard to port  | Embedded Go functions can handle any complexity                            |
| Performance with large XML       | Lazy loading or caching if needed                                          |
| External download URLs break     | Checksum verification detects tampering; caching provides offline fallback |
| Security concerns with downloads | HTTPS-only, domain warnings, checksum verification                         |

---

## Appendix A: Existing Script Analysis

### Scripts to Replace

| Script                        | Complexity | Template Approach                                               |
| ----------------------------- | ---------- | --------------------------------------------------------------- |
| `mk_lessgrind.rb`             | High       | Iterate blocks/entities, conditional multipliers, skip patterns |
| `mk_parts_recipes.rb`         | Medium     | Filter items by name pattern, generate recipes                  |
| `mk_mod_schematic_recipes.rb` | Medium     | Filter items, transform to schematics                           |
| `mkbundle.py`                 | Medium     | Template composition via `includeConfig`                        |

### Key Logic from lessgrind

- Default multiplier: 1.5x
- Terrain/ore blocks: 3x multiplier (1.5 \* 2)
- Tree stumps/cactus: 1.5x multiplier
- Skip patterns: `planted*`, `*Shapes`, `*Twitch`
- Probability cap: 1.0 maximum
- Special cases: dew collectors, cars, wood log fuel value

All this logic will be encoded in Go template functions or expressed directly in template control flow.
