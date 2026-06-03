package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

// keyTable renders the central table of key / value preview / size. It
// holds only key metadata in memory and fetches values on demand from the
// store, so a million-key database doesn't have to fit in RAM.
// onSelect fires with the row's key + freshly-fetched value.
func keyTable(sess *app.Session, filter *FilterState, onSelect func(kvstore.Entry)) *widget.Table {
	headers := []string{"Key", "Value preview", "Size"}

	read := func() []app.KeyMeta {
		keys, _ := sess.Keys()
		return applyFilter(keys, *filter)
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
			keys := read()
			if id.Row < 0 || id.Row >= len(keys) {
				l.SetText("")
				return
			}
			meta := keys[id.Row]
			switch id.Col {
			case 0:
				l.SetText(meta.Key)
			case 1:
				if v, err := sess.Value([]byte(meta.Key)); err == nil {
					l.SetText(previewValue(v))
				} else {
					l.SetText("")
				}
			case 2:
				l.Alignment = fyne.TextAlignTrailing
				l.SetText(humanSize(int64(meta.Size)))
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
		if onSelect == nil {
			return
		}
		keys := read()
		if id.Row < 0 || id.Row >= len(keys) {
			return
		}
		k := []byte(keys[id.Row].Key)
		v, err := sess.Value(k)
		if err != nil {
			return
		}
		onSelect(kvstore.Entry{Key: k, Value: v})
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
