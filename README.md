# KV-Studio

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
git clone https://github.com/mohsenm4/kv-studio.git
cd kv-studio
go mod tidy
go run ./cmd/kvstudio
```

## Project Structure

The full structure and conventions are documented in [CLAUDE.md](./CLAUDE.md).

## Development

- Run tests: `go test ./...`
- Build: `go build ./cmd/kvstudio`
- Cross-platform build: see the `/build-cross-platform` skill

## License

TBD
