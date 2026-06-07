package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	apptheme "github.com/mohsenm4/kv-explorer/internal/ui/theme"
)

// engineDot uses a "●" glyph (not a Circle shape) so it shares the label's baseline; pass the surrounding label's TextSize.
func engineDot(engine string, v fyne.ThemeVariant, size float32) fyne.CanvasObject {
	t := canvas.NewText("●", apptheme.DBAccent(engine, v))
	t.TextSize = size
	return t
}
