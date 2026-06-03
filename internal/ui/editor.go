package ui

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

// valueEditor renders the editable value pane for a single entry. The body
// is chosen by content kind: text gets an editor, images get a preview
// with Replace…, arbitrary binary gets editable hex. Save is the same
// outer button for text/hex; image saves via its own Replace flow.
func valueEditor(v fyne.ThemeVariant, sess *app.Session, parent fyne.Window, entry kvstore.Entry, onSaved func()) fyne.CanvasObject {
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

	kind, mime := DetectContent(entry.Value)
	var body fyne.CanvasObject
	var current func() ([]byte, error) // nil = outer Save disabled (image handles itself)
	var reset func()

	switch kind {
	case KindImage:
		body = imageBody(v, sess, parent, entry, onSaved)
	case KindBinary:
		body, current, reset = hexBody(v, entry.Value, mime)
	default:
		body, current, reset = textBody(entry.Value)
	}

	cancel := widget.NewButton("Cancel", reset)
	save := widget.NewButton("Save changes", func() {
		if current == nil {
			return
		}
		data, err := current()
		if err != nil {
			dialog.ShowError(err, parent)
			return
		}
		if err := sess.Store.Set(entry.Key, data); err != nil {
			dialog.ShowError(err, parent)
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
	if current == nil {
		save.Disable()
		cancel.Disable()
	}

	footer := container.NewBorder(nil, nil, layout.NewSpacer(),
		container.NewHBox(cancel, save), nil)

	return container.NewBorder(container.NewPadded(header), container.NewPadded(footer), nil, nil, body)
}

func textBody(value []byte) (fyne.CanvasObject, func() ([]byte, error), func()) {
	be := widget.NewMultiLineEntry()
	be.TextStyle = fyne.TextStyle{Monospace: true}
	be.Wrapping = fyne.TextWrapBreak
	be.SetText(displayValue(value))
	return be,
		func() ([]byte, error) { return []byte(be.Text), nil },
		func() { be.SetText(displayValue(value)) }
}

func imageBody(v fyne.ThemeVariant, sess *app.Session, parent fyne.Window, entry kvstore.Entry, onSaved func()) fyne.CanvasObject {
	muted := themeColor(v, fynetheme.ColorNamePlaceHolder)

	_, mime := DetectContent(entry.Value)
	res := fyne.NewStaticResource("value", entry.Value)
	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(200, 200))

	info := canvas.NewText(fmt.Sprintf("%s · %s", mime, humanSize(int64(len(entry.Value)))), muted)
	info.TextSize = 11

	replace := widget.NewButtonWithIcon("Replace…", fynetheme.UploadIcon(), func() {
		dialog.ShowFileOpen(func(rc fyne.URIReadCloser, err error) {
			if err != nil || rc == nil {
				return
			}
			defer rc.Close()
			data, ioErr := io.ReadAll(rc)
			if ioErr != nil {
				dialog.ShowError(ioErr, parent)
				return
			}
			if err := sess.Store.Set(entry.Key, data); err != nil {
				dialog.ShowError(err, parent)
				return
			}
			if err := sess.Refresh(); err != nil {
				fyne.LogError("refresh failed", err)
			}
			img.Resource = fyne.NewStaticResource("value", data)
			img.Refresh()
			_, newMime := DetectContent(data)
			info.Text = fmt.Sprintf("%s · %s", newMime, humanSize(int64(len(data))))
			info.Refresh()
			if onSaved != nil {
				onSaved()
			}
		}, parent)
	})

	bottom := container.NewBorder(nil, nil, container.NewPadded(info), replace, nil)
	return container.NewBorder(nil, bottom, nil, nil, img)
}

func hexBody(v fyne.ThemeVariant, value []byte, mime string) (fyne.CanvasObject, func() ([]byte, error), func()) {
	muted := themeColor(v, fynetheme.ColorNamePlaceHolder)

	text := widget.NewMultiLineEntry()
	text.TextStyle = fyne.TextStyle{Monospace: true}
	text.Wrapping = fyne.TextWrapBreak
	text.SetText(hexEditFormat(value))

	info := canvas.NewText(fmt.Sprintf("%s · %s · hex editable", mime, humanSize(int64(len(value)))), muted)
	info.TextSize = 11

	body := container.NewBorder(nil, container.NewPadded(info), nil, nil, text)

	current := func() ([]byte, error) {
		clean := strings.Map(func(r rune) rune {
			if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
				return -1
			}
			return r
		}, text.Text)
		return hex.DecodeString(clean)
	}
	reset := func() { text.SetText(hexEditFormat(value)) }
	return body, current, reset
}

func hexEditFormat(v []byte) string {
	var sb strings.Builder
	for i, b := range v {
		sb.WriteString(fmt.Sprintf("%02x", b))
		switch {
		case (i+1)%16 == 0:
			sb.WriteByte('\n')
		case (i+1)%2 == 0:
			sb.WriteByte(' ')
		}
	}
	return sb.String()
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
