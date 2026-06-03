package ui

import (
	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	fynetheme "fyne.io/fyne/v2/theme"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/kvstore"
	apptheme "github.com/mohsenm4/kv-explorer/internal/ui/theme"
)

// appHandlers is the set of cross-cutting actions that ui-level shortcuts
// fire. mainPage fills it on construction; ui.go's shortcut wiring uses
// the latest values via the pointer.
type appHandlers struct {
	addKey      func()
	editKey     func()
	deleteKey   func()
	refresh     func()
	focusFilter func()
}

func Run() {
	a := fyneapp.NewWithID("com.kvexplorer.app")

	systemVariant := a.Settings().ThemeVariant()
	themePref := "system"
	variant := systemVariant
	applyTheme(a, variant)

	w := a.NewWindow("KV-Explorer")
	w.Resize(fyne.NewSize(1280, 800))

	var session *app.Session
	var render func()
	handlers := &appHandlers{}

	setTheme := func(pref string) {
		themePref = pref
		switch pref {
		case "light":
			variant = fynetheme.VariantLight
		case "dark":
			variant = fynetheme.VariantDark
		default:
			variant = a.Settings().ThemeVariant()
		}
		applyTheme(a, variant)
		render()
	}

	openReq := func(req OpenRequest) {
		sess, err := app.OpenSession(req.Engine, req.Path, kvstore.OpenOptions{ReadOnly: req.ReadOnly})
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if session != nil {
			_ = session.Close()
		}
		session = sess
		render()
	}

	closeSess := func() {
		if session != nil {
			_ = session.Close()
			session = nil
		}
		// Clear cross-cutting handlers so shortcuts no-op on welcome.
		*handlers = appHandlers{}
		render()
	}

	toggleTheme := func() {
		if variant == fynetheme.VariantDark {
			setTheme("light")
		} else {
			setTheme("dark")
		}
	}

	openFromMain := func() {
		showOpenDatabase(w, openReq)
	}

	openSettings := func() {
		showSettings(w, themePref, SettingsHandlers{OnTheme: setTheme})
	}

	render = func() {
		if session == nil {
			*handlers = appHandlers{}
			w.SetContent(welcomePage(a, w, &variant, toggleTheme, openReq))
		} else {
			w.SetContent(mainPage(a, w, session, &variant, openFromMain, closeSess, toggleTheme, openSettings, handlers))
		}
	}
	render()

	w.SetMainMenu(mainMenu(w, openFromMain, closeSess, toggleTheme, openSettings))
	registerShortcuts(w, handlers)

	w.ShowAndRun()
}

func applyTheme(a fyne.App, v fyne.ThemeVariant) {
	a.Settings().SetTheme(apptheme.ForcedVariant(apptheme.New(), v))
}
