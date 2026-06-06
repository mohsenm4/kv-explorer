package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
)

// showDeleteKey asks the user to confirm a destructive delete of a single
// key. onDeleted runs after a successful delete.
func showDeleteKey(parent fyne.Window, sess *app.Session, key []byte, onDeleted func()) {
	th := fyne.CurrentApp().Settings().Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	fg := th.Color(fynetheme.ColorNameForeground, v)
	muted := th.Color(fynetheme.ColorNamePlaceHolder, v)

	intro := canvas.NewText("This will permanently delete:", fg)
	intro.TextSize = 12

	keyDisplay := canvas.NewText(displayKey(key), fg)
	keyDisplay.TextSize = 13
	keyDisplay.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}

	warn := canvas.NewText("The action cannot be undone.", muted)
	warn.TextSize = 11

	content := container.NewVBox(
		intro,
		container.NewPadded(keyDisplay),
		warn,
	)

	d := dialog.NewCustomConfirm("Delete key?", "Delete", "Cancel", content, func(ok bool) {
		if !ok {
			return
		}
		withProgress(parent, "Deleting…", func() error {
			if err := sess.Store.Delete(key); err != nil {
				return err
			}
			return sess.Refresh()
		}, func(err error) {
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}
			if onDeleted != nil {
				onDeleted()
			}
		})
	}, parent)
	d.Resize(fyne.NewSize(480, 240))
	d.SetConfirmImportance(widget.DangerImportance)
	d.Show()
}

// displayKey formats a key for human display — text when valid UTF-8,
// otherwise a short hex prefix.
func displayKey(k []byte) string {
	if _, mime := DetectContent(k); mime == "text/plain" || mime[:5] == "text/" {
		return string(k)
	}
	const maxBytes = 32
	if len(k) > maxBytes {
		return hexShort(k[:maxBytes]) + "…"
	}
	return hexShort(k)
}

func hexShort(b []byte) string {
	const hexc = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, x := range b {
		out[i*2] = hexc[x>>4]
		out[i*2+1] = hexc[x&0x0f]
	}
	return string(out)
}
