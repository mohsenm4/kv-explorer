package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func registerShortcuts(w fyne.Window, s *AppState) {
	c := w.Canvas()

	type binding struct {
		sc  *desktop.CustomShortcut
		run func()
	}

	bindings := []binding{
		{shortcut(fyne.KeyN, fyne.KeyModifierShortcutDefault), s.FireAddKey},
		{shortcut(fyne.KeyF, fyne.KeyModifierShortcutDefault), s.FireFocusFilter},
		{shortcut(fyne.KeyF2, 0), s.FireEditKey},
		{shortcut(fyne.KeyF5, 0), s.FireRefresh},
		{shortcut(fyne.KeyDelete, 0), s.FireDeleteKey},
		{shortcut(fyne.KeyBackspace, 0), s.FireDeleteKey},
		// Cmd+Tab is OS-reserved on macOS; use Ctrl+Tab everywhere instead.
		{shortcut(fyne.KeyTab, fyne.KeyModifierControl), s.CycleTab},
	}

	for _, b := range bindings {
		b := b
		c.AddShortcut(b.sc, func(_ fyne.Shortcut) { b.run() })
	}
}
