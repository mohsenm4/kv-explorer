package ui

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	if editing {
		valueEntry.SetText(displayValue(oldValue))
	}

	sizeReadout := widget.NewLabel(fmt.Sprintf("%d B", len(valueEntry.Text)))
	sizeReadout.Importance = widget.LowImportance
	valueEntry.OnChanged = func(s string) {
		sizeReadout.SetText(fmt.Sprintf("%d B", len(s)))
	}

	content := container.NewVBox(
		sectionLabel("Key"),
		keyEntry,
		keyMode,
		gap(6),
		sectionLabel("Value"),
		valueEntry,
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

		newValue := []byte(valueEntry.Text)
		if err := sess.Store.Set(key, newValue); err != nil {
			dialog.ShowError(err, parent)
			return
		}
		if editing && string(key) != string(oldKey) {
			// Key was renamed — drop the old one.
			if err := sess.Store.Delete(oldKey); err != nil {
				dialog.ShowError(err, parent)
				return
			}
		}
		if err := sess.Refresh(); err != nil {
			fyne.LogError("refresh failed", err)
		}
		if onSaved != nil {
			onSaved()
		}
	}, parent)
	d.Resize(fyne.NewSize(560, 460))
	d.SetConfirmImportance(widget.HighImportance)
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
