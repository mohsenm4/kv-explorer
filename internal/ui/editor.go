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

// valueEditor renders the editable value pane for a single entry. A
// Format toggle (Auto / Text / Hex / Image) lets the user pick a body
// regardless of the auto-detected content kind. Save commits the staged
// bytes; Cancel reverts to the original.
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

	detected, mime := DetectContent(entry.Value)

	bodyStack := container.NewStack()
	var current func() ([]byte, error)
	var reset func()

	rebuild := func(mode string) {
		var body fyne.CanvasObject
		switch mode {
		case "Text":
			body, current, reset = textBody(entry.Value)
		case "Hex":
			body, current, reset = hexBody(v, entry.Value, mime)
		case "Image":
			body, current, reset = imageBody(v, parent, entry.Value, mime)
		default: // Auto
			switch detected {
			case KindImage:
				body, current, reset = imageBody(v, parent, entry.Value, mime)
			case KindBinary:
				body, current, reset = hexBody(v, entry.Value, mime)
			default:
				body, current, reset = textBody(entry.Value)
			}
		}
		bodyStack.Objects = []fyne.CanvasObject{body}
		bodyStack.Refresh()
	}

	format := widget.NewRadioGroup([]string{"Auto", "Text", "Hex", "Image"}, func(s string) {
		rebuild(s)
	})
	format.Horizontal = true
	format.SetSelected("Auto")
	rebuild("Auto")

	formatLabel := canvas.NewText("Format:", muted)
	formatLabel.TextSize = 11

	detection := canvas.NewText(fmt.Sprintf("Detected: %s", mime), muted)
	detection.TextSize = 11

	formatBar := container.NewBorder(nil, nil,
		container.NewHBox(formatLabel, format),
		detection,
		nil,
	)

	cancel := widget.NewButton("Cancel", func() {
		if reset != nil {
			reset()
		}
	})
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

	footer := container.NewBorder(nil, nil, layout.NewSpacer(),
		container.NewHBox(cancel, save), nil)

	center := container.NewBorder(
		container.NewPadded(formatBar),
		nil, nil, nil,
		container.NewVScroll(bodyStack),
	)

	return container.NewBorder(container.NewPadded(header), container.NewPadded(footer), nil, nil, center)
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

// imageBody renders an image preview plus a Replace… button that stages
// new bytes. Save commits staged → store; Cancel reverts to the original.
// Replace accepts any file, so picking a text/.txt also works — after
// save the next selection will auto-detect the new content.
func imageBody(v fyne.ThemeVariant, parent fyne.Window, value []byte, mime string) (fyne.CanvasObject, func() ([]byte, error), func()) {
	muted := themeColor(v, fynetheme.ColorNamePlaceHolder)

	staged := value
	pending := false

	res := fyne.NewStaticResource("value", value)
	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(200, 200))

	info := canvas.NewText(fmt.Sprintf("%s · %s", mime, humanSize(int64(len(value)))), muted)
	info.TextSize = 11

	refreshPreview := func() {
		img.Resource = fyne.NewStaticResource("value", staged)
		img.Refresh()
		_, m := DetectContent(staged)
		suffix := ""
		if pending {
			suffix = " · pending"
		}
		info.Text = fmt.Sprintf("%s · %s%s", m, humanSize(int64(len(staged))), suffix)
		info.Refresh()
	}

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
			staged = data
			pending = true
			refreshPreview()
		}, parent)
	})

	bottom := container.NewBorder(nil, nil, container.NewPadded(info), replace, nil)
	body := container.NewBorder(nil, bottom, nil, nil, img)

	current := func() ([]byte, error) {
		pending = false
		return staged, nil
	}
	resetFn := func() {
		staged = value
		pending = false
		refreshPreview()
	}
	return body, current, resetFn
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
	resetFn := func() { text.SetText(hexEditFormat(value)) }
	return body, current, resetFn
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
