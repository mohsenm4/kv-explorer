package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/mohsenm4/kv-explorer/internal/i18n"
)

// mainMenu builds the File/Edit/View/Help menu. On macOS Fyne renders this
// in the system menu bar (top of screen); on Linux/Windows it sits inside
// the window as a strip under the title bar.
func mainMenu(s *AppState) *fyne.MainMenu {
	w := s.w

	openItem := fyne.NewMenuItem(i18n.T("menu.file.open"), s.ShowOpenDialog)
	openItem.Shortcut = shortcut(fyne.KeyO, fyne.KeyModifierShortcutDefault)

	closeItem := fyne.NewMenuItem(i18n.T("menu.file.close"), s.CloseActive)
	closeItem.Shortcut = shortcut(fyne.KeyW, fyne.KeyModifierShortcutDefault)

	settingsItem := fyne.NewMenuItem(i18n.T("menu.file.settings"), s.ShowSettings)
	settingsItem.Shortcut = shortcut(fyne.KeyComma, fyne.KeyModifierShortcutDefault)

	file := fyne.NewMenu(i18n.T("menu.file"),
		openItem,
		closeItem,
		fyne.NewMenuItemSeparator(),
		settingsItem,
	)
	edit := fyne.NewMenu(i18n.T("menu.edit"),
		disabledItem(i18n.T("menu.edit.addKey")),
		disabledItem(i18n.T("menu.edit.editKey")),
		disabledItem(i18n.T("menu.edit.deleteKey")),
	)
	view := fyne.NewMenu(i18n.T("menu.view"),
		disabledItem(i18n.T("menu.view.refresh")),
		fyne.NewMenuItem(i18n.T("menu.view.toggleTheme"), s.ToggleTheme),
	)
	help := fyne.NewMenu(i18n.T("menu.help"),
		fyne.NewMenuItem(i18n.T("menu.help.about"), func() {
			dialog.ShowInformation(i18n.T("app.name"),
				i18n.T("app.about.body"),
				w)
		}),
	)
	return fyne.NewMainMenu(file, edit, view, help)
}

func disabledItem(label string) *fyne.MenuItem {
	it := fyne.NewMenuItem(label, func() {})
	it.Disabled = true
	return it
}

func shortcut(key fyne.KeyName, mod fyne.KeyModifier) *desktop.CustomShortcut {
	return &desktop.CustomShortcut{KeyName: key, Modifier: mod}
}
