# KV-Explorer — Architecture & Maintenance Guide

A short map of where things live and where to change them. Pair this with
`docs/design/spec.md` (visual tokens) and `CLAUDE.md` (project rules).

---

## Layers at a glance

```
cmd/kvexplorer/main.go           # entry point — never put logic here
└── internal/
    ├── kvstore/                 # storage — Store interface + adapters
    │   ├── store.go             #   interface, types, OpenOptions
    │   ├── pebble/  badger/  leveldb/   # each implements Store
    ├── config/                  # ~/.kvexplorer/config.json
    ├── app/                     # wiring — only place that imports adapters
    │   ├── app.go               #   OpenStore dispatcher, CountKeys
    │   └── session.go           #   Session, KeyMeta, makePreview
    └── ui/                      # Fyne — never imports adapters directly
        ├── theme/theme.go       #   the only source of color/size tokens
        ├── state.go             #   AppState — single source of truth
        ├── ui.go                #   Run() — wire only
        ├── mainwindow.go        #   main page composition
        ├── welcome.go           #   welcome page
        ├── toolbar.go / tabstrip.go / statusbar.go / menu.go
        ├── keytable.go / prefixtree.go / editor.go / filter.go
        ├── key_dialog.go / open_dialog.go / delete_confirm.go
        ├── settings_dialog.go / shortcuts.go
        └── chip.go / dot.go / gap.go / tappable.go / recent.go
```

Import direction (never flip):

```
ui ─► app ─► kvstore                 # ui talks to app, app picks adapters
ui ─► kvstore (types only)           # OK — Store interface, EngineKind
ui ─► config                         # OK
```

---

## Where to change what

### Backend / database

| Change | File(s) |
|---|---|
| Add a new engine (e.g. RocksDB) | `internal/kvstore/<name>/<name>.go` (implement `Store`), `kvstore/store.go` (add `EngineKind`), `app/app.go` (dispatcher case), `ui/open_dialog.go` (`engineChoices`), `ui/theme/theme.go` (`DBAccent`) |
| Change `Store` interface | `kvstore/store.go` + update all three adapters + any caller |
| Add `OpenOptions` field (password, cache size…) | `kvstore/store.go`, adapters that need it, `ui/open_dialog.go` (`OpenRequest` + form control) |

### Visual / design

| Change | File |
|---|---|
| Any color or size token | `internal/ui/theme/theme.go` — every widget reads from here |
| Engine accent color | `theme.go` → `DBAccent` switch |
| Pill chip / dot / spacer behaviour | `chip.go` / `dot.go` / `gap.go` |
| Main-window layout (e.g. tree on right) | `mainwindow.go` — swap `container.NewHSplit(treeBox, center)` args |
| Welcome layout | `welcome.go` |

### State / behaviour

| Change | File |
|---|---|
| App-wide action (Export, Backup, …) | New method on `*AppState` in `state.go`. If shortcut: `shortcuts.go`. If toolbar/menu: `toolbar.go` + `menu.go` |
| New persisted state (log level, recent server…) | `config/config.go` field + getter/setter on `AppState` + `state.Persist()` |
| New page-level handler (e.g. Cmd+E on selection) | Add field to `pageHandlers` in `state.go`, populate it in `mainwindow.go`'s `SetPageHandlers`, expose `s.Fire…()`, wire shortcut |

### Dialogs / features

| Change | File |
|---|---|
| New dialog (Import CSV, Export, …) | `internal/ui/<feature>_dialog.go` exporting `show<Feature>(parent, sess, onDone)`. Add launcher method on `AppState` (`s.Show<Feature>()`) and wire into menu/toolbar |
| New Settings field | `config/config.go` (field) → `settings_dialog.go` (control in correct tab) → `OnChange` callback sets it on `state.cfg` and calls `state.Persist()` |
| New value-editor mode | `editor.go` — case in `rebuild()`, body builder fn, add to format `RadioGroup` |

### Filter / tree / preview

| Change | File |
|---|---|
| New filter mode (regex, value substring) | `filter.go` — `FilterState.Mode` field, switch in `applyFilter`, picker UI |
| Tree separator (e.g. `:` not `/`) | `prefixtree.go` — `buildPrefixTree` + `treeLabel`. Update `prefixtree_test.go` |
| Preview formatting | `app/session.go` — `makePreview` |

---

## Anticipated near-term refactors

### 1. Per-tab UI state

Today: switching tabs rebuilds `mainPage` from scratch, so filter text,
selected row, and scroll position are lost.

Fix: hold a `TabState` alongside each session:

```go
type TabState struct {
    Session  *app.Session
    Filter   FilterState
    Selected []byte
    ScrollY  float32
}
```

Store `[]*TabState` in `AppState` (replace `[]*app.Session`). Read from
it in `mainPage`. Persist nothing — purely in-memory.

### 2. Lazy preview for huge databases

`session.reloadKeys` currently iterates the entire store at open time to
build `KeyMeta.Preview`. For 1M+ keys this is ~200 MB of allocations.

Fix: make `Preview` a pointer or sentinel; compute on first cell render
and cache. Iterator pass only collects key + size cheaply.

### 3. Logger

`fyne.LogError` is scattered. Add `internal/log/log.go` with
`Info/Warn/Error`, write to `config.LogDir()` with daily rotation,
replace call sites.

### 4. Test coverage for AppState

The state struct is now testable. Mock `fyne.Window`/`fyne.App` and add
unit tests for `OpenSession`, `CloseAt`, `CycleTab`, `SetTheme`.

### 5. Background work

There are no goroutines today. If you add watchers, background
compactions, or auto-refresh, `AppState` needs a `context.Context` and
`cancel` so close-intercept can shut tasks down cleanly.

### 6. Iterator zero-copy API

The three adapters copy key + value on every `Entry()` call. For
streaming over millions of rows this is wasteful. Add `Key()` / `Value()`
methods that return slices valid until `Next()`, keep `Entry()` as the
safe-copy version, and change `session.reloadKeys` to use the no-copy
form.

---

## Philosophy (when in doubt)

- **Cross-cutting state** → method on `AppState`.
- **Visual concern** → `theme/theme.go` or a small widget in `internal/ui/`.
- **Data concern** → `internal/app/` (logic) or `internal/kvstore/` (storage).
- **Never** put code in `cmd/kvexplorer/main.go` — it is the entry point only.
- **Never** hardcode a color or size — always read from theme.
- **Never** let `ui/` import a concrete adapter package — it goes through
  `app/`.

Ask "whose responsibility is this?" before picking a file. If the answer
is "the session's", it goes in `session.go`. If "the app's", it's an
`AppState` method. If "this one widget's", it's a small ui file.
