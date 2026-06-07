package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// installEscClose wires Esc to hide d, and unregisters the shortcut when
// the dialog closes. Without the unregister, every dialog instance leaves
// a stale Hide-closure pointing at a dead dialog on the canvas, and a
// later Esc press resolves to the dead handler instead of the live one.
func installEscClose(parent fyne.Window, d dialog.Dialog) {
	sc := &desktop.CustomShortcut{KeyName: fyne.KeyEscape}
	parent.Canvas().AddShortcut(sc, func(_ fyne.Shortcut) { d.Hide() })
	d.SetOnClosed(func() {
		parent.Canvas().RemoveShortcut(sc)
	})
}

// withProgress runs fn on a goroutine while showing a non-dismissable
// progress dialog over parent. When fn returns, the dialog closes and
// done is invoked on the UI thread with the error (or nil).
//
// The pattern keeps the UI responsive but prevents interaction with
// other widgets until the operation completes, so a slow Save can't be
// interrupted midway.
func withProgress(parent fyne.Window, message string, fn func() error, done func(error)) {
	bar := widget.NewProgressBarInfinite()

	content := container.NewVBox(
		widget.NewLabel(message),
		bar,
	)

	d := dialog.NewCustomWithoutButtons("Working…", content, parent)
	d.Resize(fyne.NewSize(360, 120))
	d.Show()

	go func() {
		err := fn()
		fyne.Do(func() {
			d.Hide()
			if done != nil {
				done(err)
			}
		})
	}()
}
