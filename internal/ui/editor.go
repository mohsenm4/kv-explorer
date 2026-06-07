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
	commit := func(data []byte) {
		withProgress(parent, "Saving…", func() error {
			if err := sess.Store.Set(entry.Key, data); err != nil {
				return err
			}
			return sess.Refresh()
		}, func(err error) {
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}
			if onSaved != nil {
				onSaved()
			}
		})
	}

	save := widget.NewButton("Save changes", func() {
		if current == nil {
			return
		}
		data, err := current()
		if err != nil {
			dialog.ShowError(err, parent)
			return
		}
		// Truncation guard: the visible buffer is a slice of the
		// original in both Text (>displayValueMax) and Hex
		// (>hexEditFormatMax) modes. Saving without a confirm would
		// silently destroy the tail.
		mode := format.Selected
		isTextMode := mode == "Text" || (mode == "Auto" && detected == KindText)
		isHexMode := mode == "Hex" || (mode == "Auto" && detected == KindBinary)
		truncated := (isTextMode && len(entry.Value) > displayValueMax) ||
			(isHexMode && len(entry.Value) > hexEditFormatMax)
		if truncated && len(data) < len(entry.Value) {
			dialog.ShowConfirm(
				"Replace value?",
				fmt.Sprintf("You're about to replace a %s value with %s. Anything past the visible buffer will be lost.\n\nUse Export to keep the full content.",
					humanSize(int64(len(entry.Value))), humanSize(int64(len(data)))),
				func(yes bool) {
					if yes {
						commit(data)
					}
				}, parent)
			return
		}
		commit(data)
	})
	save.Importance = widget.HighImportance
	if current == nil {
		save.Disable()
		cancel.Disable()
	}

	export := widget.NewButtonWithIcon("Export…", fynetheme.DownloadIcon(), func() {
		saver := dialog.NewFileSave(func(wc fyne.URIWriteCloser, err error) {
			if err != nil || wc == nil {
				return
			}
			withProgress(parent, "Exporting…", func() error {
				_, werr := wc.Write(entry.Value)
				wc.Close()
				return werr
			}, func(err error) {
				if err != nil {
					dialog.ShowError(err, parent)
				}
			})
		}, parent)
		saver.SetFileName(suggestedExportName(entry.Key, entry.Value))
		saver.Show()
	})

	footer := container.NewBorder(nil, nil, export,
		container.NewHBox(cancel, save), nil)

	center := container.NewBorder(
		container.NewPadded(formatBar),
		nil, nil, nil,
		container.NewVScroll(bodyStack),
	)

	return container.NewBorder(container.NewPadded(header), container.NewPadded(footer), nil, nil, center)
}

func textBody(value []byte) (fyne.CanvasObject, func() ([]byte, error), func()) {
	displayed := displayValue(value)
	be := widget.NewMultiLineEntry()
	be.TextStyle = fyne.TextStyle{Monospace: true}
	be.Wrapping = fyne.TextWrapBreak
	be.SetText(displayed)
	// Editable even when truncated. The outer Save handler warns before
	// overwriting a much larger original.
	return be,
		func() ([]byte, error) {
			// If the user didn't touch the text, return the original
			// bytes verbatim. Without this, displayValue's JSON
			// pretty-print would silently rewrite on-disk JSON to its
			// indented form on every open-and-save.
			if be.Text == displayed {
				return value, nil
			}
			return []byte(be.Text), nil
		},
		func() { be.SetText(displayed) }
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
			var data []byte
			withProgress(parent, "Loading file…", func() error {
				var ioErr error
				data, ioErr = io.ReadAll(rc)
				rc.Close()
				return ioErr
			}, func(err error) {
				if err != nil {
					dialog.ShowError(err, parent)
					return
				}
				staged = data
				pending = true
				refreshPreview()
			})
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

// hexEditFormatMax caps how many bytes the hex editor renders so huge
// values don't choke the text widget. Edits stay valid as long as the
// user keeps total bytes ≤ this many.
const hexEditFormatMax = 4096

func hexEditFormat(v []byte) string {
	truncated := false
	if len(v) > hexEditFormatMax {
		v = v[:hexEditFormatMax]
		truncated = true
	}
	var sb strings.Builder
	sb.Grow(len(v) * 3)
	for i, b := range v {
		sb.WriteString(fmt.Sprintf("%02x", b))
		switch {
		case (i+1)%16 == 0:
			sb.WriteByte('\n')
		case (i+1)%2 == 0:
			sb.WriteByte(' ')
		}
	}
	if truncated {
		sb.WriteString(fmt.Sprintf("\n... (showing first %d bytes; edits beyond this would be lost)", hexEditFormatMax))
	}
	return sb.String()
}

// displayValueMax caps text-mode rendering so Fyne's MultiLineEntry
// doesn't have to lay out megabytes of glyphs (it isn't virtualised). 16
// KiB is the largest the widget renders without noticeable lag.
// Values larger than this are shown read-only with a note pointing at
// Export for external editing.
const displayValueMax = 16 * 1024

// suggestedExportName returns a filename derived from the key with the
// extension that matches the value's content type. Slashes in the key
// become underscores so the result is a single file. Keys that already
// look like they have an extension keep it; otherwise we sniff the bytes.
func suggestedExportName(key, value []byte) string {
	name := strings.ReplaceAll(string(key), "/", "_")
	if name == "" {
		name = "value"
	}
	if dot := strings.LastIndexByte(name, '.'); dot > 0 && dot >= len(name)-6 {
		return name
	}
	return name + extensionForBytes(value)
}

func extensionForBytes(v []byte) string {
	_, mime := DetectContent(v)
	switch {
	case strings.HasPrefix(mime, "image/jpeg"):
		return ".jpg"
	case strings.HasPrefix(mime, "image/png"):
		return ".png"
	case strings.HasPrefix(mime, "image/gif"):
		return ".gif"
	case strings.HasPrefix(mime, "image/webp"):
		return ".webp"
	case strings.HasPrefix(mime, "image/"):
		return ".img"
	case mime == "application/json":
		return ".json"
	case strings.HasPrefix(mime, "text/html"):
		return ".html"
	case strings.HasPrefix(mime, "text/xml"), mime == "application/xml":
		return ".xml"
	case strings.HasPrefix(mime, "text/"):
		return ".txt"
	case strings.HasPrefix(mime, "audio/mpeg"):
		return ".mp3"
	case strings.HasPrefix(mime, "audio/wav"), strings.HasPrefix(mime, "audio/x-wav"):
		return ".wav"
	case strings.HasPrefix(mime, "audio/"):
		return ".audio"
	case strings.HasPrefix(mime, "video/mp4"):
		return ".mp4"
	case strings.HasPrefix(mime, "video/"):
		return ".video"
	case mime == "application/pdf":
		return ".pdf"
	case mime == "application/zip":
		return ".zip"
	case mime == "application/gzip":
		return ".gz"
	default:
		return ".bin"
	}
}

// displayValue returns the value as a string, pretty-printed if it's JSON.
// Values larger than displayValueMax are truncated so the editor stays
// responsive.
func displayValue(v []byte) string {
	if len(v) > displayValueMax {
		return string(v[:displayValueMax]) +
			fmt.Sprintf("\n\n... (showing first %s of %s — use Export to save the whole value to a file)",
				humanSize(int64(displayValueMax)), humanSize(int64(len(v))))
	}
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
