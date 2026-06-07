package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/i18n"
	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

func welcomePage(s *AppState) fyne.CanvasObject {
	s.ClearPageHandlers()

	a := s.a
	v := s.Variant()
	th := a.Settings().Theme()

	fg := th.Color(fynetheme.ColorNameForeground, v)
	muted := th.Color(fynetheme.ColorNamePlaceHolder, v)
	primary := th.Color(fynetheme.ColorNamePrimary, v)

	hero := heroIcon(primary)

	title := canvas.NewText(i18n.T("app.name"), fg)
	title.TextSize = 22
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	tagline := canvas.NewText(i18n.T("app.tagline"), muted)
	tagline.TextSize = 14
	tagline.Alignment = fyne.TextAlignCenter

	open := widget.NewButtonWithIcon(i18n.T("welcome.openDatabase"), fynetheme.FolderOpenIcon(), func() {
		s.ShowOpenDialog()
	})
	open.Importance = widget.HighImportance

	recents := recentsFromConfig(s.Recents())

	openRecent := widget.NewButtonWithIcon(i18n.T("welcome.openRecent"), fynetheme.MenuDropDownIcon(), nil)
	openRecent.IconPlacement = widget.ButtonIconTrailingText
	if len(recents) == 0 {
		openRecent.Disable()
	} else {
		openRecent.OnTapped = func() {
			showRecentMenu(s.w, openRecent, recents, func(r recentEntry) {
				s.OpenSession(OpenRequest{Engine: kvstore.EngineKind(r.engine), Path: r.path})
			})
		}
	}

	actions := container.NewHBox(layout.NewSpacer(), open, openRecent, layout.NewSpacer())

	recentBlock := buildRecentBlock(v, fg, muted, recents, func(r recentEntry) {
		s.OpenSession(OpenRequest{Engine: kvstore.EngineKind(r.engine), Path: r.path})
	})

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

	return container.NewBorder(nil, welcomeStatusBar(v, s.ToggleTheme), nil, nil, center)
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
