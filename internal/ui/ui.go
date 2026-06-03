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

func Run() {
	a := fyneapp.NewWithID("com.kvexplorer.app")

	variant := a.Settings().ThemeVariant()
	applyTheme(a, variant)

	w := a.NewWindow("KV-Explorer")
	w.Resize(fyne.NewSize(1280, 800))

	var session *app.Session
	var render func()

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
		render()
	}

	toggleTheme := func() {
		if variant == fynetheme.VariantDark {
			variant = fynetheme.VariantLight
		} else {
			variant = fynetheme.VariantDark
		}
		applyTheme(a, variant)
		render()
	}

	openFromMain := func() {
		showOpenDatabase(w, openReq)
	}

	render = func() {
		if session == nil {
			w.SetContent(welcomePage(a, w, &variant, toggleTheme, openReq))
		} else {
			w.SetContent(mainPage(a, w, session, &variant, openFromMain, closeSess, toggleTheme))
		}
	}
	render()

	w.SetMainMenu(mainMenu(w, openFromMain, closeSess, toggleTheme))

	w.ShowAndRun()
}

func applyTheme(a fyne.App, v fyne.ThemeVariant) {
	a.Settings().SetTheme(apptheme.ForcedVariant(apptheme.New(), v))
}
