package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
)

// mainMenu builds the File/Edit/View/Help menu. On macOS Fyne renders this
// in the system menu bar (top of screen); on Linux/Windows it sits inside
// the window as a strip under the title bar.
func mainMenu(s *AppState) *fyne.MainMenu {
	w := s.w

	openItem := fyne.NewMenuItem("Open…", s.ShowOpenDialog)
	openItem.Shortcut = shortcut(fyne.KeyO, fyne.KeyModifierShortcutDefault)

	closeItem := fyne.NewMenuItem("Close", s.CloseActive)
	closeItem.Shortcut = shortcut(fyne.KeyW, fyne.KeyModifierShortcutDefault)

	settingsItem := fyne.NewMenuItem("Settings…", s.ShowSettings)
	settingsItem.Shortcut = shortcut(fyne.KeyComma, fyne.KeyModifierShortcutDefault)

	file := fyne.NewMenu("File",
		openItem,
		closeItem,
		fyne.NewMenuItemSeparator(),
		settingsItem,
	)
	edit := fyne.NewMenu("Edit",
		disabledItem("Add Key…"),
		disabledItem("Edit Key…"),
		disabledItem("Delete Key"),
	)
	view := fyne.NewMenu("View",
		disabledItem("Refresh"),
		fyne.NewMenuItem("Toggle Theme", s.ToggleTheme),
	)
	help := fyne.NewMenu("Help",
		fyne.NewMenuItem("About KV-Explorer", func() {
			dialog.ShowInformation("KV-Explorer",
				"A desktop GUI tool for managing key-value databases.\nPebbleDB · BadgerDB · LevelDB",
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
