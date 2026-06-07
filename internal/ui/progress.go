package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// installEscClose wires Esc to hide d and unregisters on close — otherwise stale Hide-closures pile up on the canvas and later Esc presses fire dead handlers.
func installEscClose(parent fyne.Window, d dialog.Dialog) {
	sc := &desktop.CustomShortcut{KeyName: fyne.KeyEscape}
	parent.Canvas().AddShortcut(sc, func(_ fyne.Shortcut) { d.Hide() })
	d.SetOnClosed(func() {
		parent.Canvas().RemoveShortcut(sc)
	})
}

// withProgress runs fn on a goroutine behind a modal progress dialog so a slow operation can't be interrupted midway; done fires on the UI thread.
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
