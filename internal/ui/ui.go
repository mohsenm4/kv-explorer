package ui

import (
	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	fynetheme "fyne.io/fyne/v2/theme"

	apptheme "github.com/mohsenm4/kv-explorer/internal/ui/theme"
)

func Run(version string) {
	a := fyneapp.NewWithID("com.kvexplorer.app")
	w := a.NewWindow("KV-Explorer")

	state := NewAppState(a, w)
	state.SetVersion(version)
	state.ApplyTheme()
	state.ApplyInitialWindowSize()

	var render func()
	render = func() {
		if state.Active() == nil {
			w.SetContent(welcomePage(state))
			return
		}
		w.SetContent(mainPage(state))
	}
	state.SetNotify(render)
	render()

	w.SetMainMenu(mainMenu(state))
	registerShortcuts(w, state)

	w.SetCloseIntercept(func() {
		state.SaveWindowSize()
		state.CloseAll()
		w.Close()
	})

	w.ShowAndRun()
}

func applyTheme(a fyne.App, v fyne.ThemeVariant) {
	a.Settings().SetTheme(apptheme.ForcedVariant(apptheme.New(), v))
}

func variantFor(pref string, system fyne.ThemeVariant) fyne.ThemeVariant {
	switch pref {
	case "light":
		return fynetheme.VariantLight
	case "dark":
		return fynetheme.VariantDark
	default:
		return system
	}
}
