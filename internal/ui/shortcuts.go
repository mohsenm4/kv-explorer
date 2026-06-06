package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// registerShortcuts wires canvas-level shortcuts to AppState. Per-page
// actions (add/edit/delete/refresh/focus) flow through Fire* methods so
// they no-op when there's no session.
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
		// Cmd+Tab is owned by the OS on macOS, so cycle on Ctrl+Tab
		// everywhere instead. Settings → Shortcuts pane mirrors this.
		{shortcut(fyne.KeyTab, fyne.KeyModifierControl), s.CycleTab},
	}

	for _, b := range bindings {
		b := b
		c.AddShortcut(b.sc, func(_ fyne.Shortcut) { b.run() })
	}
}
