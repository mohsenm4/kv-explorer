package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	fynetheme "fyne.io/fyne/v2/theme"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/config"
	"github.com/mohsenm4/kv-explorer/internal/i18n"
	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

type AppState struct {
	a fyne.App
	w fyne.Window

	cfg       config.Config
	sessions  []*app.Session
	active    int
	themePref string
	variant   fyne.ThemeVariant
	version   string

	page   pageHandlers
	notify func()
}

type pageHandlers struct {
	addKey      func()
	editKey     func()
	deleteKey   func()
	refresh     func()
	focusFilter func()
}

func NewAppState(a fyne.App, w fyne.Window) *AppState {
	cfg, _ := config.Load()
	pref := cfg.Theme
	if pref == "" {
		pref = "system"
	}
	return &AppState{
		a:         a,
		w:         w,
		cfg:       cfg,
		themePref: pref,
		variant:   variantFor(pref, a.Settings().ThemeVariant()),
	}
}

func (s *AppState) SetNotify(fn func()) { s.notify = fn }

func (s *AppState) SetVersion(v string) { s.version = v }
func (s *AppState) Version() string     { return s.version }

func (s *AppState) Notify() {
	if s.notify != nil {
		s.notify()
	}
}

func (s *AppState) Variant() fyne.ThemeVariant { return s.variant }
func (s *AppState) ThemePref() string          { return s.themePref }

func (s *AppState) ApplyTheme() {
	applyTheme(s.a, s.variant)
}

func (s *AppState) SetTheme(pref string) {
	s.themePref = pref
	s.cfg.Theme = pref
	s.variant = variantFor(pref, s.a.Settings().ThemeVariant())
	s.ApplyTheme()
	s.Persist()
	s.Notify()
}

func (s *AppState) ToggleTheme() {
	if s.variant == fynetheme.VariantDark {
		s.SetTheme("light")
	} else {
		s.SetTheme("dark")
	}
}

func (s *AppState) Sessions() []*app.Session { return s.sessions }
func (s *AppState) ActiveIdx() int           { return s.active }

func (s *AppState) Active() *app.Session {
	if s.active < 0 || s.active >= len(s.sessions) {
		return nil
	}
	return s.sessions[s.active]
}

func (s *AppState) OpenSession(req OpenRequest) {
	var sess *app.Session
	withProgress(s.w, i18n.T("progress.opening"), func() error {
		var err error
		sess, err = app.OpenSession(req.Engine, req.Path, kvstore.OpenOptions{ReadOnly: req.ReadOnly})
		return err
	}, func(err error) {
		if err != nil {
			dialog.ShowError(err, s.w)
			return
		}
		if !req.NewTab && len(s.sessions) > 0 {
			_ = s.sessions[s.active].Close()
			s.sessions[s.active] = sess
		} else {
			s.sessions = append(s.sessions, sess)
			s.active = len(s.sessions) - 1
		}
		s.cfg.AddRecent(req.Path, string(req.Engine))
		s.Persist()
		s.Notify()
	})
}

func (s *AppState) CloseAt(i int) {
	if i < 0 || i >= len(s.sessions) {
		return
	}
	_ = s.sessions[i].Close()
	s.sessions = append(s.sessions[:i], s.sessions[i+1:]...)
	if s.active >= len(s.sessions) {
		s.active = len(s.sessions) - 1
	}
	if s.active < 0 {
		s.active = 0
	}
	s.page = pageHandlers{}
	s.Notify()
}

func (s *AppState) CloseActive() {
	if len(s.sessions) > 0 {
		s.CloseAt(s.active)
	}
}

func (s *AppState) CloseAll() {
	for _, sess := range s.sessions {
		_ = sess.Close()
	}
	s.sessions = nil
}

func (s *AppState) SelectTab(i int) {
	if i < 0 || i >= len(s.sessions) || i == s.active {
		return
	}
	s.active = i
	s.page = pageHandlers{}
	s.Notify()
}

func (s *AppState) CycleTab() {
	if len(s.sessions) <= 1 {
		return
	}
	s.SelectTab((s.active + 1) % len(s.sessions))
}

func (s *AppState) ShowOpenDialog() {
	showOpenDatabase(s.w, s.OpenSession)
}

func (s *AppState) ShowOpenDialogNewTab() {
	showOpenDatabase(s.w, func(req OpenRequest) {
		req.NewTab = true
		s.OpenSession(req)
	})
}

func (s *AppState) ShowSettings() {
	showSettings(s.w, s.themePref, i18n.Chosen(), SettingsHandlers{
		OnTheme:    s.SetTheme,
		OnLanguage: s.SetLanguage,
	})
}

func (s *AppState) SetLanguage(pref string) {
	i18n.SetLanguage(pref)
	s.cfg.Language = pref
	s.Persist()
	s.w.SetTitle(i18n.T("app.name"))
	s.Notify()
}

func (s *AppState) SetPageHandlers(h pageHandlers) { s.page = h }

func (s *AppState) ClearPageHandlers() { s.page = pageHandlers{} }

func (s *AppState) FireAddKey()      { fire(s.page.addKey) }
func (s *AppState) FireEditKey()     { fire(s.page.editKey) }
func (s *AppState) FireDeleteKey()   { fire(s.page.deleteKey) }
func (s *AppState) FireRefresh()     { fire(s.page.refresh) }
func (s *AppState) FireFocusFilter() { fire(s.page.focusFilter) }

func fire(fn func()) {
	if fn != nil {
		fn()
	}
}

func (s *AppState) Recents() []config.Recent { return s.cfg.Recents }

func (s *AppState) ApplyInitialWindowSize() {
	wpx, hpx := float32(1280), float32(800)
	if s.cfg.WindowWidth > 400 {
		wpx = s.cfg.WindowWidth
	}
	if s.cfg.WindowHeight > 300 {
		hpx = s.cfg.WindowHeight
	}
	s.w.Resize(fyne.NewSize(wpx, hpx))
}

func (s *AppState) SaveWindowSize() {
	size := s.w.Canvas().Size()
	s.cfg.WindowWidth = size.Width
	s.cfg.WindowHeight = size.Height
	s.Persist()
}

func (s *AppState) Persist() {
	if err := config.Save(s.cfg); err != nil {
		fyne.LogError("config save", err)
	}
}
