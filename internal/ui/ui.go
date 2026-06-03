package ui

import (
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

	hero := container.NewVBox(
		container.NewCenter(iconBox),
		title,
		tagline,
		container.NewCenter(open),
	)

	toggleLabel := "Dark"
	if v == fynetheme.VariantDark {
		toggleLabel = "Light"
	}
	toggle := widget.NewButton(toggleLabel, onToggle)
	toggle.Importance = widget.LowImportance

	topRow := container.NewHBox(layout.NewSpacer(), toggle)
	return container.NewBorder(topRow, nil, nil, nil, container.NewCenter(hero))
}
