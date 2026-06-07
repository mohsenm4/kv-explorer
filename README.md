# KV-Explorer

A desktop GUI for inspecting, editing, and comparing embedded key-value
databases. One window, one workflow for all three of the popular
Go-native KV engines.

A modern rewrite of KV-Toolbox, built with [Fyne](https://fyne.io) and
heavy use of AI-assisted development.

## Supported engines

| Engine    | Notes                                          |
| --------- | ---------------------------------------------- |
| PebbleDB  | LSM-tree engine derived from RocksDB           |
| BadgerDB  | Popular Go-native KV store                     |
| LevelDB   | Classic lightweight engine (requires CGO)      |

## Features

- **Multiple databases at once** — open Pebble, Badger, and LevelDB
  stores side-by-side in tabs, each color-coded by engine.
- **Prefix tree + key table** — browse keys hierarchically on the left,
  see key / preview / size in a sortable table on the right.
- **Value editor** — view and edit values as UTF-8, raw Hex, or
  pretty-printed JSON. Unsaved changes are clearly marked.
- **Powerful filtering** — filter by key prefix, key/value substring, or
  key regex; constraints AND together and update as you type.
- **Add / edit / delete keys** — with confirmation for destructive ops
  and batch selection for bulk delete.
- **Read-only mode** — open a database without risk of accidental
  writes.
- **Persistent state** — window size, theme, and a recent-databases
  list survive restarts (stored in `~/.kvexplorer/`).
- **Keyboard-first** — every common action has a shortcut
  (`Cmd/Ctrl+O` open, `Cmd/Ctrl+F` filter, `F2` edit, `Delete` remove,
  `F5` refresh, …).
- **Light / dark / system theme** — toggle from the status bar.

## Install

### Download a release

Pre-built binaries for macOS (Intel/ARM), Linux, and Windows are
published on the [Releases](https://github.com/mohsenm4/kv-explorer/releases)
page. Download the archive for your platform, extract, and run
`kvexplorer`.

### Build from source

Requirements:

- Go 1.25 or newer
- A C toolchain with CGO enabled (LevelDB requires it)
- Linux only: `gcc`, `libgl1-mesa-dev`, `xorg-dev` for the Fyne GUI

```bash
git clone https://github.com/mohsenm4/kv-explorer.git
cd kv-explorer
go build ./cmd/kvexplorer
./kvexplorer
```

Or run without building:

```bash
go run ./cmd/kvexplorer
```

## Usage at a glance

1. Launch the app — the **Welcome** screen lists recent databases.
2. **Open Database…** — pick the engine, point at a directory, choose
   whether to open in a new tab and whether to mount read-only.
3. Browse the prefix tree on the left or use the filter bar to narrow
   the table.
4. Click a row to load its value into the editor; switch between
   UTF-8 / Hex / JSON view; **Save** to commit, **Cancel** to discard.
5. **Add / Edit / Delete** from the toolbar or via shortcuts.

Full keyboard reference and behavior details live in
[CLAUDE.md](./CLAUDE.md).

## Documentation

- [CLAUDE.md](./CLAUDE.md) — project rules, full feature reference, and
  conventions.
- [docs/architecture.md](./docs/architecture.md) — code map: where every
  layer lives and where to change what.
- [docs/design/](./docs/design/) — visual tokens, Fyne mapping, and
  accessibility rules. The Figma source is linked from the design
  README.

## Development

```bash
go test -race -cover ./...   # tests
go vet ./...                 # static analysis
gofmt -w .                   # format
```

CI runs lint, test, and build on Linux, macOS, and Windows for every
push and pull request (see [`.github/workflows/ci.yml`](.github/workflows/ci.yml)).
Tagged releases (`v*`) build and publish cross-platform binaries
automatically.

Contributions are welcome. Please open an issue first for anything
larger than a small fix so we can agree on the approach.

## License

TBD
