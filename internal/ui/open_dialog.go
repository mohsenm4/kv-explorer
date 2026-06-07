package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

type OpenRequest struct {
	Engine   kvstore.EngineKind
	Path     string
	NewTab   bool
	ReadOnly bool
}

var engineChoices = []struct {
	label string
	kind  kvstore.EngineKind
}{
	{"PebbleDB", kvstore.EnginePebble},
	{"BadgerDB", kvstore.EngineBadger},
	{"LevelDB", kvstore.EngineLevelDB},
}

func showOpenDatabase(parent fyne.Window, onConfirm func(OpenRequest)) {
	labels := make([]string, 0, len(engineChoices))
	for _, e := range engineChoices {
		labels = append(labels, e.label)
	}

	engine := widget.NewRadioGroup(labels, nil)
	engine.SetSelected(labels[0])

	path := widget.NewEntry()
	path.SetPlaceHolder("/path/to/database")
	path.TextStyle = fyne.TextStyle{Monospace: true}

	pick := widget.NewButtonWithIcon("", fynetheme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			path.SetText(uri.Path())
			// Engine mismatch is caught by the confirm step below; we
			// deliberately don't auto-switch the radio so the user's
			// choice stays explicit.
		}, parent)
	})

	pathRow := container.NewBorder(nil, nil, nil, pick, path)

	newTab := widget.NewCheck("Open in new tab", nil)
	newTab.SetChecked(true)
	readOnly := widget.NewCheck("Read-only", nil)

	body := container.NewVBox(
		sectionLabel("Engine"),
		engine,
		gap(8),
		sectionLabel("Path"),
		pathRow,
		gap(8),
		newTab,
		readOnly,
	)
	content := container.New(layout.NewCustomPaddedLayout(4, 4, 8, 8), body)

	d := dialog.NewCustomConfirm("Open Database", "Open", "Cancel", content, func(confirmed bool) {
		if !confirmed {
			return
		}
		kind := kvstore.EnginePebble
		for _, e := range engineChoices {
			if e.label == engine.Selected {
				kind = e.kind
				break
			}
		}
		req := OpenRequest{
			Engine:   kind,
			Path:     path.Text,
			NewTab:   newTab.Checked,
			ReadOnly: readOnly.Checked,
		}
		// Safety net for typed paths: if the folder looks like a
		// different engine, ask before trying to open with the wrong one.
		if detected, ok := kvstore.DetectEngine(path.Text); ok && detected != kind {
			dialog.ShowConfirm(
				"Engine mismatch",
				fmt.Sprintf("This folder looks like a %s database but you chose %s.\n\nOpen anyway?",
					engineLabelFor(detected), engine.Selected),
				func(yes bool) {
					if yes {
						onConfirm(req)
					}
				}, parent)
			return
		}
		onConfirm(req)
	}, parent)
	d.Resize(fyne.NewSize(520, 380))
	d.SetConfirmImportance(widget.HighImportance)
	installEscClose(parent, d)
	d.Show()
}

func engineLabelFor(k kvstore.EngineKind) string {
	for _, e := range engineChoices {
		if e.kind == k {
			return e.label
		}
	}
	return string(k)
}

func sectionLabel(text string) fyne.CanvasObject {
	th := fyne.CurrentApp().Settings().Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	t := canvas.NewText(text, th.Color(fynetheme.ColorNamePlaceHolder, v))
	t.TextSize = 11
	return t
}
