package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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

	left := placeholderPane("Tree (Step 8)")
	table := keyTable(sess, func(e tableEntry) {
		// TODO Step 9: surface the value in the editor
		_ = e
	})
	editor := placeholderPane("Editor (Step 9)")
	center := container.NewBorder(nil, editor, nil, nil, table)
	split := container.NewHSplit(left, center)
	split.Offset = 0.22

	status := mainStatusBar(v, sess, onToggle)

	top := container.NewVBox(accent, bar, tabs)
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
