package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	fynetheme "fyne.io/fyne/v2/theme"

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
		w.SetContent(welcomePage(a, w, &variant, func() {
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
