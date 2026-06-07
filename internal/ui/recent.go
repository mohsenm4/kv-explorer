package ui

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/config"
	"github.com/mohsenm4/kv-explorer/internal/i18n"
)

type recentEntry struct {
	path   string
	engine string
	when   time.Time
}

func recentsFromConfig(rs []config.Recent) []recentEntry {
	out := make([]recentEntry, 0, len(rs))
	for _, r := range rs {
		out = append(out, recentEntry{path: r.Path, engine: r.Engine, when: r.OpenedAt})
	}
	return out
}

func showRecentMenu(parent fyne.Window, anchor fyne.CanvasObject, entries []recentEntry, onPick func(recentEntry)) {
	if len(entries) == 0 {
		return
	}
	items := make([]*fyne.MenuItem, 0, len(entries))
	for _, e := range entries {
		entry := e
		label := fmt.Sprintf("[%s]  %s", entry.engine, middleTruncate(entry.path, 56))
		items = append(items, fyne.NewMenuItem(label, func() {
			if onPick != nil {
				onPick(entry)
			}
		}))
	}
	menu := fyne.NewMenu("", items...)
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(anchor)
	pos.Y += anchor.Size().Height
	widget.ShowPopUpMenuAtPosition(menu, parent.Canvas(), pos)
}

func buildRecentBlock(v fyne.ThemeVariant, fg, muted color.Color, entries []recentEntry, onPick func(recentEntry)) fyne.CanvasObject {
	if len(entries) == 0 {
		return container.NewWithoutLayout()
	}
	heading := canvas.NewText(i18n.T("welcome.recent"), muted)
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

// middleTruncate collapses the middle of a path to "…" so both prefix and leaf stay visible; n is the target visible length.
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
		return i18n.T("time.justNow")
	case d < time.Hour:
		return i18n.Tf("time.minutesAgo", map[string]any{"Count": int(d.Minutes())})
	case d < 24*time.Hour:
		return i18n.Tf("time.hoursAgo", map[string]any{"Count": int(d.Hours())})
	case d < 48*time.Hour:
		return i18n.T("time.yesterday")
	case d < 7*24*time.Hour:
		return i18n.Tf("time.daysAgo", map[string]any{"Count": int(d.Hours() / 24)})
	default:
		return t.Format("Jan 2")
	}
}
