package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

// keyTable renders the central table of key / value preview / size. Cells
// pull from sess.Entries() filtered by the shared FilterState on every
// render, so calling sess.Refresh()/table.Refresh() picks up edits and
// filter changes alike.
// onSelect fires when the user picks a row.
func keyTable(sess *app.Session, filter *FilterState, onSelect func(kvstore.Entry)) *widget.Table {
	headers := []string{"Key", "Value preview", "Size"}

	read := func() []kvstore.Entry {
		entries, _ := sess.Entries()
		return applyFilter(entries, *filter)
	}

	table := widget.NewTableWithHeaders(
		func() (int, int) { return len(read()), 3 },
		func() fyne.CanvasObject {
			l := widget.NewLabel("")
			l.TextStyle = fyne.TextStyle{Monospace: true}
			l.Truncation = fyne.TextTruncateEllipsis
			return l
		},
		func(id widget.TableCellID, c fyne.CanvasObject) {
			l := c.(*widget.Label)
			l.TextStyle = fyne.TextStyle{Monospace: true}
			l.Alignment = fyne.TextAlignLeading
			entries := read()
			if id.Row < 0 || id.Row >= len(entries) {
				l.SetText("")
				return
			}
			e := entries[id.Row]
			switch id.Col {
			case 0:
				l.SetText(string(e.Key))
			case 1:
				l.SetText(previewValue(e.Value))
			case 2:
				l.Alignment = fyne.TextAlignTrailing
				l.SetText(humanSize(int64(len(e.Value))))
			}
		},
	)

	table.CreateHeader = func() fyne.CanvasObject {
		l := widget.NewLabel("")
		l.TextStyle = fyne.TextStyle{Bold: true}
		return l
	}
	table.UpdateHeader = func(id widget.TableCellID, c fyne.CanvasObject) {
		l := c.(*widget.Label)
		l.TextStyle = fyne.TextStyle{Bold: true}
		l.Alignment = fyne.TextAlignLeading
		if id.Row == -1 && id.Col >= 0 && id.Col < len(headers) {
			if id.Col == 2 {
				l.Alignment = fyne.TextAlignTrailing
			}
			l.SetText(headers[id.Col])
		}
	}

	table.SetColumnWidth(0, 360)
	table.SetColumnWidth(1, 320)
	table.SetColumnWidth(2, 80)

	table.OnSelected = func(id widget.TableCellID) {
		entries := read()
		if onSelect != nil && id.Row >= 0 && id.Row < len(entries) {
			onSelect(entries[id.Row])
		}
	}

	return table
}

// previewValue trims a value down to a single-line preview for the table.
func previewValue(v []byte) string {
	const max = 120
	s := string(v)
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' || s[i] == '\r' || s[i] == '\t' {
			s = s[:i] + " " + s[i+1:]
		}
	}
	if len(s) > max {
		return s[:max] + "…"
	}
	return s
}
