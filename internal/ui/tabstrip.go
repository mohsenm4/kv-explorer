package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
	apptheme "github.com/mohsenm4/kv-explorer/internal/ui/theme"
)

// tabStrip renders the row of database tabs. For now there's only one
// session, so just a single active tab plus a "+" affordance.
func tabStrip(v fyne.ThemeVariant, sess *app.Session) fyne.CanvasObject {
	accent := apptheme.DBAccent(string(sess.Engine), v)
	fg := themeColor(v, fynetheme.ColorNameForeground)

	dot := engineDot(string(sess.Engine), v, 14)

	name := canvas.NewText(engineDisplayName(sess.Engine), fg)
	name.TextSize = 14

	underline := canvas.NewRectangle(accent)
	underline.SetMinSize(fyne.NewSize(0, 2))

	head := container.NewHBox(dot, name)
	tab := container.NewVBox(head, underline)

	plus := widget.NewButtonWithIcon("", fynetheme.ContentAddIcon(), nil)
	plus.Importance = widget.LowImportance
	plus.Disable() // multi-tab is Step 17

	return container.NewHBox(container.NewPadded(tab), plus, layout.NewSpacer())
}
