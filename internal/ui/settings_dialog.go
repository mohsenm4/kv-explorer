package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// SettingsHandlers groups callbacks the Settings dialog can fire so it
// stays decoupled from the rest of the UI state.
type SettingsHandlers struct {
	OnTheme func(string) // "light" | "dark" | "system"
}

// showSettings opens the tabbed Settings dialog.
func showSettings(parent fyne.Window, current string, handlers SettingsHandlers) {
	tabs := container.NewAppTabs(
		container.NewTabItem("Appearance", pane(appearancePane(current, handlers))),
		container.NewTabItem("General", pane(generalPane())),
		container.NewTabItem("Editor", pane(editorPane())),
		container.NewTabItem("Shortcuts", pane(shortcutsPane())),
		container.NewTabItem("About", pane(aboutPane())),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	d := dialog.NewCustom("Settings", "Close", tabs, parent)
	d.Resize(fyne.NewSize(720, 480))
	d.Show()
}

// pane wraps a settings tab body so its content doesn't touch the
// vertical tab strip and gets breathing room on every side.
func pane(content fyne.CanvasObject) fyne.CanvasObject {
	return container.NewPadded(container.NewPadded(content))
}

func appearancePane(current string, h SettingsHandlers) fyne.CanvasObject {
	themeRadio := widget.NewRadioGroup(
		[]string{"Light", "Dark", "Follow system"},
		func(s string) {
			if h.OnTheme == nil {
				return
			}
			switch s {
			case "Light":
				h.OnTheme("light")
			case "Dark":
				h.OnTheme("dark")
			case "Follow system":
				h.OnTheme("system")
			}
		},
	)
	switch current {
	case "dark":
		themeRadio.SetSelected("Dark")
	case "system":
		themeRadio.SetSelected("Follow system")
	default:
		themeRadio.SetSelected("Light")
	}

	density := widget.NewRadioGroup([]string{"Compact", "Comfortable"}, nil)
	density.SetSelected("Comfortable")

	zebra := widget.NewCheck("Show zebra rows", nil)
	mono := widget.NewCheck("Use monospace everywhere", nil)

	return container.NewVBox(
		sectionLabel("Theme"),
		themeRadio,
		gap(8),
		sectionLabel("Density"),
		density,
		gap(8),
		zebra,
		mono,
	)
}

func generalPane() fyne.CanvasObject {
	configPath := widget.NewEntry()
	configPath.SetText("~/.kvexplorer/config.json")
	configPath.Disable()

	logPath := widget.NewEntry()
	logPath.SetText("~/.kvexplorer/logs/")
	logPath.Disable()

	level := widget.NewSelect([]string{"debug", "info", "warn", "error"}, nil)
	level.SetSelected("info")

	return container.NewVBox(
		sectionLabel("Config file"),
		configPath,
		gap(8),
		sectionLabel("Log directory"),
		logPath,
		gap(8),
		sectionLabel("Log level"),
		level,
	)
}

func editorPane() fyne.CanvasObject {
	font := widget.NewSelect([]string{"System default", "SF Mono", "JetBrains Mono", "Consolas"}, nil)
	font.SetSelected("System default")

	size := widget.NewSelect([]string{"12", "13", "14", "16"}, nil)
	size.SetSelected("13")

	pretty := widget.NewCheck("Pretty-print JSON when displaying", nil)
	pretty.SetChecked(true)

	return container.NewVBox(
		sectionLabel("Font family"),
		font,
		gap(8),
		sectionLabel("Font size"),
		size,
		gap(8),
		pretty,
	)
}

func shortcutsPane() fyne.CanvasObject {
	rows := [][2]string{
		{"Open Database…", "Ctrl/Cmd + O"},
		{"Close current tab", "Ctrl/Cmd + W"},
		{"Add key", "Ctrl/Cmd + N"},
		{"Focus filter", "Ctrl/Cmd + F"},
		{"Save value edits", "Ctrl/Cmd + S"},
		{"Delete selected key", "Delete"},
		{"Edit selected key", "F2"},
		{"Refresh", "F5"},
		{"Open Settings", "Ctrl/Cmd + ,"},
		{"Cycle tabs", "Ctrl + Tab"},
		{"Close dialog / clear filter / cancel edit", "Esc"},
	}
	list := container.NewVBox()
	for _, r := range rows {
		action := widget.NewLabel(r[0])
		shortcut := widget.NewLabel(r[1])
		shortcut.TextStyle = fyne.TextStyle{Monospace: true}
		row := container.NewBorder(nil, nil, action, shortcut, nil)
		list.Add(row)
	}
	return container.NewVScroll(list)
}

func aboutPane() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("KV-Explorer", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	version := widget.NewLabelWithStyle("Version 0.1.0", fyne.TextAlignCenter, fyne.TextStyle{})
	tagline := widget.NewLabelWithStyle(
		"Inspect, edit, and compare key-value databases.",
		fyne.TextAlignCenter, fyne.TextStyle{})
	engines := widget.NewLabelWithStyle(
		"PebbleDB · BadgerDB · LevelDB",
		fyne.TextAlignCenter, fyne.TextStyle{})

	icon := widget.NewIcon(fynetheme.StorageIcon())
	iconBox := container.NewGridWrap(fyne.NewSize(64, 64), icon)

	return container.NewVBox(
		gap(16),
		container.NewCenter(iconBox),
		title,
		version,
		gap(8),
		tagline,
		engines,
	)
}
