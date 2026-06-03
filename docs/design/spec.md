# KV-Explorer — UI Design Specification

A complete design system for KV-Explorer's desktop UI, written so a designer
can recreate it in Figma and a developer can implement it with Fyne v2 without
guessing.

This spec covers: brand, color, typography, spacing, components, screens,
states, and accessibility. Two themes — light and dark — are described in
parallel so any token has a value in both.

> **Implementation note.** All tokens map 1-to-1 to Fyne's `theme.ColorName`
> and `theme.SizeName`. The theme layer must read from this token list —
> no colors or sizes may be hardcoded elsewhere in the UI.

---

## 1. Product Identity

- **Name**: KV-Explorer
- **Tagline**: "Inspect, edit, and compare key-value databases."
- **Tone**: Technical, calm, no-nonsense. Closer to a database tool than a
  consumer app. The UI never shouts — it informs.
- **Target user**: Backend engineers, SREs, and embedded-system developers
  who need to look inside PebbleDB / BadgerDB / LevelDB stores.

---

## 2. Color Tokens

All colors are defined per theme. The token name (left column) is what the
code and the Figma library both reference.

### 2.1 Light Theme

| Token                  | Hex       | Used for                                          |
| ---------------------- | --------- | ------------------------------------------------- |
| `background`           | `#FAFAFA` | Window background                                 |
| `surface`              | `#FFFFFF` | Cards, panels, dialog body                        |
| `surfaceSubtle`        | `#F3F4F6` | Toolbar, status bar, secondary surfaces           |
| `foreground`           | `#0F172A` | Primary text                                      |
| `foregroundMuted`      | `#64748B` | Secondary text, captions, placeholders            |
| `foregroundDisabled`   | `#CBD5E1` | Disabled text and icons                           |
| `border`               | `#E2E8F0` | Default border and separator                      |
| `borderStrong`         | `#CBD5E1` | Focused input border, table grid lines            |
| `primary`              | `#4F46E5` | Primary action buttons, links, focus ring         |
| `primaryHover`         | `#4338CA` | Primary on hover                                  |
| `primaryForeground`    | `#FFFFFF` | Text/icon on top of primary                       |
| `selection`            | `#EEF2FF` | Selected row / tree item background               |
| `selectionStrong`      | `#C7D2FE` | Selected row on dark mode contrast / focused row  |
| `success`              | `#10B981` | Positive status, connected DB indicator           |
| `warning`              | `#F59E0B` | Warning toast, unsaved-change marker              |
| `error`                | `#EF4444` | Destructive button, error border, error toast    |
| `hoverOverlay`         | `rgba(0,0,0,0.04)` | Hover overlay on any tappable region    |
| `pressedOverlay`       | `rgba(0,0,0,0.08)` | Pressed overlay on any tappable region  |
| `scrollbar`            | `#94A3B8` | Scrollbar thumb                                   |
| `shadow`               | `rgba(15,23,42,0.08)` | Dialog/menu drop shadow                |

### 2.2 Dark Theme

| Token                  | Hex                       | Used for                              |
| ---------------------- | ------------------------- | ------------------------------------- |
| `background`           | `#0B1220`                 | Window background                     |
| `surface`              | `#111827`                 | Cards, panels, dialog body            |
| `surfaceSubtle`        | `#1F2937`                 | Toolbar, status bar                   |
| `foreground`           | `#F1F5F9`                 | Primary text                          |
| `foregroundMuted`      | `#94A3B8`                 | Secondary text                        |
| `foregroundDisabled`   | `#475569`                 | Disabled text                         |
| `border`               | `#1F2937`                 | Default border                        |
| `borderStrong`         | `#334155`                 | Focused input border, grid lines      |
| `primary`              | `#818CF8`                 | Primary actions                       |
| `primaryHover`         | `#6366F1`                 | Primary on hover                      |
| `primaryForeground`    | `#0B1220`                 | Text on top of primary                |
| `selection`            | `#1E1B4B`                 | Selected row background               |
| `selectionStrong`      | `#3730A3`                 | Focused selected row                  |
| `success`              | `#34D399`                 | Positive status                       |
| `warning`              | `#FBBF24`                 | Warning                               |
| `error`                | `#F87171`                 | Destructive / error                   |
| `hoverOverlay`         | `rgba(255,255,255,0.06)`  | Hover overlay                         |
| `pressedOverlay`       | `rgba(255,255,255,0.10)`  | Pressed overlay                       |
| `scrollbar`            | `#475569`                 | Scrollbar thumb                       |
| `shadow`               | `rgba(0,0,0,0.40)`        | Dialog/menu drop shadow               |

### 2.3 Database Accent Colors

Each supported engine has its own accent. The accent appears in: tab indicator,
status badge, and the database chip on rows.

| Engine    | Token         | Light hex | Dark hex   |
| --------- | ------------- | --------- | ---------- |
| PebbleDB  | `dbPebble`    | `#0EA5E9` | `#38BDF8`  |
| BadgerDB  | `dbBadger`    | `#F59E0B` | `#FBBF24`  |
| LevelDB   | `dbLeveldb`   | `#10B981` | `#34D399`  |

---

## 3. Typography

Fyne uses the platform's default font. Figma should mirror this so the spec
matches the on-screen result.

| Token         | Family                                              | Size (pt) | Weight    | Line-height |
| ------------- | --------------------------------------------------- | --------- | --------- | ----------- |
| `displayLg`   | System default (SF Pro / Segoe UI / Cantarell)       | 22        | Semibold  | 1.30        |
| `heading`     | System default                                       | 18        | Semibold  | 1.35        |
| `subheading`  | System default                                       | 14        | Semibold  | 1.40        |
| `body`        | System default                                       | 14        | Regular   | 1.45        |
| `bodyStrong`  | System default                                       | 14        | Medium    | 1.45        |
| `caption`     | System default                                       | 11        | Regular   | 1.40        |
| `mono`        | SF Mono / JetBrains Mono / Consolas                  | 13        | Regular   | 1.50        |
| `monoSmall`   | SF Mono / JetBrains Mono / Consolas                  | 11        | Regular   | 1.45        |

- **Mono is mandatory for keys and values.** Anything that displays raw bytes
  (key path, value preview, hex view) uses `mono`. All other text uses the
  system font.

---

## 4. Spacing & Sizing

A 4-pt base grid. Every measurement is a multiple of 4. Figma should configure
its layout grid to match.

| Token       | Value | Used for                                  |
| ----------- | ----- | ----------------------------------------- |
| `space.0`   | 0pt   | Reset                                     |
| `space.1`   | 4pt   | Inner padding, icon gap                   |
| `space.2`   | 8pt   | Default control padding                   |
| `space.3`   | 12pt  | Group spacing                             |
| `space.4`   | 16pt  | Section padding                           |
| `space.6`   | 24pt  | Dialog padding                            |
| `space.8`   | 32pt  | Empty-state vertical spacing              |

| Token           | Value | Used for                                |
| --------------- | ----- | --------------------------------------- |
| `radius.sm`     | 4pt   | Buttons, inputs, chips                  |
| `radius.md`     | 6pt   | Cards, panels                           |
| `radius.lg`     | 8pt   | Dialogs                                 |
| `border.width`  | 1pt   | Default input/card border               |
| `focus.width`   | 2pt   | Focus ring                              |
| `iconSize.sm`   | 16pt  | Inline icon inside text                 |
| `iconSize.md`   | 20pt  | Toolbar icon, button leading icon       |
| `iconSize.lg`   | 32pt  | Empty-state icon                        |
| `iconSize.xl`   | 64pt  | About / welcome screen                  |

| Token             | Value     | Used for                              |
| ----------------- | --------- | ------------------------------------- |
| `control.height`  | 32pt      | Buttons, inputs, selects              |
| `row.height`      | 28pt      | Table/list row                        |
| `toolbar.height`  | 40pt      | Top toolbar                           |
| `statusbar.height`| 24pt      | Bottom status bar                     |
| `tab.height`      | 36pt      | Database tab strip                    |
| `window.minSize`  | 1024×720  | Minimum window                        |
| `window.default`  | 1280×800  | First-run window size                 |

---

## 5. Iconography

- **Source**: Fyne's bundled icon set + the same Lucide-style outlined family
  for any custom icon (so visual weight stays consistent).
- **Stroke width**: 1.5pt at 20pt icon size.
- **Color**: Always `foreground` or `foregroundMuted`; never colored except
  for status (success, warning, error) and DB accent chips.
- **Required icons** (in the Figma component library):
  - File: `folder-open`, `file-plus`, `save`, `download`, `upload`
  - Edit: `plus`, `trash`, `pencil`, `copy`, `clipboard`
  - View: `search`, `filter`, `refresh`, `eye`, `eye-off`
  - DB: `database`, `database-plus`, `database-x`
  - State: `check`, `x`, `info`, `alert-triangle`, `alert-circle`
  - Nav: `chevron-right`, `chevron-down`, `chevron-up`, `chevron-left`
  - Misc: `settings`, `sun`, `moon`, `keyboard`, `help-circle`

---

## 6. Component Library

Every component below maps to a Fyne widget (or a small composition of them).
Figma should publish each as a Component with variants for state.

### 6.1 Button

Fyne primitive: `widget.Button` with `Importance`.

| Variant       | Background        | Foreground          | Border         | Use                          |
| ------------- | ----------------- | ------------------- | -------------- | ---------------------------- |
| Primary       | `primary`         | `primaryForeground` | none           | Confirm action in a dialog   |
| Secondary     | `surface`         | `foreground`        | `border`       | Cancel, secondary action     |
| Ghost         | transparent       | `foreground`        | none           | Toolbar text button          |
| Destructive   | `error`           | `#FFFFFF`           | none           | Delete confirm               |
| Icon-only     | transparent       | `foreground`        | none           | Toolbar icon, row actions    |

States: `default`, `hover` (apply `hoverOverlay`), `pressed`
(apply `pressedOverlay`), `focused` (2pt `primary` ring offset 2pt), `disabled`
(opacity 50%, no hover).

Sizing: height `control.height` (32pt), horizontal padding `space.3` (12pt),
gap between icon and label `space.1` (4pt).

### 6.2 Input (Entry)

Fyne primitive: `widget.Entry`, `widget.PasswordEntry`, `widget.MultiLineEntry`.

- Background: `surface`
- Border: 1pt `border`, becomes `borderStrong` on hover, `primary` on focus
- Padding: `space.2` (8pt) horizontal, `space.1` (4pt) vertical
- Placeholder color: `foregroundMuted`
- Invalid state: 1pt `error` border, helper text in `error` below the field
- Helper/error text: `caption`
- Multi-line: minimum 3 lines, max grows to fill, monospace font (since it
  typically holds JSON or raw values)

### 6.3 Select / Dropdown

Fyne primitive: `widget.Select`.

- Same dimensions and border behavior as Input.
- Trailing `chevron-down` icon at `iconSize.sm`, padded `space.2` from text.
- Pop-up panel: `surface` with `radius.md`, `shadow`, max-height ~280pt with
  internal scroll.

### 6.4 Checkbox & RadioGroup

Fyne primitives: `widget.Check`, `widget.RadioGroup`.

- Box / circle: 16×16, 1pt `border`, fills with `primary` when checked.
- Tick / dot color: `primaryForeground`.
- Label `body`, gap `space.2`.

### 6.5 Tabs

Fyne primitive: `container.AppTabs`.

- Strip height: `tab.height` (36pt).
- Each tab: padding `space.3` horizontal, gap `space.2` between icon and label.
- Inactive: `foregroundMuted` text, no underline.
- Active: `foreground` text, 2pt underline in the DB accent color
  (so the user sees at a glance which engine is open).
- Hover: `hoverOverlay`.

### 6.6 Toolbar

Fyne primitive: `widget.Toolbar`.

- Height: `toolbar.height` (40pt).
- Background: `surfaceSubtle`.
- Bottom border: 1pt `border`.
- Items: icon buttons at `iconSize.md`, padded `space.2`.
- Separator: 1pt vertical line in `border`, height `iconSize.md`.
- Spacer: flexible, pushes following items right.

### 6.7 Status Bar

Custom composition (HBox with labels).

- Height: `statusbar.height` (24pt).
- Background: `surfaceSubtle`.
- Top border: 1pt `border`.
- Sections, left to right:
  1. DB engine name + colored dot (engine accent)
  2. Key count (e.g. "1,234 keys")
  3. On-disk size (e.g. "4.2 MB")
  4. Spacer
  5. Open path (truncated middle, full path on hover tooltip)
  6. Theme toggle icon button

### 6.8 Tree (Key Prefix Tree)

Fyne primitive: `widget.Tree`.

- Row height: `row.height` (28pt).
- Indent per level: `space.4` (16pt).
- Caret icon: `chevron-right` (collapsed) / `chevron-down` (expanded), 16pt,
  `foregroundMuted`.
- Folder icon (intermediate node): `folder` outline at `iconSize.sm`,
  `foregroundMuted`.
- Leaf icon (actual key): `file` outline at `iconSize.sm`, `foregroundMuted`.
- Selected row: background `selection`, text `foreground`.
- Focused selected row: background `selectionStrong`.
- Hover: `hoverOverlay`.

### 6.9 Table (Key List)

Fyne primitive: `widget.Table`.

- Header height: 32pt, background `surfaceSubtle`, text `subheading`,
  bottom border 1pt `border`.
- Row height: `row.height` (28pt).
- Zebra striping: off by default; can be enabled in Settings.
- Cell padding: `space.2` horizontal.
- Default columns:
  | Column         | Width    | Notes                              |
  | -------------- | -------- | ---------------------------------- |
  | Key            | flex 60% | `mono`, truncate-with-ellipsis     |
  | Value preview  | flex 35% | `mono`, single line                |
  | Size           | 80pt     | right-aligned, `caption`           |
- Right-click context menu: Edit, Copy key, Copy value, Delete.

### 6.10 Split Pane

Fyne primitive: `container.Split`.

- Divider: 4pt wide, color `border`, becomes `primary` on hover.
- Drag cursor on hover.

### 6.11 Card

Custom composition (border + padding).

- Background: `surface`
- Border: 1pt `border`
- Radius: `radius.md` (6pt)
- Padding: `space.4` (16pt)
- Optional header row (HBox at top): title in `subheading` + trailing actions.

### 6.12 Dialog (Modal)

Fyne primitive: `dialog.Custom` and friends.

- Backdrop: `rgba(0,0,0,0.40)` (both themes).
- Panel: `surface`, `radius.lg` (8pt), `shadow`.
- Width: 480pt for forms, 640pt for confirms with key/value preview,
  720pt for the "Open Database" picker.
- Padding: `space.6` (24pt).
- Title: `heading`, bottom margin `space.4`.
- Footer: right-aligned actions, gap `space.2`. Primary action on the right,
  Cancel to its left.
- Close icon (`x` at `iconSize.sm`) in the top-right of the title row.

### 6.13 Toast / Snackbar

Custom (transient overlay).

- Position: bottom-center, 24pt above status bar.
- Width: auto, max 480pt.
- Padding: `space.2` `space.4`.
- Background: `foreground` (inverted color for contrast).
- Text: `surface` color on top.
- Variants: success (left icon `check`, accent `success`), warning, error.
- Auto-dismiss: 3.5s; manual dismiss via `x` button.

### 6.14 Database Chip

Custom (small badge used in row metadata and dialogs).

- Pill shape: radius 999pt, padding `space.1` `space.2`.
- Background: engine accent at 20% opacity.
- Foreground: engine accent at full opacity.
- Text: `caption`, weight Medium, engine name in uppercase ("PEBBLE").

### 6.15 Progress

Fyne primitives: `widget.ProgressBar`, `widget.Activity`.

- ProgressBar: 4pt height, track `border`, fill `primary`.
- Activity (spinner): 24pt, color `primary`, 8 dots / segmented arc.

### 6.16 Search Field

Composition: Entry + leading `search` icon + trailing `x` clear button.

- Leading icon `search` `iconSize.sm`, color `foregroundMuted`, padded
  `space.2` from text.
- Trailing clear button visible only when not empty.
- Submits on Enter; live-filters on each change after 250ms debounce.

---

## 7. Screen Wireframes

### 7.1 Welcome (No Database Open)

```text
┌──────────────────────────────────────────────────────────────┐
│ MainMenu: File   Edit   View   Help                          │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│                                                              │
│                        ┌───────────┐                         │
│                        │  database │   (icon, 64pt)          │
│                        └───────────┘                         │
│                                                              │
│                       KV-Explorer                            │
│              Inspect, edit, and compare KV stores.           │
│                                                              │
│            [ Open Database… ]   [ Open Recent ▾ ]            │
│                                                              │
│              Recent:                                         │
│                · ~/data/users.pebble        2h ago           │
│                · ~/work/cache.badger        yesterday        │
│                · /tmp/scratch.ldb           3d ago           │
│                                                              │
│                                                              │
├──────────────────────────────────────────────────────────────┤
│ No database open                                  ☀  Theme  │
└──────────────────────────────────────────────────────────────┘
```

- Centered hero with logo, name, tagline, and primary CTA.
- Recent list shows up to 5 entries, each clickable to reopen.
- Status bar shows "No database open"; theme toggle stays visible.

### 7.2 Main Window (Database Open)

```text
┌──────────────────────────────────────────────────────────────────────────┐
│ MainMenu: File   Edit   View   Help                                      │
├──────────────────────────────────────────────────────────────────────────┤
│ [📁 Open] [✕ Close] | [+ Add] [✎ Edit] [🗑 Delete] | [↻ Refresh] [⚙]    │  (toolbar 40)
├──────────────────────────────────────────────────────────────────────────┤
│ ▎PebbleDB ▎ ▏BadgerDB ▏ ▏LevelDB ▏ [+]                                   │  (tab strip 36)
├──────────────────────────────────────────────────────────────────────────┤
│              │  🔍 [ Filter keys…                          ]   [ × ]     │
│  Prefixes    ├──────────────────────────────────────────────────────────┤
│              │  Key                       Value preview              Size │
│  ▾ users/    │  users/0001                {"name":"ali","age":30}     42B│
│    · 0001    │  users/0002                {"name":"sara","age":24}    44B│ ◄ selected
│    · 0002    │  users/0003                {"name":"reza","age":40}    43B│
│  ▸ sessions/ │  users/0004                {"name":"mina","age":27}    45B│
│  ▸ logs/     │  ...                                                       │
│  ▸ cache/    │                                                            │
│              ├──────────────────────────────────────────────────────────┤
│              │  Value: users/0002                          [ ⌄ ] [ ⤢ ]  │
│              │  ┌────────────────────────────────────────────────────┐  │
│              │  │ {                                                  │  │
│              │  │   "name": "sara",                                  │  │
│              │  │   "age":  24,                                      │  │
│              │  │   "tags": ["admin", "beta"]                        │  │
│              │  │ }                                                  │  │
│              │  └────────────────────────────────────────────────────┘  │
│              │  [ Cancel ]                              [ Save changes ]│
├──────────────┴──────────────────────────────────────────────────────────┤
│ ● PEBBLE   1,234 keys   4.2 MB   ~/data/users.pebble              ☀  ⚙ │  (status 24)
└──────────────────────────────────────────────────────────────────────────┘
```

- **Top toolbar** — primary actions, grouped by separator.
- **Tab strip** — one tab per open database; "+" opens another.
- **Left pane** (Tree) — keys grouped by `/` prefix; clicking a leaf
  selects the matching row in the table.
- **Right pane** (Table over Editor, vertical split) — filter on top,
  table in the middle, value editor at the bottom (collapsible).
- **Status bar** — engine dot + counts + path + theme toggle.

### 7.3 Open Database Dialog

```text
                    ┌──────────────────────────────────────┐
                    │  Open Database                    ×  │
                    ├──────────────────────────────────────┤
                    │  Engine                              │
                    │  ( ) PebbleDB                        │
                    │  (•) BadgerDB                        │
                    │  ( ) LevelDB                         │
                    │                                      │
                    │  Path                                │
                    │  [ /Users/.../cache.badger    ] [📁] │
                    │                                      │
                    │  ☑ Open in new tab                   │
                    │  ☐ Read-only                         │
                    │                                      │
                    ├──────────────────────────────────────┤
                    │                  [ Cancel ]  [ Open ]│
                    └──────────────────────────────────────┘
```

### 7.4 Add / Edit Key Dialog

```text
            ┌────────────────────────────────────────────────────┐
            │  Edit key                                       ×  │
            ├────────────────────────────────────────────────────┤
            │  Key                                               │
            │  [ users/0002                                    ] │
            │  ⚠ Key already exists                              │
            │                                                    │
            │  Value (UTF-8)                                     │
            │  ┌──────────────────────────────────────────────┐  │
            │  │ {                                            │  │
            │  │   "name": "sara",                            │  │
            │  │ }                                            │  │
            │  └──────────────────────────────────────────────┘  │
            │  44 B                                              │
            │                                                    │
            │  Format: ( ) Hex  (•) UTF-8  ( ) JSON              │
            │                                                    │
            ├────────────────────────────────────────────────────┤
            │                       [ Cancel ]  [ Save changes ] │
            └────────────────────────────────────────────────────┘
```

### 7.5 Delete Confirmation

```text
                    ┌──────────────────────────────────────┐
                    │  Delete key?                      ×  │
                    ├──────────────────────────────────────┤
                    │  This will permanently delete:       │
                    │                                      │
                    │  users/0002                          │
                    │                                      │
                    │  The action cannot be undone.        │
                    │                                      │
                    ├──────────────────────────────────────┤
                    │            [ Cancel ]  [ 🗑 Delete ]  │
                    └──────────────────────────────────────┘
```

Primary button uses the Destructive variant.

### 7.6 Settings Dialog

```text
        ┌──────────────────────────────────────────────────────┐
        │  Settings                                         ×  │
        ├────────────────┬─────────────────────────────────────┤
        │ ▸ Appearance   │  Theme                              │
        │   General      │   ( ) Light                         │
        │   Editor       │   ( ) Dark                          │
        │   Shortcuts    │   (•) Follow system                 │
        │   About        │                                     │
        │                │  Density                            │
        │                │   ( ) Compact                       │
        │                │   (•) Comfortable                   │
        │                │                                     │
        │                │  ☑ Show zebra rows                  │
        │                │  ☐ Use monospace everywhere         │
        │                │                                     │
        ├────────────────┴─────────────────────────────────────┤
        │                                  [ Close ]  [ Save ] │
        └──────────────────────────────────────────────────────┘
```

Left side: navigation list. Right side: settings group. The same pattern
applies for "General" (paths, log level), "Editor" (font size, JSON
pretty-print on save), "Shortcuts" (key-binding list), "About" (version,
license, links).

### 7.7 Empty State (No Keys)

```text
┌──────────────────────────────────────────────────────────────┐
│  🔍 [ Filter keys…                                ]          │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│                      ┌──────────┐                            │
│                      │   file   │   (32pt muted icon)        │
│                      └──────────┘                            │
│                                                              │
│                   No keys in this store                      │
│              Add a key to get started.                       │
│                                                              │
│                     [ + Add key ]                            │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### 7.8 Error State (Open Failed)

```text
                    ┌──────────────────────────────────────┐
                    │  ⚠ Couldn't open database         ×  │
                    ├──────────────────────────────────────┤
                    │  Path: /tmp/scratch.ldb              │
                    │                                      │
                    │  leveldb: file does not exist        │
                    │                                      │
                    │  [▾ Show details]                    │
                    │                                      │
                    ├──────────────────────────────────────┤
                    │                  [ Close ]  [ Retry ]│
                    └──────────────────────────────────────┘
```

---

## 8. State Catalog

For every interactive component, Figma must publish these variants:

| State        | Trigger              | Visual change                                    |
| ------------ | -------------------- | ------------------------------------------------ |
| `default`    | resting              | base tokens                                      |
| `hover`      | pointer over         | `hoverOverlay` on top of background              |
| `pressed`    | pointer down         | `pressedOverlay` on top of background            |
| `focused`    | keyboard focus       | 2pt `primary` ring, offset 2pt                   |
| `disabled`   | non-interactive      | opacity 50%, no hover/press                      |
| `selected`   | selection (list/row) | background `selection` (or `selectionStrong` if focused) |
| `error`      | validation failed    | 1pt `error` border, helper text in `error`       |
| `loading`    | async pending        | spinner replaces leading icon; control disabled  |

---

## 9. Keyboard Shortcuts

These show up in the menu and in the Settings ▸ Shortcuts panel. The UI design
must visibly support them (every action reachable by keyboard).

| Shortcut          | Action                                        |
| ----------------- | --------------------------------------------- |
| `Ctrl/Cmd + O`    | Open Database…                                |
| `Ctrl/Cmd + W`    | Close current tab                             |
| `Ctrl/Cmd + N`    | Add key                                       |
| `Ctrl/Cmd + F`    | Focus filter                                  |
| `Ctrl/Cmd + S`    | Save value edits                              |
| `Delete`          | Delete selected key (with confirm)            |
| `F2`              | Edit selected key                             |
| `F5`              | Refresh                                       |
| `Ctrl/Cmd + ,`    | Open Settings                                 |
| `Ctrl/Cmd + Tab`  | Cycle database tabs                           |
| `Esc`             | Close dialog / clear filter / cancel edit     |

---

## 10. Accessibility

- Minimum touch target: 32×32. Toolbar icon buttons are 32×32 with the icon
  centered at 20pt.
- Color contrast: every foreground/background pair in §2 meets WCAG AA
  (4.5:1) for body text, 3:1 for large text and UI components.
- Focus is always visible — never `outline: none` without an equivalent
  ring. The focus ring is the 2pt `primary` color.
- Every icon-only button has a tooltip and an `aria-label` equivalent.
- Tree, Table, and Tabs are keyboard-navigable (arrow keys, Home/End,
  Page Up/Down).
- Don't rely on color alone — the engine accent always pairs with the engine
  name; the success/error toast always pairs with an icon.

---

## 11. Figma File Structure (recommended)

Set the Figma file up like this so handoff stays clean:

```text
KV-Explorer
├── 🎨 Tokens
│   ├── Color — Light
│   ├── Color — Dark
│   ├── Typography
│   ├── Spacing
│   └── Radius & Shadow
├── 🧩 Components
│   ├── Button (5 variants × 6 states)
│   ├── Input (3 variants × 5 states)
│   ├── Select, Check, Radio
│   ├── Tabs, Toolbar, StatusBar
│   ├── Tree row, Table row, Header
│   ├── Card, Dialog, Toast
│   ├── DB Chip (3 engines)
│   ├── Icon (set described in §5)
│   └── Progress, Activity
├── 📱 Screens — Light
│   ├── Welcome
│   ├── Main (DB open)
│   ├── Empty store
│   ├── Open dialog
│   ├── Add/Edit dialog
│   ├── Delete confirm
│   ├── Settings (4 tabs)
│   └── Error dialog
├── 📱 Screens — Dark
│   └── (mirror of Light)
└── 📐 Specs
    └── Annotated copies of the main screen with sizes
```

- Every screen must exist in both themes.
- Components use Figma Variables tied to the Color tokens so the theme switch
  is a single dropdown.

---

## 12. What Fyne Constrains (so Figma doesn't go off-spec)

Fyne is canvas-based — the designer should respect what Fyne can render
cheaply, otherwise the implementation diverges from the mockup.

- **No drop shadows on inline elements.** Shadow is only on dialogs/menus.
- **No background gradients on buttons/inputs.** Solid colors only.
- **Radius is global per widget kind**, not per corner. Don't design only-top
  rounded corners.
- **No custom hover animations longer than 120ms.** Fyne supports fades but
  not spring physics.
- **Icons are monochrome** by default; multi-color icons are possible but
  require custom `CanvasObject` work — avoid unless necessary.
- **System dialogs (native file picker)** look like the OS, not the app
  theme. The "Open path" picker therefore won't perfectly match the Figma
  mockup on every platform — design the entry point but accept the native
  picker for the actual file chooser.

---

## 13. Deliverables Checklist

Before handoff, confirm the Figma file has:

- [ ] Color, typography, spacing, radius tokens published as Variables.
- [ ] Each component in §6 published with all variants from §8.
- [ ] All 8 screens from §7 in both light and dark.
- [ ] Annotated version of the main screen with sizes referencing tokens
      (e.g. "Toolbar height: `toolbar.height` = 40pt").
- [ ] Icon set from §5 imported as a Figma component.
- [ ] Three database accent colors applied consistently in tabs, chips,
      and status dot.
