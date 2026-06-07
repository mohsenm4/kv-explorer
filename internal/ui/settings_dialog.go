package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/i18n"
)

// SettingsHandlers groups callbacks the Settings dialog can fire so it
// stays decoupled from the rest of the UI state.
type SettingsHandlers struct {
	OnTheme    func(string) // "light" | "dark" | "system"
	OnLanguage func(string) // BCP-47 tag, or "" for system default
}

// showSettings opens the tabbed Settings dialog.
func showSettings(parent fyne.Window, currentTheme, currentLang string, handlers SettingsHandlers) {
	tabs := container.NewAppTabs(
		container.NewTabItem(i18n.T("settings.tab.appearance"), pane(appearancePane(currentTheme, currentLang, handlers))),
		container.NewTabItem(i18n.T("settings.tab.general"), pane(generalPane())),
		container.NewTabItem(i18n.T("settings.tab.editor"), pane(editorPane())),
		container.NewTabItem(i18n.T("settings.tab.shortcuts"), pane(shortcutsPane())),
		container.NewTabItem(i18n.T("settings.tab.about"), pane(aboutPane())),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	d := dialog.NewCustom(i18n.T("settings.title"), i18n.T("settings.close"), tabs, parent)
	d.Resize(fyne.NewSize(720, 480))
	d.Show()
}

// pane wraps a settings tab body so its content doesn't touch the
// vertical tab strip and gets breathing room on every side.
func pane(content fyne.CanvasObject) fyne.CanvasObject {
	return container.NewPadded(container.NewPadded(content))
}

func appearancePane(currentTheme, currentLang string, h SettingsHandlers) fyne.CanvasObject {
	light := i18n.T("settings.appearance.theme.light")
	dark := i18n.T("settings.appearance.theme.dark")
	system := i18n.T("settings.appearance.theme.system")

	themeRadio := widget.NewRadioGroup(
		[]string{light, dark, system},
		func(s string) {
			if h.OnTheme == nil {
				return
			}
			switch s {
			case light:
				h.OnTheme("light")
			case dark:
				h.OnTheme("dark")
			case system:
				h.OnTheme("system")
			}
		},
	)
	switch currentTheme {
	case "dark":
		themeRadio.SetSelected(dark)
	case "system":
		themeRadio.SetSelected(system)
	default:
		themeRadio.SetSelected(light)
	}

	density := widget.NewRadioGroup([]string{
		i18n.T("settings.appearance.density.compact"),
		i18n.T("settings.appearance.density.comfortable"),
	}, nil)
	density.SetSelected(i18n.T("settings.appearance.density.comfortable"))

	zebra := widget.NewCheck(i18n.T("settings.appearance.zebra"), nil)
	mono := widget.NewCheck(i18n.T("settings.appearance.mono"), nil)

	langChoices := i18n.Available()
	labels := make([]string, len(langChoices))
	labelByCode := map[string]string{}
	codeByLabel := map[string]string{}
	for i, c := range langChoices {
		lbl := c.Label
		if c.Code == i18n.SystemTag {
			lbl = i18n.T("lang.systemDefault")
		}
		labels[i] = lbl
		labelByCode[c.Code] = lbl
		codeByLabel[lbl] = c.Code
	}
	langSelect := widget.NewSelect(labels, func(s string) {
		if h.OnLanguage == nil {
			return
		}
		h.OnLanguage(codeByLabel[s])
	})
	if lbl, ok := labelByCode[currentLang]; ok {
		langSelect.SetSelected(lbl)
	} else {
		langSelect.SetSelected(labelByCode[i18n.SystemTag])
	}

	return container.NewVBox(
		sectionLabel(i18n.T("settings.appearance.theme")),
		themeRadio,
		gap(8),
		sectionLabel(i18n.T("settings.appearance.language")),
		langSelect,
		gap(8),
		sectionLabel(i18n.T("settings.appearance.density")),
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
		sectionLabel(i18n.T("settings.general.configFile")),
		configPath,
		gap(8),
		sectionLabel(i18n.T("settings.general.logDirectory")),
		logPath,
		gap(8),
		sectionLabel(i18n.T("settings.general.logLevel")),
		level,
	)
}

func editorPane() fyne.CanvasObject {
	font := widget.NewSelect([]string{"System default", "SF Mono", "JetBrains Mono", "Consolas"}, nil)
	font.SetSelected("System default")

	size := widget.NewSelect([]string{"12", "13", "14", "16"}, nil)
	size.SetSelected("13")

	pretty := widget.NewCheck(i18n.T("settings.editor.prettyJson"), nil)
	pretty.SetChecked(true)

	return container.NewVBox(
		sectionLabel(i18n.T("settings.editor.fontFamily")),
		font,
		gap(8),
		sectionLabel(i18n.T("settings.editor.fontSize")),
		size,
		gap(8),
		pretty,
	)
}

func shortcutsPane() fyne.CanvasObject {
	rows := [][2]string{
		{i18n.T("settings.shortcut.openDatabase"), "Ctrl/Cmd + O"},
		{i18n.T("settings.shortcut.closeTab"), "Ctrl/Cmd + W"},
		{i18n.T("settings.shortcut.addKey"), "Ctrl/Cmd + N"},
		{i18n.T("settings.shortcut.focusFilter"), "Ctrl/Cmd + F"},
		{i18n.T("settings.shortcut.saveEdits"), "Ctrl/Cmd + S"},
		{i18n.T("settings.shortcut.deleteKey"), "Delete"},
		{i18n.T("settings.shortcut.editKey"), "F2"},
		{i18n.T("settings.shortcut.refresh"), "F5"},
		{i18n.T("settings.shortcut.openSettings"), "Ctrl/Cmd + ,"},
		{i18n.T("settings.shortcut.cycleTabs"), "Ctrl + Tab"},
		{i18n.T("settings.shortcut.escape"), "Esc"},
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
	title := widget.NewLabelWithStyle(i18n.T("app.name"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	version := widget.NewLabelWithStyle(
		i18n.Tf("app.version", map[string]any{"Version": "0.1.0"}),
		fyne.TextAlignCenter, fyne.TextStyle{})
	tagline := widget.NewLabelWithStyle(
		i18n.T("app.tagline"),
		fyne.TextAlignCenter, fyne.TextStyle{})
	engines := widget.NewLabelWithStyle(
		i18n.T("app.about.engines"),
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
