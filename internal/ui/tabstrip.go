package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
	apptheme "github.com/mohsenm4/kv-explorer/internal/ui/theme"
)

type TabBar struct {
	Sessions []*app.Session
	Active   int
	OnSelect func(int)
	OnClose  func(int)
	OnAdd    func()
}

func tabStrip(v fyne.ThemeVariant, bar TabBar) fyne.CanvasObject {
	row := container.NewHBox()
	for i, s := range bar.Sessions {
		idx := i
		row.Add(buildTab(v, s, i == bar.Active,
			func() { bar.OnSelect(idx) },
			func() { bar.OnClose(idx) },
		))
	}

	plus := widget.NewButtonWithIcon("", fynetheme.ContentAddIcon(), bar.OnAdd)
	plus.Importance = widget.LowImportance
	row.Add(plus)
	row.Add(layout.NewSpacer())
	return row
}

func buildTab(v fyne.ThemeVariant, sess *app.Session, active bool, onSelect, onClose func()) fyne.CanvasObject {
	accent := apptheme.DBAccent(string(sess.Engine), v)
	fg := themeColor(v, fynetheme.ColorNameForeground)
	muted := themeColor(v, fynetheme.ColorNamePlaceHolder)

	dot := engineDot(string(sess.Engine), v, 14)

	nameColor := fg
	if !active {
		nameColor = muted
	}
	name := canvas.NewText(engineDisplayName(sess.Engine), nameColor)
	name.TextSize = 14

	closeBtn := widget.NewButtonWithIcon("", fynetheme.CancelIcon(), onClose)
	closeBtn.Importance = widget.LowImportance

	head := container.NewHBox(dot, name, closeBtn)

	var underline *canvas.Rectangle
	if active {
		underline = canvas.NewRectangle(accent)
	} else {
		underline = canvas.NewRectangle(color.Transparent)
	}
	underline.SetMinSize(fyne.NewSize(0, 2))

	return newTappable(container.NewVBox(container.NewPadded(head), underline), onSelect)
}
