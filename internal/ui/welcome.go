package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func welcomePage(a fyne.App, w fyne.Window, variant *fyne.ThemeVariant, onToggle func()) fyne.CanvasObject {
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
		showOpenDatabase(w, func(req OpenRequest) {
			fmt.Printf("open requested: %+v\n", req) // wired up in Step 5
		})
	})
	open.Importance = widget.HighImportance

	openRecent := widget.NewButtonWithIcon("Open Recent", fynetheme.MenuDropDownIcon(), func() {
		// TODO Step 4: recent dropdown menu
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
