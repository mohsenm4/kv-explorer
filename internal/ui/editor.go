package ui

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/i18n"
	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

func valueEditor(v fyne.ThemeVariant, sess *app.Session, parent fyne.Window, entry kvstore.Entry, onSaved func(), onExtSession func(*externalEditSession)) fyne.CanvasObject {
	muted := themeColor(v, fynetheme.ColorNamePlaceHolder)
	fg := themeColor(v, fynetheme.ColorNameForeground)

	label := canvas.NewText(i18n.T("editor.label"), muted)
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
	isJSON := mime == "application/json"

	bodyStack := container.NewStack()
	var current func() ([]byte, error)
	var reset func()

	rebuild := func(mode string) {
		var body fyne.CanvasObject
		switch mode {
		case "Text":
			body, current, reset = textBody(entry.Value)
		case "Tree":
			body, current, reset = jsonTreeBody(v, entry.Value)
		case "Hex":
			body, current, reset = hexBody(v, entry.Value, mime)
		case "Image":
			body, current, reset = imageBody(v, parent, entry.Value, mime)
		default: // Auto
			switch {
			case detected == KindImage:
				body, current, reset = imageBody(v, parent, entry.Value, mime)
			case detected == KindBinary:
				body, current, reset = hexBody(v, entry.Value, mime)
			case isJSON:
				body, current, reset = jsonTreeBody(v, entry.Value)
			default:
				body, current, reset = textBody(entry.Value)
			}
		}
		bodyStack.Objects = []fyne.CanvasObject{body}
		bodyStack.Refresh()
	}

	formatIDs := []string{"Auto", "Text", "Tree", "Hex", "Image"}
	formatLabels := []string{
		i18n.T("editor.format.auto"),
		i18n.T("editor.format.text"),
		i18n.T("editor.format.tree"),
		i18n.T("editor.format.hex"),
		i18n.T("editor.format.image"),
	}
	labelToID := map[string]string{}
	for i, id := range formatIDs {
		labelToID[formatLabels[i]] = id
	}
	format := widget.NewRadioGroup(formatLabels, func(s string) {
		rebuild(labelToID[s])
	})
	format.Horizontal = true
	format.SetSelected(formatLabels[0])
	rebuild("Auto")

	formatLabel := canvas.NewText(i18n.T("editor.format"), muted)
	formatLabel.TextSize = 11

	detection := canvas.NewText(i18n.Tf("editor.detected", map[string]any{"Mime": mime}), muted)
	detection.TextSize = 11

	formatBar := container.NewBorder(nil, nil,
		container.NewHBox(formatLabel, format),
		detection,
		nil,
	)

	cancel := widget.NewButton(i18n.T("editor.cancel"), func() {
		if reset != nil {
			reset()
		}
	})
	commit := func(data []byte) {
		withProgress(parent, i18n.T("progress.saving"), func() error {
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

	save := widget.NewButton(i18n.T("editor.save"), func() {
		if current == nil {
			return
		}
		data, err := current()
		if err != nil {
			dialog.ShowError(err, parent)
			return
		}
		// In Text/Hex modes the visible buffer is truncated; saving without a confirm would lose the tail.
		mode := labelToID[format.Selected]
		isTextMode := mode == "Text" || (mode == "Auto" && detected == KindText)
		isHexMode := mode == "Hex" || (mode == "Auto" && detected == KindBinary)
		truncated := (isTextMode && len(entry.Value) > displayValueMax) ||
			(isHexMode && len(entry.Value) > hexEditFormatMax)
		if truncated && len(data) < len(entry.Value) {
			dialog.ShowConfirm(
				i18n.T("editor.replaceConfirm.title"),
				i18n.Tf("editor.replaceConfirm.body", map[string]any{
					"Old": humanSize(int64(len(entry.Value))),
					"New": humanSize(int64(len(data))),
				}),
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

	export := widget.NewButtonWithIcon(i18n.T("editor.export"), fynetheme.DownloadIcon(), func() {
		saver := dialog.NewFileSave(func(wc fyne.URIWriteCloser, err error) {
			if err != nil || wc == nil {
				return
			}
			withProgress(parent, i18n.T("progress.exporting"), func() error {
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

	pendingBanner := newPendingBanner(i18n.T("editor.externalPending"))
	pendingBanner.Hide()

	var currentExt *externalEditSession

	applyBtn := widget.NewButtonWithIcon(i18n.T("editor.apply"), fynetheme.ConfirmIcon(), nil)
	applyBtn.Importance = widget.HighImportance
	applyBtn.Hide()

	discardBtn := widget.NewButton(i18n.T("editor.discard"), nil)
	discardBtn.Hide()

	hidePending := func() {
		pendingBanner.Hide()
		applyBtn.Hide()
		discardBtn.Hide()
	}

	applyBtn.OnTapped = func() {
		if currentExt == nil {
			return
		}
		data, err := os.ReadFile(currentExt.Path())
		if err != nil {
			dialog.ShowError(err, parent)
			return
		}
		withProgress(parent, i18n.T("progress.saving"), func() error {
			if err := sess.Store.Set(entry.Key, data); err != nil {
				return err
			}
			return sess.Refresh()
		}, func(err error) {
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}
			hidePending()
			if onSaved != nil {
				onSaved()
			}
		})
	}

	discardBtn.OnTapped = hidePending

	doOpenExt := func() {
		if currentExt != nil {
			currentExt.Close()
			currentExt = nil
		}
		es, err := startExternalEditSession(entry.Key, entry.Value, func() {
			pendingBanner.Show()
			applyBtn.Show()
			discardBtn.Show()
		})
		if err != nil {
			dialog.ShowError(err, parent)
			return
		}
		currentExt = es
		if onExtSession != nil {
			onExtSession(es)
		}
	}

	openExt := widget.NewButtonWithIcon(i18n.T("editor.openInVSCode"), fynetheme.ComputerIcon(), func() {
		if len(entry.Value) > externalSizeWarn {
			dialog.ShowConfirm(
				i18n.T("editor.openInVSCode"),
				i18n.Tf("editor.externalSizeWarn", map[string]any{
					"Size": humanSize(int64(len(entry.Value))),
				}),
				func(ok bool) {
					if ok {
						doOpenExt()
					}
				}, parent)
			return
		}
		doOpenExt()
	})

	pendingArea := container.NewHBox(pendingBanner, applyBtn, discardBtn)

	footer := container.NewBorder(nil, nil, container.NewHBox(export, openExt),
		container.NewHBox(cancel, save), container.NewCenter(pendingArea))

	center := container.NewBorder(
		container.NewPadded(formatBar),
		nil, nil, nil,
		container.NewVScroll(bodyStack),
	)

	return container.NewBorder(container.NewPadded(header), container.NewPadded(footer), nil, nil, center)
}

// Visible amber pill so the user notices the file changed externally. Color is fixed (not theme-derived) so it stands out in both light and dark.
func newPendingBanner(text string) *fyne.Container {
	bg := canvas.NewRectangle(color.NRGBA{R: 0xea, G: 0x86, B: 0x0c, A: 0xff})
	bg.CornerRadius = 8

	label := canvas.NewText(text, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	label.TextStyle = fyne.TextStyle{Bold: true}
	label.TextSize = 13

	icon := canvas.NewText("⚠", color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	icon.TextStyle = fyne.TextStyle{Bold: true}
	icon.TextSize = 14

	row := container.NewHBox(icon, label)
	return container.NewStack(bg, container.NewPadded(row))
}

func textBody(value []byte) (fyne.CanvasObject, func() ([]byte, error), func()) {
	displayed := displayValue(value)
	be := widget.NewMultiLineEntry()
	be.TextStyle = fyne.TextStyle{Monospace: true}
	be.Wrapping = fyne.TextWrapBreak
	be.SetText(displayed)

	body := fyne.CanvasObject(be)
	if hint := timestampHint(value); hint != "" {
		muted := themeColor(fyne.CurrentApp().Settings().ThemeVariant(), fynetheme.ColorNamePlaceHolder)
		label := canvas.NewText(hint, muted)
		label.TextSize = 11
		body = container.NewBorder(nil, container.NewPadded(label), nil, nil, be)
	}

	return body,
		func() ([]byte, error) {
			// Untouched text returns the original bytes so JSON pretty-print doesn't silently rewrite on save.
			if be.Text == displayed {
				return value, nil
			}
			return []byte(be.Text), nil
		},
		func() { be.SetText(displayed) }
}

func imageBody(v fyne.ThemeVariant, parent fyne.Window, value []byte, mime string) (fyne.CanvasObject, func() ([]byte, error), func()) {
	muted := themeColor(v, fynetheme.ColorNamePlaceHolder)

	staged := value
	pending := false

	res := fyne.NewStaticResource("value", value)
	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(200, 200))

	info := canvas.NewText(imageInfoText(mime, value, false), muted)
	info.TextSize = 11

	refreshPreview := func() {
		img.Resource = fyne.NewStaticResource("value", staged)
		img.Refresh()
		_, m := DetectContent(staged)
		info.Text = imageInfoText(m, staged, pending)
		info.Refresh()
	}

	replace := widget.NewButtonWithIcon(i18n.T("editor.replace"), fynetheme.UploadIcon(), func() {
		dialog.ShowFileOpen(func(rc fyne.URIReadCloser, err error) {
			if err != nil || rc == nil {
				return
			}
			var data []byte
			withProgress(parent, i18n.T("progress.loadingFile"), func() error {
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

	info := canvas.NewText(i18n.Tf("editor.hexInfo", map[string]any{
		"Mime": mime,
		"Size": humanSize(int64(len(value))),
	}), muted)
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

// imageInfoText picks the i18n template with dimensions when DecodeConfig
// succeeds, otherwise falls back to mime + size.
func imageInfoText(mime string, value []byte, pending bool) string {
	w, h, ok := imageDimensions(value)
	data := map[string]any{"Mime": mime, "Size": humanSize(int64(len(value)))}
	if ok {
		data["Width"] = w
		data["Height"] = h
		if pending {
			return i18n.Tf("editor.imageInfoDimsPending", data)
		}
		return i18n.Tf("editor.imageInfoDims", data)
	}
	if pending {
		return i18n.Tf("editor.imageInfoPending", data)
	}
	return i18n.Tf("editor.imageInfo", data)
}

func imageDimensions(v []byte) (int, int, bool) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(v))
	if err != nil {
		return 0, 0, false
	}
	return cfg.Width, cfg.Height, true
}

// Cap on hex editor render so MultiLineEntry (non-virtualised) doesn't choke on huge values.
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

// 16 KiB is the largest MultiLineEntry renders without lag (it isn't virtualised).
const displayValueMax = 16 * 1024

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
	// http.DetectContentType often labels JSON as text/plain; check directly.
	if len(v) > 0 && len(v) < 4<<20 && json.Valid(v) {
		return ".json"
	}
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

func emptyEditor(v fyne.ThemeVariant) fyne.CanvasObject {
	muted := themeColor(v, fynetheme.ColorNamePlaceHolder)
	t := canvas.NewText(i18n.T("editor.placeholder"), muted)
	t.TextSize = 12
	t.Alignment = fyne.TextAlignCenter
	return container.NewCenter(t)
}
