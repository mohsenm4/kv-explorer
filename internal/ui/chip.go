package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	apptheme "github.com/mohsenm4/kv-explorer/internal/ui/theme"
)

func engineChip(engine string, v fyne.ThemeVariant) fyne.CanvasObject {
	bg := canvas.NewRectangle(apptheme.DBAccentTint(engine, v))
	bg.CornerRadius = 999

	label := canvas.NewText(strings.ToUpper(engine), apptheme.DBAccent(engine, v))
	label.TextSize = 10
	label.TextStyle = fyne.TextStyle{Bold: true}
	label.Alignment = fyne.TextAlignCenter

	padded := container.New(layout.NewCustomPaddedLayout(2, 2, 10, 10), label)
	stack := container.NewStack(bg, padded)
	return container.NewGridWrap(fyne.NewSize(72, 22), stack)
}
