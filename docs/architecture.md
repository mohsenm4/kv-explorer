# Architecture

This document explains how KV-Explorer is organized and how data flows through
the system. It is meant for new contributors and stakeholders who want to
understand the project without reading the source.

## Layered View

KV-Explorer is split into four layers. Each layer depends only on the layer
below it — never the other way around.

```text
┌─────────────────────────────────────────┐
│  UI layer  (internal/ui, Fyne v2)        │
│  windows, widgets, theme                 │
└──────────────────┬──────────────────────┘
                   │ calls
┌──────────────────▼──────────────────────┐
│  Logic layer  (internal/logic)           │
│  filter, search, batch ops, validation   │
└──────────────────┬──────────────────────┘
                   │ uses
┌──────────────────▼──────────────────────┐
│  KVStore interface  (internal/databases) │
│  shared contract for every adapter       │
└──────────────────┬──────────────────────┘
                   │ implemented by
┌──────────────────▼──────────────────────┐
│  Adapters  (pebble / badger / leveldb)   │
│  per-engine code, isolated in subpkgs    │
└─────────────────────────────────────────┘
```

The key idea: the UI and logic layers do not know which database is in use.
They only talk to the `KVStore` interface. Swapping engines (or adding new
ones) does not require any change above the adapter layer.

## The KVStore Contract

The single most important file in the codebase is `internal/databases/kvstore.go`.
It defines the interface every adapter must implement:

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

Each method has a specific contract:

| Method   | Contract                                                                |
| -------- | ----------------------------------------------------------------------- |
| Open     | Opens the store at `path`. Creates the directory if missing.            |
| Close    | Releases all resources. Safe to call multiple times.                    |
| Get      | Returns the value for `key`. Returns `ErrNotFound` if absent.           |
| Set      | Writes `key` -> `value`. Overwrites if `key` already exists.            |
| Delete   | Removes `key`. No error if `key` did not exist.                         |
| Iterate  | Scans keys with the given prefix. Stop early when `fn` returns `false`. |
| Stats    | Returns key count and on-disk size at the moment of the call.           |

All adapters must agree on error semantics. A missing key returns the same
sentinel error across PebbleDB, BadgerDB, and LevelDB.

## Database Comparison

| Engine    | Best for                       | Notes                                     |
| --------- | ------------------------------ | ----------------------------------------- |
| PebbleDB  | Write-heavy, large datasets    | LSM-tree, pure Go, derived from RocksDB   |
| BadgerDB  | Mixed workloads, embedded use  | LSM-tree, pure Go, value-log architecture |
| LevelDB   | Read-heavy, small to medium    | Classic LSM-tree, requires CGO            |

KV-Explorer does not pick a winner — it lets the user compare them side by side.

## Data Flow Example: Setting a Key

When a user sets a key from the UI:

1. The UI calls `logic.WriteKey(store, key, value)` with the active store.
2. The logic layer validates input (size limits, encoding) and calls
   `store.Set(key, value)`.
3. The adapter translates the call into the engine's native API
   (e.g. `pebble.DB.Set`, `badger.Txn.Set`, `leveldb.DB.Put`).
4. On success, the UI refreshes the visible key list.
5. On failure, the adapter returns a wrapped error; the UI shows it.

The UI never imports any of the adapter packages directly. It only sees the
`KVStore` interface returned from a factory in `internal/databases`.

## Extension Points

The architecture is designed so that the following changes touch as little
code as possible:

- **Adding a fourth database**: implement `KVStore` in a new subpackage under
  `internal/databases/<name>/`. Register it in the factory. No UI changes.
- **Adding a new filter type**: implement it in `internal/logic/`. The UI
  picks it up through the existing filter registry.
- **Theming**: edit `internal/ui/theme/`. No widget code changes.
- **New widgets**: add them in `internal/ui/components/`. Reusable across
  windows.

## Configuration

User settings (window size, last opened path, theme choice) are stored in
`~/.kvexplorer/config.json`. Logs go to `~/.kvexplorer/logs/`. Neither location
is touched by the test suite.

## What This Document Is Not

- It is not API documentation. See the Go doc comments for that.
- It is not a roadmap. See GitHub issues.
- It is not exhaustive. See `CLAUDE.md` for conventions and commands.
