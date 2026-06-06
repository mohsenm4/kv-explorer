package ui

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

// showAddKey opens the Add key dialog. onSaved fires after a successful
// write so the parent can refresh views.
func showAddKey(parent fyne.Window, sess *app.Session, onSaved func()) {
	showKeyDialog(parent, sess, "Add key", nil, nil, onSaved)
}

// showEditKey opens the Edit key dialog prefilled with the entry's
// current key and value.
func showEditKey(parent fyne.Window, sess *app.Session, entry kvstore.Entry, onSaved func()) {
	showKeyDialog(parent, sess, "Edit key", entry.Key, entry.Value, onSaved)
}

func showKeyDialog(parent fyne.Window, sess *app.Session, title string, oldKey, oldValue []byte, onSaved func()) {
	editing := oldKey != nil

	keyEntry := widget.NewEntry()
	keyEntry.TextStyle = fyne.TextStyle{Monospace: true}
	keyEntry.SetPlaceHolder("e.g. users/0042")

	keyMode := widget.NewRadioGroup([]string{"Text", "Hex"}, nil)
	keyMode.Horizontal = true
	keyMode.SetSelected("Text")

	if editing {
		// Default key mode based on whether oldKey is valid UTF-8 without
		// control bytes — hash-like binary keys land in Hex.
		if _, mime := DetectContent(oldKey); strings.HasPrefix(mime, "text/") {
			keyEntry.SetText(string(oldKey))
		} else {
			keyMode.SetSelected("Hex")
			keyEntry.SetText(hex.EncodeToString(oldKey))
		}
	}

	valueEntry := widget.NewMultiLineEntry()
	valueEntry.TextStyle = fyne.TextStyle{Monospace: true}
	valueEntry.Wrapping = fyne.TextWrapBreak
	valueEntry.SetPlaceHolder("Type a value or click \"Use file…\" to upload one")

	// File-upload state: when non-nil, the dialog will write these bytes
	// instead of the value entry's text. For Edit on a binary value we
	// pre-stage the existing bytes so we never try to render megabytes of
	// garbled text inside MultiLineEntry.
	var staged []byte
	autoStage := false
	if editing {
		if kind, _ := DetectContent(oldValue); kind != KindText {
			staged = oldValue
			autoStage = true
		} else {
			valueEntry.SetText(displayValue(oldValue))
		}
	}

	initialSize := len(valueEntry.Text)
	if autoStage {
		initialSize = len(staged)
	}
	sizeReadout := widget.NewLabel(fmt.Sprintf("%d B", initialSize))
	sizeReadout.Importance = widget.LowImportance
	valueEntry.OnChanged = func(s string) {
		sizeReadout.SetText(fmt.Sprintf("%d B", len(s)))
	}

	fileInfo := widget.NewLabel("")
	fileInfo.Importance = widget.LowImportance

	clearFile := widget.NewButton("Clear", nil)
	clearFile.Hide()
	clearFile.OnTapped = func() {
		staged = nil
		fileInfo.SetText("")
		clearFile.Hide()
		valueEntry.Enable()
		sizeReadout.SetText(fmt.Sprintf("%d B", len(valueEntry.Text)))
	}

	useFile := widget.NewButtonWithIcon("Use file…", fynetheme.UploadIcon(), func() {
		dialog.ShowFileOpen(func(rc fyne.URIReadCloser, err error) {
			if err != nil || rc == nil {
				return
			}
			name := rc.URI().Name()
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
				_, mime := DetectContent(data)
				fileInfo.SetText(fmt.Sprintf("Loaded %s · %s · %s", name, mime, humanSize(int64(len(data)))))
				clearFile.Show()
				valueEntry.Disable()
				sizeReadout.SetText(fmt.Sprintf("%d B", len(data)))
			})
		}, parent)
	})

	exportBtn := widget.NewButtonWithIcon("Export…", fynetheme.DownloadIcon(), func() {
		data := staged
		if data == nil {
			data = []byte(valueEntry.Text)
		}
		if len(data) == 0 {
			return
		}
		saver := dialog.NewFileSave(func(wc fyne.URIWriteCloser, err error) {
			if err != nil || wc == nil {
				return
			}
			withProgress(parent, "Exporting…", func() error {
				_, werr := wc.Write(data)
				wc.Close()
				return werr
			}, func(err error) {
				if err != nil {
					dialog.ShowError(err, parent)
				}
			})
		}, parent)
		nameKey := oldKey
		if nameKey == nil {
			nameKey = []byte(keyEntry.Text)
		}
		if len(nameKey) == 0 {
			nameKey = []byte("value")
		}
		saver.SetFileName(suggestedExportName(nameKey, data))
		saver.Show()
	})

	fileRow := container.NewBorder(nil, nil, useFile, container.NewHBox(exportBtn, clearFile), fileInfo)

	if autoStage {
		_, mime := DetectContent(staged)
		fileInfo.SetText(fmt.Sprintf("Current value: %s · %s", mime, humanSize(int64(len(staged)))))
		clearFile.Show()
		valueEntry.Disable()
	}

	content := container.NewVBox(
		sectionLabel("Key"),
		keyEntry,
		keyMode,
		gap(6),
		sectionLabel("Value"),
		valueEntry,
		fileRow,
		sizeReadout,
	)

	confirmLabel := "Add"
	if editing {
		confirmLabel = "Save changes"
	}

	d := dialog.NewCustomConfirm(title, confirmLabel, "Cancel", content, func(ok bool) {
		if !ok {
			return
		}
		key, err := parseKey(keyEntry.Text, keyMode.Selected)
		if err != nil {
			dialog.ShowError(err, parent)
			return
		}
		if len(key) == 0 {
			dialog.ShowError(errors.New("key cannot be empty"), parent)
			return
		}

		// Duplicate check — only matters when adding or when the edit
		// changed the key.
		if !editing || string(key) != string(oldKey) {
			if _, err := sess.Store.Get(key); err == nil {
				dialog.ShowError(fmt.Errorf("key %q already exists", string(key)), parent)
				return
			}
		}

		newValue := staged
		if newValue == nil {
			newValue = []byte(valueEntry.Text)
		}
		withProgress(parent, "Saving…", func() error {
			if err := sess.Store.Set(key, newValue); err != nil {
				return err
			}
			if editing && string(key) != string(oldKey) {
				if err := sess.Store.Delete(oldKey); err != nil {
					return err
				}
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
	}, parent)
	d.Resize(fyne.NewSize(560, 460))
	d.SetConfirmImportance(widget.HighImportance)
	// Esc closes the dialog
	parent.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyEscape},
		func(_ fyne.Shortcut) { d.Hide() })
	d.Show()
}

// parseKey converts the raw entry text into bytes given the selected mode.
func parseKey(text, mode string) ([]byte, error) {
	if mode == "Hex" {
		clean := strings.Map(func(r rune) rune {
			if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
				return -1
			}
			return r
		}, text)
		b, err := hex.DecodeString(clean)
		if err != nil {
			return nil, fmt.Errorf("hex key: %w", err)
		}
		return b, nil
	}
	return []byte(text), nil
}
