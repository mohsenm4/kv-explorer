package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// mainMenu builds the File/Edit/View/Help menu. On macOS Fyne renders this
// in the system menu bar (top of screen); on Linux/Windows it sits inside
// the window as a strip under the title bar.
func mainMenu(w fyne.Window, openDialog, closeSession, toggleTheme func()) *fyne.MainMenu {
	file := fyne.NewMenu("File",
		fyne.NewMenuItem("Open…", openDialog),
		fyne.NewMenuItem("Close", closeSession),
	)
	edit := fyne.NewMenu("Edit",
		disabledItem("Add Key…"),
		disabledItem("Edit Key…"),
		disabledItem("Delete Key"),
	)
	view := fyne.NewMenu("View",
		disabledItem("Refresh"),
		fyne.NewMenuItem("Toggle Theme", toggleTheme),
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
