package ui

import (
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/i18n"
)

// FilterState carries the current filter query. Mode is a placeholder for
// the future picker (prefix / substring / regex / value) — for now every
// query is treated as a case-sensitive key substring.
type FilterState struct {
	Query string
}

// filterRow returns the search row that sits above the key table plus the
// Entry widget itself so shortcuts can focus it.
// onChange fires (debounced) after the user stops typing.
func filterRow(state *FilterState, onChange func()) (fyne.CanvasObject, *widget.Entry) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(i18n.T("filter.placeholder"))

	var timer *time.Timer
	entry.OnChanged = func(s string) {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(250*time.Millisecond, func() {
			fyne.Do(func() {
				state.Query = s
				if onChange != nil {
					onChange()
				}
			})
		})
	}

	leading := widget.NewIcon(fynetheme.SearchIcon())
	search := container.NewBorder(nil, nil, container.NewPadded(leading), nil, entry)

	filterBtn := widget.NewButtonWithIcon(i18n.T("filter.button"), fynetheme.MenuExpandIcon(), func() {
		// TODO: picker for prefix / substring / regex / value modes
	})
	filterBtn.IconPlacement = widget.ButtonIconLeadingText

	return container.NewBorder(nil, nil, nil, filterBtn, search), entry
}

// applyFilter returns the subset of keys that match the filter state.
func applyFilter(keys []app.KeyMeta, state FilterState) []app.KeyMeta {
	q := state.Query
	if q == "" {
		return keys
	}
	out := keys[:0:0]
	for _, k := range keys {
		if strings.Contains(k.Key, q) {
			out = append(out, k)
		}
	}
	return out
}
