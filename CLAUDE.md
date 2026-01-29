# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

7dtd-modtools is a CLI tool suite for "7 Days to Die" mod developers. It manages "modlets" - a standardized format combining game mod configurations with Go text/template templating.

**Core Commands:**

- `modlet new <name>` - Create new modlet skeleton
- `modlet pack <mod-dir>` - Pack mod directory into portable `.tmpl` template
- `modlet unpack <templates>` - Unpack templates to generate mod files

## Development Commands

```bash
# Run tests with coverage
make test                    # or: go test -v -cover ./...

# Lint and format
make check                   # Runs trunk check (depends on format)
make format                  # Runs trunk fmt (depends on tidy)

# Build
make build                   # All platforms (runs check + test first)
go build -o bin/7dmt ./main.go  # Quick local build

# Clean
make clean                   # Remove build artifacts and caches
```

**Run a single test:**

```bash
go test -v -run TestName ./modlet/tests/
```

## Architecture

main.go # Entry point → cmd.Execute()

```text
├── cmd/                     # CLI layer (Cobra commands)
│   ├── root.go              # Root command, persistent flags, config loading
│   └── modlet/              # Modlet command group
│       ├── modlet.go        # Parent command with shared flags
│       ├── new/new.go       # modlet new subcommand
│       ├── pack/pack.go     # modlet pack subcommand
│       └── unpack/unpack.go # modlet unpack subcommand
├── modlet/                  # Business logic
│   ├── modlet.go            # CmdArgs type and command dispatchers
│   ├── new.go               # CreateNew implementation
│   ├── pack.go              # Pack implementation
│   ├── unpack.go            # Unpack implementation
│   ├── functions.go         # Template functions (modlet, output, write, set, mult)
│   └── tests/               # Unit tests
└── lib/logger/              # pterm-based logging
```

**Key Pattern:** CLI commands in `cmd/` delegate to business logic in `modlet/` via `CmdArgs` struct methods.

## Testing Conventions

Tests use an in-memory filesystem (`afero.MemMapFs`) for isolation:

```go
func setup(t *testing.T) *assert.Assertions {
    logger.Testing = true           // Enables panic instead of os.Exit
    modlet.FS = &afero.Afero{Fs: afero.NewMemMapFs()}
    return assert.New(t)
}
```

Always set `logger.Testing = true` in tests to prevent `os.Exit()` calls.

## Key Dependencies

- **Cobra/Viper** - CLI framework and configuration
- **Afero** - Filesystem abstraction (swap real FS for memory FS in tests)
- **7dtd-gamedata** - Game data parsing library (ModInfo, modlet structures) - see below
- **pterm** - Terminal output formatting

## Related Project: 7dtd-gamedata

The `github.com/donovanmods/7dtd-gamedata` library is a sibling project located at `~/Projects/apps/7dtd-gamedata`. It provides core game data parsing (ModInfo, Modlet structs) and is primarily built for use by this tool. Both projects share the same owner, so the library can be modified alongside this tool as needed.

## Template Functions

The template system (`functions.go`) provides:

- `{{ modlet "name" }}` - Initialize modlet
- `{{ output "path" }}` - Set output file path
- `{{ write }}` - Flush buffer to file
- `{{ set "xpath" "value" }}` - XPath modification
- `{{ mult "xpath" "by=1.5" }}` - Value multiplication (min/max bounds not yet implemented)
