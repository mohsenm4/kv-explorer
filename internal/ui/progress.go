package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

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
