package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func statusBar(v fyne.ThemeVariant, message string, onToggle func()) fyne.CanvasObject {
	th := fyne.CurrentApp().Settings().Theme()
	bgCol := th.Color(fynetheme.ColorNameHeaderBackground, v)
	muted := th.Color(fynetheme.ColorNamePlaceHolder, v)

	bg := canvas.NewRectangle(bgCol)

	msg := canvas.NewText(message, muted)
	msg.TextSize = 11

	toggleLabel := "Light"
	if v == fynetheme.VariantDark {
		toggleLabel = "Dark"
	}
	toggle := widget.NewButton(toggleLabel, onToggle)
	toggle.Importance = widget.MediumImportance

	row := container.NewBorder(nil, nil, container.NewPadded(msg), container.NewPadded(toggle), layout.NewSpacer())

	return container.NewStack(bg, row)
}
