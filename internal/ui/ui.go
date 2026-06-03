package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	apptheme "github.com/mohsenm4/kv-explorer/internal/ui/theme"
)

func Run() {
	a := app.NewWithID("com.kvexplorer.app")

	variant := a.Settings().ThemeVariant()
	applyTheme(a, variant)

	w := a.NewWindow("KV-Explorer")
	w.Resize(fyne.NewSize(1280, 800))

	var render func()
	render = func() {
		w.SetContent(welcome(a, &variant, func() {
			if variant == fynetheme.VariantDark {
				variant = fynetheme.VariantLight
			} else {
				variant = fynetheme.VariantDark
			}
			applyTheme(a, variant)
			render()
		}))
	}
	render()

	w.ShowAndRun()
}

func applyTheme(a fyne.App, v fyne.ThemeVariant) {
	a.Settings().SetTheme(apptheme.ForcedVariant(apptheme.New(), v))
}

func welcome(a fyne.App, variant *fyne.ThemeVariant, onToggle func()) fyne.CanvasObject {
	th := a.Settings().Theme()
	v := *variant

	icon := widget.NewIcon(fynetheme.StorageIcon())
	iconBox := container.NewGridWrap(fyne.NewSize(64, 64), icon)

	title := canvas.NewText("KV-Explorer", th.Color(fynetheme.ColorNameForeground, v))
	title.TextSize = 22
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	tagline := canvas.NewText("Inspect, edit, and compare key-value databases.", th.Color(fynetheme.ColorNamePlaceHolder, v))
	tagline.TextSize = 14
	tagline.Alignment = fyne.TextAlignCenter

	open := widget.NewButtonWithIcon("Open Database…", fynetheme.FolderOpenIcon(), func() {
		// TODO Step 4: open database dialog
	})
	open.Importance = widget.HighImportance

	recents := fakeRecents()
	hero := container.NewVBox(
		container.NewCenter(iconBox),
		title,
		tagline,
		container.NewCenter(open),
	)
	if len(recents) > 0 {
		hero.Add(widget.NewSeparator())
		hero.Add(recentSection(v, recents, func(r recentEntry) {
			fmt.Println("recent clicked:", r.path) // wired up in Step 5
		}))
	}

	heroBox := container.NewGridWrap(fyne.NewSize(480, hero.MinSize().Height), hero)

	toggleLabel := "Dark"
	if v == fynetheme.VariantDark {
		toggleLabel = "Light"
	}
	toggle := widget.NewButton(toggleLabel, onToggle)
	toggle.Importance = widget.LowImportance

	topRow := container.NewHBox(layout.NewSpacer(), toggle)
	return container.NewBorder(topRow, nil, nil, nil, container.NewCenter(heroBox))
}

type recentEntry struct {
	path   string
	engine string // matches kvstore.EngineKind values
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

func recentSection(v fyne.ThemeVariant, entries []recentEntry, onPick func(recentEntry)) fyne.CanvasObject {
	heading := widget.NewLabelWithStyle("Recent", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	rows := container.NewVBox()
	for _, e := range entries {
		rows.Add(recentRow(v, e, onPick))
	}
	return container.NewVBox(heading, rows)
}

func recentRow(v fyne.ThemeVariant, e recentEntry, onPick func(recentEntry)) fyne.CanvasObject {
	dot := canvas.NewCircle(apptheme.DBAccent(e.engine, v))
	dot.Resize(fyne.NewSize(10, 10))
	dotBox := container.NewGridWrap(fyne.NewSize(10, 10), dot)

	pathBtn := widget.NewButton(e.path, func() { onPick(e) })
	pathBtn.Importance = widget.LowImportance
	pathBtn.Alignment = widget.ButtonAlignLeading

	when := widget.NewLabel(relTime(e.when))
	when.Importance = widget.LowImportance

	return container.NewBorder(nil, nil, dotBox, when, pathBtn)
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
