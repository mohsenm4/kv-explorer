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
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/i18n"
	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

func showAddKey(parent fyne.Window, sess *app.Session, onSaved func()) {
	showKeyDialog(parent, sess, i18n.T("addKey.title"), nil, nil, onSaved)
}

func showEditKey(parent fyne.Window, sess *app.Session, entry kvstore.Entry, onSaved func()) {
	showKeyDialog(parent, sess, i18n.T("editKey.title"), entry.Key, entry.Value, onSaved)
}

func showKeyDialog(parent fyne.Window, sess *app.Session, title string, oldKey, oldValue []byte, onSaved func()) {
	editing := oldKey != nil

	keyEntry := widget.NewEntry()
	keyEntry.TextStyle = fyne.TextStyle{Monospace: true}
	keyEntry.SetPlaceHolder(i18n.T("keyDialog.keyPlaceholder"))

	textLabel := i18n.T("editor.format.text")
	hexLabel := i18n.T("editor.format.hex")
	keyMode := widget.NewRadioGroup([]string{textLabel, hexLabel}, nil)
	keyMode.Horizontal = true
	keyMode.SetSelected(textLabel)

	if editing {
		if _, mime := DetectContent(oldKey); strings.HasPrefix(mime, "text/") {
			keyEntry.SetText(string(oldKey))
		} else {
			keyMode.SetSelected(hexLabel)
			keyEntry.SetText(hex.EncodeToString(oldKey))
		}
	}

	valueEntry := widget.NewMultiLineEntry()
	valueEntry.TextStyle = fyne.TextStyle{Monospace: true}
	valueEntry.Wrapping = fyne.TextWrapBreak
	valueEntry.SetPlaceHolder(i18n.T("keyDialog.valuePlaceholder"))

	// On Edit, pre-stage binary values so we never try to render megabytes of garbled text in MultiLineEntry.
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

	clearFile := widget.NewButton(i18n.T("keyDialog.clear"), nil)
	clearFile.Hide()
	clearFile.OnTapped = func() {
		staged = nil
		fileInfo.SetText("")
		clearFile.Hide()
		valueEntry.Enable()
		sizeReadout.SetText(fmt.Sprintf("%d B", len(valueEntry.Text)))
	}

	useFile := widget.NewButtonWithIcon(i18n.T("keyDialog.useFile"), fynetheme.UploadIcon(), func() {
		dialog.ShowFileOpen(func(rc fyne.URIReadCloser, err error) {
			if err != nil || rc == nil {
				return
			}
			name := rc.URI().Name()
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
				_, mime := DetectContent(data)
				fileInfo.SetText(i18n.Tf("keyDialog.loadedInfo", map[string]any{
					"Name": name,
					"Mime": mime,
					"Size": humanSize(int64(len(data))),
				}))
				clearFile.Show()
				valueEntry.Disable()
				sizeReadout.SetText(fmt.Sprintf("%d B", len(data)))
			})
		}, parent)
	})

	exportBtn := widget.NewButtonWithIcon(i18n.T("keyDialog.export"), fynetheme.DownloadIcon(), func() {
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
			withProgress(parent, i18n.T("progress.exporting"), func() error {
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
		fileInfo.SetText(i18n.Tf("keyDialog.currentInfo", map[string]any{
			"Mime": mime,
			"Size": humanSize(int64(len(staged))),
		}))
		clearFile.Show()
		valueEntry.Disable()
	}

	content := container.NewVBox(
		sectionLabel(i18n.T("keyDialog.keyLabel")),
		keyEntry,
		keyMode,
		gap(6),
		sectionLabel(i18n.T("keyDialog.valueLabel")),
		valueEntry,
		fileRow,
		sizeReadout,
	)

	confirmLabel := i18n.T("keyDialog.confirmAdd")
	if editing {
		confirmLabel = i18n.T("keyDialog.confirmSave")
	}

	d := dialog.NewCustomConfirm(title, confirmLabel, i18n.T("keyDialog.cancel"), content, func(ok bool) {
		if !ok {
			return
		}
		modeID := "Text"
		if keyMode.Selected == hexLabel {
			modeID = "Hex"
		}
		key, err := parseKey(keyEntry.Text, modeID)
		if err != nil {
			dialog.ShowError(err, parent)
			return
		}
		if len(key) == 0 {
			dialog.ShowError(errors.New(i18n.T("keyDialog.error.emptyKey")), parent)
			return
		}

		if !editing || string(key) != string(oldKey) {
			if _, err := sess.Store.Get(key); err == nil {
				dialog.ShowError(errors.New(i18n.Tf("keyDialog.error.duplicateKey", map[string]any{"Key": fmt.Sprintf("%q", string(key))})), parent)
				return
			}
		}

		newValue := staged
		if newValue == nil {
			newValue = []byte(valueEntry.Text)
		}
		withProgress(parent, i18n.T("progress.saving"), func() error {
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
	installEscClose(parent, d)
	d.Show()
}

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
			return nil, fmt.Errorf("%s: %w", i18n.T("keyDialog.error.hexKey"), err)
		}
		return b, nil
	}
	return []byte(text), nil
}
