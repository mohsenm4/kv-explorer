package ui

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type recentEntry struct {
	path   string
	engine string
	when   time.Time
}

func fakeRecents() []recentEntry {
	now := time.Now()
	return []recentEntry{
		{"~/data/users.pebble", "pebble", now.Add(-2 * time.Hour)},
		{"~/work/cache.badger", "badger", now.Add(-26 * time.Hour)},
		{"/tmp/scratch.ldb", "leveldb", now.Add(-3 * 24 * time.Hour)},
	}
}

func buildRecentBlock(v fyne.ThemeVariant, fg, muted color.Color, entries []recentEntry) fyne.CanvasObject {
	if len(entries) == 0 {
		return container.NewWithoutLayout()
	}
	heading := canvas.NewText("RECENT", muted)
	heading.TextSize = 11
	heading.TextStyle = fyne.TextStyle{Bold: true}

	rows := container.NewVBox()
	for _, e := range entries {
		rows.Add(recentRow(v, fg, muted, e))
	}
	return container.NewVBox(heading, rows)
}

func recentRow(v fyne.ThemeVariant, fg, muted color.Color, e recentEntry) fyne.CanvasObject {
	chip := engineChip(e.engine, v)
	chipBox := container.New(layout.NewCustomPaddedLayout(0, 0, 0, 12), chip)

	path := canvas.NewText(e.path, fg)
	path.TextSize = 13
	path.TextStyle = fyne.TextStyle{Monospace: true}

	when := canvas.NewText(relTime(e.when), muted)
	when.TextSize = 11
	whenBox := container.New(layout.NewCustomPaddedLayout(0, 0, 12, 0), when)

	row := container.NewBorder(nil, nil, chipBox, whenBox, path)
	padded := container.New(layout.NewCustomPaddedLayout(6, 6, 0, 0), row)
	return newTappable(padded, func() {
		fmt.Println("recent clicked:", e.path) // wired up in Step 5
	})
}

func relTime(t time.Time) string {
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
