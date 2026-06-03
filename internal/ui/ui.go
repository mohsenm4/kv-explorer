package ui

import (
	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	fynetheme "fyne.io/fyne/v2/theme"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/config"
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

	cfg, _ := config.Load()

	systemVariant := a.Settings().ThemeVariant()
	themePref := cfg.Theme
	if themePref == "" {
		themePref = "system"
	}
	variant := variantFor(themePref, systemVariant)
	applyTheme(a, variant)

	w := a.NewWindow("KV-Explorer")
	winW, winH := float32(1280), float32(800)
	if cfg.WindowWidth > 400 {
		winW = cfg.WindowWidth
	}
	if cfg.WindowHeight > 300 {
		winH = cfg.WindowHeight
	}
	w.Resize(fyne.NewSize(winW, winH))

	var session *app.Session
	var render func()
	handlers := &appHandlers{}

	persist := func() {
		if err := config.Save(cfg); err != nil {
			fyne.LogError("config save", err)
		}
	}

	setTheme := func(pref string) {
		themePref = pref
		cfg.Theme = pref
		variant = variantFor(pref, a.Settings().ThemeVariant())
		applyTheme(a, variant)
		persist()
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
		cfg.AddRecent(req.Path, string(req.Engine))
		persist()
		render()
	}

	closeSess := func() {
		if session != nil {
			_ = session.Close()
			session = nil
		}
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
			w.SetContent(welcomePage(a, w, &variant, toggleTheme, openReq, recentsFromConfig(cfg.Recents)))
		} else {
			w.SetContent(mainPage(a, w, session, &variant, openFromMain, closeSess, toggleTheme, openSettings, handlers))
		}
	}
	render()

	w.SetMainMenu(mainMenu(w, openFromMain, closeSess, toggleTheme, openSettings))
	registerShortcuts(w, handlers)

	w.SetCloseIntercept(func() {
		size := w.Canvas().Size()
		cfg.WindowWidth = size.Width
		cfg.WindowHeight = size.Height
		persist()
		if session != nil {
			_ = session.Close()
		}
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
