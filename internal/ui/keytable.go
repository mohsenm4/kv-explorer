package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/i18n"
	"github.com/mohsenm4/kv-explorer/internal/kvstore"
	apptheme "github.com/mohsenm4/kv-explorer/internal/ui/theme"
)

// Filtered rows are cached between cell renders so applyFilter doesn't run 90+ times per refresh.
func keyTable(sess *app.Session, filter *FilterState, v fyne.ThemeVariant, onSelect func(kvstore.Entry)) *widget.Table {
	headers := []string{
		i18n.T("table.header.key"),
		i18n.T("table.header.valuePreview"),
		i18n.T("table.header.size"),
	}

	type cacheKey struct {
		rev int
		q   string
	}
	var cached []app.KeyMeta
	var ck cacheKey

	read := func() []app.KeyMeta {
		keys, err := sess.Keys()
		if err != nil {
			fyne.LogError("load keys for table", err)
		}
		cur := cacheKey{sess.Rev(), filter.Query}
		if cached == nil || cur != ck {
			cached = applyFilter(keys, *filter)
			ck = cur
		}
		return cached
	}

	table := widget.NewTableWithHeaders(
		func() (int, int) { return len(read()), 3 },
		func() fyne.CanvasObject {
			badge := canvas.NewText("", color.Transparent)
			badge.TextSize = 11
			badge.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
			l := widget.NewLabel("")
			l.TextStyle = fyne.TextStyle{Monospace: true}
			l.Truncation = fyne.TextTruncateEllipsis
			return container.NewBorder(nil, nil, badge, nil, l)
		},
		func(id widget.TableCellID, c fyne.CanvasObject) {
			badge, l := badgeAndLabel(c)
			l.TextStyle = fyne.TextStyle{Monospace: true}
			l.Alignment = fyne.TextAlignLeading
			badge.Text = ""
			badge.Color = color.Transparent
			rows := read()
			if id.Row < 0 || id.Row >= len(rows) {
				l.SetText("")
				badge.Refresh()
				return
			}
			meta := rows[id.Row]
			switch id.Col {
			case 0:
				l.SetText(meta.Key)
			case 1:
				if meta.Kind != "" {
					badge.Text = "[" + meta.Kind + "] "
					badge.Color = apptheme.KindAccent(meta.Kind, v)
				}
				l.SetText(meta.Preview)
			case 2:
				l.Alignment = fyne.TextAlignTrailing
				l.SetText(humanSize(int64(meta.Size)))
			}
			badge.Refresh()
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
	table.SetColumnWidth(1, 360)
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
		val, err := sess.Value(k)
		if err != nil {
			return
		}
		onSelect(kvstore.Entry{Key: k, Value: val})
	}

	return table
}

// badgeAndLabel extracts the two children we put in the cell container. We walk
// Objects rather than relying on a fixed index so a layout change in the
// factory doesn't silently break rendering.
func badgeAndLabel(c fyne.CanvasObject) (*canvas.Text, *widget.Label) {
	cont := c.(*fyne.Container)
	var badge *canvas.Text
	var l *widget.Label
	for _, o := range cont.Objects {
		switch t := o.(type) {
		case *canvas.Text:
			badge = t
		case *widget.Label:
			l = t
		}
	}
	return badge, l
}
