package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/i18n"
)

// ToolbarActions groups the callbacks the top toolbar can fire. nil
// callbacks render their button as disabled.
type ToolbarActions struct {
	OnOpen     func()
	OnClose    func()
	OnAdd      func()
	OnEdit     func()
	OnDelete   func()
	OnRefresh  func()
	OnSettings func()
}

// toolbarHandles lets the parent flip Edit/Delete enable state when the
// table selection changes.
type toolbarHandles struct {
	bar       fyne.CanvasObject
	editBtn   *widget.Button
	deleteBtn *widget.Button
}

func buildToolbar(actions ToolbarActions) toolbarHandles {
	open := toolbarButton(i18n.T("toolbar.open"), fynetheme.FolderOpenIcon(), actions.OnOpen)
	closeBtn := toolbarButton(i18n.T("toolbar.close"), fynetheme.CancelIcon(), actions.OnClose)
	add := toolbarButton(i18n.T("toolbar.add"), fynetheme.ContentAddIcon(), actions.OnAdd)
	edit := toolbarButton(i18n.T("toolbar.edit"), fynetheme.DocumentCreateIcon(), actions.OnEdit)
	del := toolbarButton(i18n.T("toolbar.delete"), fynetheme.DeleteIcon(), actions.OnDelete)
	refresh := toolbarButton(i18n.T("toolbar.refresh"), fynetheme.ViewRefreshIcon(), actions.OnRefresh)
	settings := toolbarButton("", fynetheme.SettingsIcon(), actions.OnSettings)

	row := container.NewHBox(
		open, closeBtn,
		toolbarSep(),
		add, edit, del,
		toolbarSep(),
		refresh,
		layout.NewSpacer(),
		settings,
	)
	bg := canvas.NewRectangle(fyne.CurrentApp().Settings().Theme().Color(
		fynetheme.ColorNameHeaderBackground,
		fyne.CurrentApp().Settings().ThemeVariant(),
	))
	return toolbarHandles{
		bar:       container.NewStack(bg, container.NewPadded(row)),
		editBtn:   edit,
		deleteBtn: del,
	}
}

func toolbarButton(label string, icon fyne.Resource, action func()) *widget.Button {
	b := widget.NewButtonWithIcon(label, icon, action)
	b.Importance = widget.LowImportance
	if action == nil {
		b.Disable()
	}
	return b
}

func toolbarSep() fyne.CanvasObject {
	r := canvas.NewRectangle(fyne.CurrentApp().Settings().Theme().Color(
		fynetheme.ColorNameSeparator,
		fyne.CurrentApp().Settings().ThemeVariant(),
	))
	r.SetMinSize(fyne.NewSize(1, 20))
	return container.NewCenter(r)
}
