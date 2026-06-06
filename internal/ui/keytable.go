package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

// keyTable renders the central table of key / value preview / size. It
// caches the filtered slice between cell renders so applyFilter doesn't
// run 90+ times per refresh.
// onSelect fires with the row's key + freshly-fetched value.
func keyTable(sess *app.Session, filter *FilterState, onSelect func(kvstore.Entry)) *widget.Table {
	headers := []string{"Key", "Value preview", "Size"}

	type cacheKey struct {
		n int
		q string
	}
	var cached []app.KeyMeta
	var ck cacheKey

	read := func() []app.KeyMeta {
		keys, _ := sess.Keys()
		cur := cacheKey{len(keys), filter.Query}
		if cached == nil || cur != ck {
			cached = applyFilter(keys, *filter)
			ck = cur
		}
		return cached
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
			rows := read()
			if id.Row < 0 || id.Row >= len(rows) {
				l.SetText("")
				return
			}
			meta := rows[id.Row]
			switch id.Col {
			case 0:
				l.SetText(meta.Key)
			case 1:
				l.SetText(meta.Preview)
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
		rows := read()
		if id.Row < 0 || id.Row >= len(rows) {
			return
		}
		k := []byte(rows[id.Row].Key)
		v, err := sess.Value(k)
		if err != nil {
			return
		}
		onSelect(kvstore.Entry{Key: k, Value: v})
	}

	return table
}
