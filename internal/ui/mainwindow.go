package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/kvstore"
	apptheme "github.com/mohsenm4/kv-explorer/internal/ui/theme"
)

func mainPage(a fyne.App, w fyne.Window, sess *app.Session, variant *fyne.ThemeVariant, onOpen func(), onClose, onToggle func()) fyne.CanvasObject {
	v := *variant

	accent := canvas.NewRectangle(apptheme.DBAccent(string(sess.Engine), v))
	accent.SetMinSize(fyne.NewSize(0, 3))

	bar := buildToolbar(onOpen, onClose)
	tabs := tabStrip(v, sess)

	editorBox := container.NewStack(emptyEditor(v))
	filter := &FilterState{}

	var table *widget.Table
	table = keyTable(sess, filter, func(e kvstore.Entry) {
		editorBox.Objects = []fyne.CanvasObject{valueEditor(v, sess, e, func() {
			table.Refresh()
		})}
		editorBox.Refresh()
	})

	filterUI := filterRow(filter, func() {
		table.Refresh()
	})

	left := prefixTree(sess, func(key []byte) {
		val, err := sess.Store.Get(key)
		if err != nil {
			return
		}
		editorBox.Objects = []fyne.CanvasObject{valueEditor(v, sess, kvstore.Entry{Key: key, Value: val}, func() {
			table.Refresh()
		})}
		editorBox.Refresh()
	})

	tableWithFilter := container.NewBorder(container.NewPadded(filterUI), nil, nil, nil, table)
	center := container.NewVSplit(tableWithFilter, editorBox)
	center.Offset = 0.62

	split := container.NewHSplit(left, center)
	split.Offset = 0.22

	status := mainStatusBar(v, sess, onToggle)

	sep := canvas.NewRectangle(themeColor(v, fynetheme.ColorNameSeparator))
	sep.SetMinSize(fyne.NewSize(0, 1))

	top := container.NewVBox(accent, bar, tabs, sep)
	return container.NewBorder(top, status, nil, nil, split)
}

func placeholderPane(label string) fyne.CanvasObject {
	l := widget.NewLabel(label)
	l.Importance = widget.LowImportance
	return container.NewCenter(l)
}

func engineDisplayName(k kvstore.EngineKind) string {
	for _, e := range engineChoices {
		if e.kind == k {
			return e.label
		}
	}
	return string(k)
}
