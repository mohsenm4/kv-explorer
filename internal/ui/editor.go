package ui

import (
	"bytes"
	"encoding/json"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

// valueEditor renders the editable value pane for a single entry. onSaved
// runs after a successful write so the parent can refresh the table.
func valueEditor(v fyne.ThemeVariant, sess *app.Session, entry kvstore.Entry, onSaved func()) fyne.CanvasObject {
	muted := themeColor(v, fynetheme.ColorNamePlaceHolder)
	fg := themeColor(v, fynetheme.ColorNameForeground)

	label := canvas.NewText("Value", muted)
	label.TextSize = 12

	keyText := canvas.NewText(string(entry.Key), fg)
	keyText.TextSize = 12
	keyText.TextStyle = fyne.TextStyle{Monospace: true}

	chip := engineChip(string(sess.Engine), v)
	header := container.NewBorder(nil, nil,
		container.NewHBox(label, keyText),
		chip,
		nil,
	)

	body := widget.NewMultiLineEntry()
	body.TextStyle = fyne.TextStyle{Monospace: true}
	body.Wrapping = fyne.TextWrapBreak
	body.SetText(displayValue(entry.Value))

	cancel := widget.NewButton("Cancel", func() {
		body.SetText(displayValue(entry.Value))
	})
	save := widget.NewButton("Save changes", func() {
		if err := sess.Store.Set(entry.Key, []byte(body.Text)); err != nil {
			fyne.LogError("save failed", err)
			return
		}
		if err := sess.Refresh(); err != nil {
			fyne.LogError("refresh failed", err)
		}
		if onSaved != nil {
			onSaved()
		}
	})
	save.Importance = widget.HighImportance

	footer := container.NewBorder(nil, nil, layout.NewSpacer(),
		container.NewHBox(cancel, save), nil)

	return container.NewBorder(container.NewPadded(header), container.NewPadded(footer), nil, nil, body)
}

// displayValue returns the value as a string, pretty-printed if it's JSON.
func displayValue(v []byte) string {
	var pretty bytes.Buffer
	if json.Indent(&pretty, v, "", "  ") == nil && pretty.Len() > 0 {
		return pretty.String()
	}
	return string(v)
}

// emptyEditor is the placeholder shown before any row is selected.
func emptyEditor(v fyne.ThemeVariant) fyne.CanvasObject {
	muted := themeColor(v, fynetheme.ColorNamePlaceHolder)
	t := canvas.NewText("Select a key to view its value", muted)
	t.TextSize = 12
	t.Alignment = fyne.TextAlignCenter
	return container.NewCenter(t)
}
