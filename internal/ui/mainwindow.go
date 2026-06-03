package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/kvstore"
	apptheme "github.com/mohsenm4/kv-explorer/internal/ui/theme"
)

func mainPage(a fyne.App, w fyne.Window, bar TabBar, variant *fyne.ThemeVariant, onOpen, onClose, onToggle, onSettings func(), handlers *appHandlers) fyne.CanvasObject {
	v := *variant
	sess := bar.Sessions[bar.Active]

	accent := canvas.NewRectangle(apptheme.DBAccent(string(sess.Engine), v))
	accent.SetMinSize(fyne.NewSize(0, 3))

	editorBox := container.NewStack(emptyEditor(v))
	treeBox := container.NewStack()
	filter := &FilterState{}

	var current *kvstore.Entry
	var table *widget.Table
	var toolbar toolbarHandles

	loadEditorFor := func(e kvstore.Entry) {
		current = &e
		editorBox.Objects = []fyne.CanvasObject{valueEditor(v, sess, w, e, func() {
			table.Refresh()
		})}
		editorBox.Refresh()
		if toolbar.editBtn != nil {
			toolbar.editBtn.Enable()
			toolbar.deleteBtn.Enable()
		}
	}

	clearSelection := func() {
		current = nil
		editorBox.Objects = []fyne.CanvasObject{emptyEditor(v)}
		editorBox.Refresh()
		if toolbar.editBtn != nil {
			toolbar.editBtn.Disable()
			toolbar.deleteBtn.Disable()
		}
	}

	rebuildTree := func() {
		treeBox.Objects = []fyne.CanvasObject{prefixTree(sess, func(key []byte) {
			val, err := sess.Store.Get(key)
			if err != nil {
				return
			}
			loadEditorFor(kvstore.Entry{Key: key, Value: val})
		})}
		treeBox.Refresh()
	}

	refreshAll := func() {
		_ = sess.Refresh()
		table.Refresh()
		rebuildTree()
		clearSelection()
	}

	table = keyTable(sess, filter, loadEditorFor)
	rebuildTree()

	filterUI, filterEntry := filterRow(filter, func() {
		table.Refresh()
	})

	actions := ToolbarActions{
		OnOpen:  onOpen,
		OnClose: onClose,
		OnAdd: func() {
			showAddKey(w, sess, refreshAll)
		},
		OnEdit: func() {
			if current == nil {
				return
			}
			showEditKey(w, sess, *current, refreshAll)
		},
		OnDelete: func() {
			if current == nil {
				return
			}
			showDeleteKey(w, sess, current.Key, refreshAll)
		},
		OnRefresh:  refreshAll,
		OnSettings: onSettings,
	}
	toolbar = buildToolbar(actions)
	toolbar.editBtn.Disable()
	toolbar.deleteBtn.Disable()

	handlers.addKey = actions.OnAdd
	handlers.editKey = actions.OnEdit
	handlers.deleteKey = actions.OnDelete
	handlers.refresh = refreshAll
	handlers.focusFilter = func() { w.Canvas().Focus(filterEntry) }

	tabs := tabStrip(v, bar)

	tableWithFilter := container.NewBorder(container.NewPadded(filterUI), nil, nil, nil, table)
	center := container.NewVSplit(tableWithFilter, editorBox)
	center.Offset = 0.62

	split := container.NewHSplit(treeBox, center)
	split.Offset = 0.22

	status := mainStatusBar(v, sess, onToggle)

	sep := canvas.NewRectangle(themeColor(v, fynetheme.ColorNameSeparator))
	sep.SetMinSize(fyne.NewSize(0, 1))

	top := container.NewVBox(accent, toolbar.bar, tabs, sep)
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
