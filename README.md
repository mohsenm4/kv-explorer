# KV-Explorer

A desktop GUI tool for managing, inspecting, and comparing key-value databases.
A complete, modern rewrite of KV-Toolbox with AI assistance.

## Supported Databases

- PebbleDB
- BadgerDB
- LevelDB

## Requirements

- Go 1.22 or newer
- CGO enabled (required by LevelDB)
- macOS / Linux / Windows

## Quick Start

```bash
git clone https://github.com/mohsenm4/kv-explorer.git
cd kv-explorer
go mod tidy
go run ./cmd/kvexplorer
```

## Project Structure

The full structure and conventions are documented in [CLAUDE.md](./CLAUDE.md).

## Design

UI design lives in Figma:
<https://www.figma.com/design/QAe3rjM2Mrb5dCUZgDr4Dm/Untitled?node-id=1-1997&p=f>

Engineering reference (tokens, Fyne mapping, accessibility rules) is in
[`docs/design/`](./docs/design/).

## Development

- Run tests: `go test ./...`
- Build: `go build ./cmd/kvexplorer`
- Cross-platform build: see the `/build-cross-platform` skill

## License

TBD
