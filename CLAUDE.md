# KV-Studio

A desktop GUI tool for managing and inspecting key-value databases.
This project is a complete, modern rewrite of "KV-Toolbox" with AI assistance.

## Supported Databases

- **PebbleDB** — Low-level LSM-Tree engine (derived from RocksDB)
- **BadgerDB** — Popular Go-native KV store
- **LevelDB** — Classic, simple, lightweight

## Architecture

The project follows the standard Go layout with `cmd/` and `internal/` separation:

```text
kv-studio/
├── cmd/kvstudio/           # Application entry point (main package)
├── internal/
│   ├── databases/          # Per-database adapters
│   │   ├── pebble/         # PebbleDB implementation
│   │   ├── badger/         # BadgerDB implementation
│   │   └── leveldb/        # LevelDB implementation
│   ├── ui/                 # UI layer (Fyne)
│   │   ├── mainwindow/     # Main window
│   │   ├── components/     # Reusable widgets
│   │   └── theme/          # Theme and styling
│   ├── logic/              # Business logic (filter, search, ...)
│   ├── config/             # Read/write user settings
│   └── utils/              # Generic helpers
├── docs/                   # Project documentation
└── .claude/                # Claude Code configuration (skills, agents, settings)
```

## Core Architectural Principle: One Interface for All Databases

Every database adapter must implement a shared interface called `KVStore`:

```go
type KVStore interface {
    Open(path string) error
    Close() error
    Get(key []byte) ([]byte, error)
    Set(key, value []byte) error
    Delete(key []byte) error
    Iterate(prefix []byte, fn func(key, value []byte) bool) error
    Stats() Stats
}
```

This pattern keeps the UI and logic layers completely decoupled from any specific database implementation.

## Common Commands

| Command                    | Purpose                            |
| -------------------------- | ---------------------------------- |
| `go build ./cmd/kvstudio`  | Build the binary                   |
| `go run ./cmd/kvstudio`    | Run the app in development mode    |
| `go test ./...`            | Run the full test suite            |
| `go vet ./...`             | Static analysis                    |
| `gofmt -w .`               | Format code                        |

## Coding Standards

- **Language**: Go 1.22+
- **GUI**: Fyne v2
- **Error handling**: Use `errors.Is`/`errors.As` and wrap with `fmt.Errorf("...: %w", err)`
- **Naming**: UpperCamelCase for exported, lowerCamelCase for internal
- **Tests**: Every `internal/...` package must have an accompanying `*_test.go` file
- **Comments**: Only where the "why" is non-obvious — not to describe "what"
- **No `panic` in the main path** — propagate errors and let the UI decide

## Project Conventions

1. Each database adapter lives in its own package (`internal/databases/<name>`).
2. The UI must never import database packages directly — always go through the interface.
3. User settings are stored in `~/.kvstudio/config.json`.
4. Logs are written to `~/.kvstudio/logs/` with daily rotation.
5. No secrets or machine-specific paths may be committed to the repo.

## Current Status

The project is in the **bootstrap** phase. The base structure has been created;
no core modules have been implemented yet. The first step is to implement the
`KVStore` interface and the three database adapters.

## Notes for AI Assistants

- Before making large architectural changes, enter **plan mode**.
- Before every commit, ensure `go vet` and `go test ./...` pass.
- When a file in `internal/databases/<x>/` changes, verify the other adapters
  still satisfy the interface uniformly.
- For specialized tasks, use the subagents defined in `.claude/agents/`.

## Commit Conventions

- **Do not add AI attribution to commits.** No `Co-Authored-By: Claude ...`,
  no `Generated with ...` trailers, no signatures referencing any AI tool.
- Keep commit messages short, plain English, and human-sounding. One line is
  usually enough; two if context is needed.
- Never modify `git config user.name` or `user.email`. Commits stay under the
  repo owner's identity.
