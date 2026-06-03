package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	apptheme "github.com/mohsenm4/kv-explorer/internal/ui/theme"
)

func Run() {
	a := app.NewWithID("com.kvexplorer.app")
	a.Settings().SetTheme(apptheme.New())

	w := a.NewWindow("KV-Explorer")
	w.Resize(fyne.NewSize(1280, 800))
	w.SetContent(welcome(a, w))
	w.ShowAndRun()
}

func welcome(a fyne.App, w fyne.Window) fyne.CanvasObject {
	icon := widget.NewIcon(fynetheme.StorageIcon())
	iconBox := container.NewGridWrap(fyne.NewSize(64, 64), icon)

	th := a.Settings().Theme()
	v := a.Settings().ThemeVariant()

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

	hero := container.NewVBox(
		container.NewCenter(iconBox),
		title,
		tagline,
		container.NewCenter(open),
	)

	return container.NewCenter(hero)
}
