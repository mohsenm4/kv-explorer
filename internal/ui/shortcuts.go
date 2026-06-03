package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// registerShortcuts wires canvas-level shortcuts to the cross-cutting
// handlers. Each handler is nil-checked so shortcuts no-op when there's
// no session.
func registerShortcuts(w fyne.Window, h *appHandlers) {
	c := w.Canvas()

	type binding struct {
		sc  *desktop.CustomShortcut
		run func()
	}

	bindings := []binding{
		{shortcut(fyne.KeyN, fyne.KeyModifierShortcutDefault), func() {
			if h.addKey != nil {
				h.addKey()
			}
		}},
		{shortcut(fyne.KeyF, fyne.KeyModifierShortcutDefault), func() {
			if h.focusFilter != nil {
				h.focusFilter()
			}
		}},
		{shortcut(fyne.KeyF2, 0), func() {
			if h.editKey != nil {
				h.editKey()
			}
		}},
		{shortcut(fyne.KeyF5, 0), func() {
			if h.refresh != nil {
				h.refresh()
			}
		}},
		{shortcut(fyne.KeyDelete, 0), func() {
			if h.deleteKey != nil {
				h.deleteKey()
			}
		}},
		{shortcut(fyne.KeyBackspace, 0), func() {
			if h.deleteKey != nil {
				h.deleteKey()
			}
		}},
	}

	for _, b := range bindings {
		b := b
		c.AddShortcut(b.sc, func(_ fyne.Shortcut) { b.run() })
	}
}
