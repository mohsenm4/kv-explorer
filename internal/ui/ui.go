package ui

import (
	"fmt"
	"image/color"
	"strings"
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
		w.SetContent(welcomePage(a, &variant, func() {
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

func welcomePage(a fyne.App, variant *fyne.ThemeVariant, onToggle func()) fyne.CanvasObject {
	th := a.Settings().Theme()
	v := *variant

	fg := th.Color(fynetheme.ColorNameForeground, v)
	muted := th.Color(fynetheme.ColorNamePlaceHolder, v)
	primary := th.Color(fynetheme.ColorNamePrimary, v)

	hero := heroIcon(primary)

	title := canvas.NewText("KV-Explorer", fg)
	title.TextSize = 22
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	tagline := canvas.NewText("Inspect, edit, and compare key-value databases.", muted)
	tagline.TextSize = 14
	tagline.Alignment = fyne.TextAlignCenter

	open := widget.NewButtonWithIcon("Open Database…", fynetheme.FolderOpenIcon(), func() {
		// TODO Step 4: open database dialog
	})
	open.Importance = widget.HighImportance

	openRecent := widget.NewButtonWithIcon("Open Recent", fynetheme.MenuDropDownIcon(), func() {
		// TODO Step 4: recent dropdown
	})
	openRecent.IconPlacement = widget.ButtonIconTrailingText

	actions := container.NewHBox(layout.NewSpacer(), open, openRecent, layout.NewSpacer())

	recentBlock := buildRecentBlock(v, fg, muted, fakeRecents())

	heroStack := container.NewVBox(
		container.NewCenter(hero),
		title,
		tagline,
		container.NewPadded(actions),
		container.NewPadded(widget.NewSeparator()),
		recentBlock,
	)

	heroSized := container.NewGridWrap(fyne.NewSize(520, heroStack.MinSize().Height), heroStack)
	center := container.NewCenter(heroSized)

	return container.NewBorder(nil, statusBar(v, "No database open", onToggle), nil, nil, center)
}

func heroIcon(primary color.Color) fyne.CanvasObject {
	bg := canvas.NewRectangle(primary)
	bg.CornerRadius = 14

	sym := canvas.NewText("⌘", color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF})
	sym.TextSize = 32
	sym.TextStyle = fyne.TextStyle{Bold: true}
	sym.Alignment = fyne.TextAlignCenter

	stack := container.NewStack(bg, container.NewCenter(sym))
	return container.NewGridWrap(fyne.NewSize(64, 64), stack)
}

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

func engineChip(engine string, v fyne.ThemeVariant) fyne.CanvasObject {
	bg := canvas.NewRectangle(apptheme.DBAccentTint(engine, v))
	bg.CornerRadius = 999

	label := canvas.NewText(strings.ToUpper(engine), apptheme.DBAccent(engine, v))
	label.TextSize = 10
	label.TextStyle = fyne.TextStyle{Bold: true}
	label.Alignment = fyne.TextAlignCenter

	padded := container.New(layout.NewCustomPaddedLayout(2, 2, 10, 10), label)
	stack := container.NewStack(bg, padded)
	return container.NewGridWrap(fyne.NewSize(72, 22), stack)
}

func statusBar(v fyne.ThemeVariant, message string, onToggle func()) fyne.CanvasObject {
	th := fyne.CurrentApp().Settings().Theme()
	bgCol := th.Color(fynetheme.ColorNameHeaderBackground, v)
	muted := th.Color(fynetheme.ColorNamePlaceHolder, v)

	bg := canvas.NewRectangle(bgCol)

	msg := canvas.NewText(message, muted)
	msg.TextSize = 11

	toggleLabel := "Light"
	if v == fynetheme.VariantDark {
		toggleLabel = "Dark"
	}
	toggle := widget.NewButton(toggleLabel, onToggle)
	toggle.Importance = widget.MediumImportance

	row := container.NewBorder(nil, nil, container.NewPadded(msg), container.NewPadded(toggle), layout.NewSpacer())

	return container.NewStack(bg, row)
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

// tappable wraps any canvas object and makes it tap-aware.
type tappable struct {
	widget.BaseWidget
	content fyne.CanvasObject
	onTap   func()
}

func newTappable(content fyne.CanvasObject, onTap func()) *tappable {
	t := &tappable{content: content, onTap: onTap}
	t.ExtendBaseWidget(t)
	return t
}

func (t *tappable) Tapped(*fyne.PointEvent) {
	if t.onTap != nil {
		t.onTap()
	}
}

func (t *tappable) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.content)
}
