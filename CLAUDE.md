# KV-Explorer

A desktop GUI tool for managing and inspecting key-value databases.
This project is a complete, modern rewrite of "KV-Toolbox" with AI assistance.

> **Process**: features come first, code shape follows. When the UI needs a
> new capability, we add the smallest backend piece that supports it. We do
> not prescribe folder layout up-front; we let it grow from what the
> features actually need.

## Supported Databases

- **PebbleDB** — Low-level LSM-Tree engine (derived from RocksDB)
- **BadgerDB** — Popular Go-native KV store
- **LevelDB** — Classic, simple, lightweight

## What the User Can Do (UI Features)

The desktop app's user-facing capabilities. Backend exists to serve these —
nothing else.

### Screens

- **Welcome** — Shown when no database is open. Shows the product name,
  a tagline, an **Open Database…** action, an **Open Recent** dropdown,
  and a list of recently opened databases.
- **Main window** — Shown once a database is open:
  - Top toolbar (see below)
  - One tab per open database, each with the engine's accent color
  - Left pane: prefix-based tree of keys
  - Center: table of key / value preview / size
  - Bottom of center: value editor (collapsible)
  - Bottom of window: status bar (see below)
- **Open Database dialog** — Engine selector (Pebble / Badger / LevelDB),
  path picker, "open in new tab" toggle, "read-only" toggle.
- **Add key dialog** — Key field, value editor with UTF-8 / Hex / JSON
  format toggle, byte-size readout.
- **Edit key dialog** — Same as Add, prefilled, with validation for
  duplicate keys.
- **Delete confirmation** — Shows the key being deleted; destructive style.
- **Settings dialog** — Tabbed: Appearance (theme, density), General
  (paths, log level), Editor (font, pretty-print), Shortcuts (key bindings),
  About.

### Toolbar actions (top of main window)

| Open | Close | ─ | Add | Edit | Delete | ─ | Refresh | Settings |

### Status bar (bottom of main window)

Engine dot + name | key count | on-disk size | open path (truncated) | theme toggle | settings icon

### Filtering and matching

A filter input above the key table. Filters by **key prefix**, **key substring**,
**value substring**, or **key regex**. Multiple constraints are AND-ed. The
table updates as the user types (debounced).

### Editing values

The value editor supports UTF-8, raw Hex view, and JSON pretty-print.
Save and Cancel buttons commit or discard the edit. Unsaved changes show a
visible marker.

### Batch operations

Multiple keys can be selected in the table for batch delete. Batch operations
stop at the first failing key and report which one.

### Database comparison

Multiple databases can be open simultaneously, one tab each. Switching tabs
swaps the entire body (tree + table + editor) for the other engine's data.

### Persistent state

User settings survive restarts:

- Window size and position
- Theme choice (light / dark / system)
- Recently opened databases (path + engine + timestamp)

Settings live at `~/.kvexplorer/config.json`. Logs live at
`~/.kvexplorer/logs/` with daily rotation.

### Keyboard shortcuts

- `Cmd/Ctrl + O` — Open Database
- `Cmd/Ctrl + W` — Close current tab
- `Cmd/Ctrl + N` — Add key
- `Cmd/Ctrl + F` — Focus filter
- `Cmd/Ctrl + S` — Save value edits
- `Delete` — Delete selected key (with confirm)
- `F2` — Edit selected key
- `F5` — Refresh
- `Cmd/Ctrl + ,` — Settings
- `Cmd/Ctrl + Tab` — Cycle tabs
- `Esc` — Close dialog / clear filter / cancel edit

Full visual design (colors, components, screen layouts) is in
[`docs/design/`](./docs/design/).

## Core Principle: One Interface for All Databases

Every database adapter implements the same Go interface:

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

The UI sees only this interface. Adding a fourth engine means adding an
adapter — no UI changes.

## Project Layout

Standard Go layout (`cmd/` + `internal/`). The current package map and the
import rules between layers live in [`docs/architecture.md`](./docs/architecture.md);
that document is authoritative. This file states principles, not paths.

## Common Commands

| Command                      | Purpose                         |
| ---------------------------- | ------------------------------- |
| `go build ./cmd/kvexplorer`  | Build the binary                |
| `go run ./cmd/kvexplorer`    | Run the app in development mode |
| `go test ./...`              | Run the full test suite         |
| `go vet ./...`               | Static analysis                 |
| `gofmt -w .`                 | Format code                     |

## Coding Standards

- **Language**: Go 1.22+
- **GUI**: Fyne v2
- **Error handling**: Use `errors.Is`/`errors.As` and wrap with `fmt.Errorf("...: %w", err)`
- **Naming**: UpperCamelCase for exported, lowerCamelCase for internal
- **Tests**: Every `internal/...` package must have an accompanying `*_test.go` file
- **Comments**: Only where the "why" is non-obvious — not to describe "what"
- **No `panic` in the main path** — propagate errors and let the UI decide

## Conventions

1. **Adapters are siblings.** None imports another. Each implements
   `KVStore` independently.
2. **The UI never imports an adapter directly** — only the `KVStore` interface.
3. **Theme tokens are the source of truth for color and size.** Never
   hardcode in a widget.
4. **No secrets or machine-specific paths** in the repo.

## Current Status

Foundation in place: theme system, welcome screen, `KVStore` interface,
three adapter stubs, filter + batch logic, config load/save. Next: real
adapter implementations and wiring the Open Database flow to the UI.

## Notes for AI Assistants

- Before making large architectural changes, enter **plan mode**.
- Before every commit, ensure `go vet` and `go test ./...` pass.
- When a file in `internal/databases/<x>/` changes, verify the other adapters
  still satisfy the interface uniformly.
- **Do not blindly follow patterns from the previous KV-Toolbox codebase.**
  Decide what's right for KV-Explorer from first principles — let user-facing
  features dictate what backend code exists.
- For specialized tasks, use the subagents defined in `.claude/agents/`.

## Commit Conventions

- **Do not add AI attribution to commits.** No `Co-Authored-By: Claude ...`,
  no `Generated with ...` trailers, no signatures referencing any AI tool.
- Keep commit messages short, plain English, and human-sounding. One line is
  usually enough; two if context is needed.
- Never modify `git config user.name` or `user.email`. Commits stay under the
  repo owner's identity.
