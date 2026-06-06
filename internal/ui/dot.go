package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	apptheme "github.com/mohsenm4/kv-explorer/internal/ui/theme"
)

// engineDot renders a small filled circle in the engine's accent color
// as a "●" glyph so it shares the same baseline as the label it sits next
// to. Pass the same size as the surrounding label's TextSize for tight
// vertical alignment.
func engineDot(engine string, v fyne.ThemeVariant, size float32) fyne.CanvasObject {
	t := canvas.NewText("●", apptheme.DBAccent(engine, v))
	t.TextSize = size
	return t
}
