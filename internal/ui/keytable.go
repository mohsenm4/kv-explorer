package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

type tableEntry struct {
	key   []byte
	value []byte
}

func loadEntries(s kvstore.Store) ([]tableEntry, error) {
	it, err := s.Iter(nil)
	if err != nil {
		return nil, err
	}
	defer it.Close()
	var out []tableEntry
	for it.Next() {
		e := it.Entry()
		out = append(out, tableEntry{key: e.Key, value: e.Value})
	}
	return out, nil
}

// keyTable renders the central table of key / value preview / size.
// onSelect is called when the user picks a row (wired into the editor
// in Step 9).
func keyTable(sess *app.Session, onSelect func(tableEntry)) fyne.CanvasObject {
	entries, _ := loadEntries(sess.Store)

	headers := []string{"Key", "Value preview", "Size"}

	table := widget.NewTableWithHeaders(
		func() (int, int) { return len(entries), 3 },
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
			e := entries[id.Row]
			switch id.Col {
			case 0:
				l.SetText(string(e.key))
			case 1:
				l.SetText(previewValue(e.value))
			case 2:
				l.Alignment = fyne.TextAlignTrailing
				l.SetText(humanSize(int64(len(e.value))))
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
	// Collapse newlines so multi-line JSON shows as one row.
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
