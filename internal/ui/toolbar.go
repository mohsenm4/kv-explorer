package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func buildToolbar(onOpen, onClose func()) fyne.CanvasObject {
	open := toolbarButton("Open", fynetheme.FolderOpenIcon(), onOpen)
	closeBtn := toolbarButton("Close", fynetheme.CancelIcon(), onClose)
	add := toolbarButton("Add", fynetheme.ContentAddIcon(), nil)
	edit := toolbarButton("Edit", fynetheme.DocumentCreateIcon(), nil)
	del := toolbarButton("Delete", fynetheme.DeleteIcon(), nil)
	refresh := toolbarButton("Refresh", fynetheme.ViewRefreshIcon(), nil)
	settings := toolbarButton("", fynetheme.SettingsIcon(), nil)

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
	return container.NewStack(bg, container.NewPadded(row))
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
