package ui

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"github.com/mohsenm4/kv-explorer/internal/config"
)

type recentEntry struct {
	path   string
	engine string
	when   time.Time
}

// recentsFromConfig converts the persisted entries into the ui-level form.
func recentsFromConfig(rs []config.Recent) []recentEntry {
	out := make([]recentEntry, 0, len(rs))
	for _, r := range rs {
		out = append(out, recentEntry{path: r.Path, engine: r.Engine, when: r.OpenedAt})
	}
	return out
}

func fakeRecents() []recentEntry {
	now := time.Now()
	return []recentEntry{
		{"~/data/users.pebble", "pebble", now.Add(-2 * time.Hour)},
		{"~/work/cache.badger", "badger", now.Add(-26 * time.Hour)},
		{"/tmp/scratch.ldb", "leveldb", now.Add(-3 * 24 * time.Hour)},
	}
}

func buildRecentBlock(v fyne.ThemeVariant, fg, muted color.Color, entries []recentEntry, onPick func(recentEntry)) fyne.CanvasObject {
	if len(entries) == 0 {
		return container.NewWithoutLayout()
	}
	heading := canvas.NewText("RECENT", muted)
	heading.TextSize = 11
	heading.TextStyle = fyne.TextStyle{Bold: true}

	rows := container.NewVBox()
	for _, e := range entries {
		rows.Add(recentRow(v, fg, muted, e, onPick))
	}
	return container.NewVBox(heading, rows)
}

func recentRow(v fyne.ThemeVariant, fg, muted color.Color, e recentEntry, onPick func(recentEntry)) fyne.CanvasObject {
	chip := engineChip(e.engine, v)
	chipBox := container.New(layout.NewCustomPaddedLayout(0, 0, 0, 12), chip)

	path := canvas.NewText(middleTruncate(e.path, 48), fg)
	path.TextSize = 13
	path.TextStyle = fyne.TextStyle{Monospace: true}

	when := canvas.NewText(relTime(e.when), muted)
	when.TextSize = 11
	whenBox := container.New(layout.NewCustomPaddedLayout(0, 0, 12, 0), when)

	row := container.NewBorder(nil, nil, chipBox, whenBox, path)
	padded := container.New(layout.NewCustomPaddedLayout(6, 6, 0, 0), row)
	return newTappable(padded, func() {
		if onPick != nil {
			onPick(e)
		}
	})
}

// middleTruncate shortens a long path by collapsing its middle to "…"
// so both the prefix (e.g. "~/Desktop") and the leaf (db name) stay
// visible. n is the target visible character count.
func middleTruncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n < 5 {
		return s[:n]
	}
	half := (n - 1) / 2
	return s[:half] + "…" + s[len(s)-(n-1-half):]
}

func relTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 48*time.Hour:
		return "yesterday"
	case d < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	default:
		return t.Format("Jan 2")
	}
}
