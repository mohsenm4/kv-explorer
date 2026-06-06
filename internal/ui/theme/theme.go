package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	fynetheme "fyne.io/fyne/v2/theme"
)

type Theme struct{}

func New() fyne.Theme { return Theme{} }

// ForcedVariant wraps a base theme and forces every color lookup to use the
// given variant, ignoring whatever variant Fyne would otherwise pass in.
func ForcedVariant(base fyne.Theme, v fyne.ThemeVariant) fyne.Theme {
	return &forcedVariantTheme{base: base, variant: v}
}

type forcedVariantTheme struct {
	base    fyne.Theme
	variant fyne.ThemeVariant
}

func (t *forcedVariantTheme) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return t.base.Color(n, t.variant)
}
func (t *forcedVariantTheme) Font(s fyne.TextStyle) fyne.Resource  { return t.base.Font(s) }
func (t *forcedVariantTheme) Icon(n fyne.ThemeIconName) fyne.Resource { return t.base.Icon(n) }
func (t *forcedVariantTheme) Size(n fyne.ThemeSizeName) float32     { return t.base.Size(n) }

func (Theme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if variant == fynetheme.VariantDark {
		if c, ok := darkPalette[name]; ok {
			return c
		}
	} else {
		if c, ok := lightPalette[name]; ok {
			return c
		}
	}
	return fynetheme.DefaultTheme().Color(name, variant)
}

func (Theme) Font(style fyne.TextStyle) fyne.Resource {
	return fynetheme.DefaultTheme().Font(style)
}

func (Theme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return fynetheme.DefaultTheme().Icon(name)
}

func (Theme) Size(name fyne.ThemeSizeName) float32 {
	if s, ok := sizes[name]; ok {
		return s
	}
	return fynetheme.DefaultTheme().Size(name)
}

var sizes = map[fyne.ThemeSizeName]float32{
	fynetheme.SizeNamePadding:            8,
	fynetheme.SizeNameInnerPadding:       4,
	fynetheme.SizeNameInlineIcon:         16,
	fynetheme.SizeNameSeparatorThickness: 1,
	fynetheme.SizeNameInputBorder:        1,
	fynetheme.SizeNameText:               14,
	fynetheme.SizeNameHeadingText:        18,
	fynetheme.SizeNameSubHeadingText:     14,
	fynetheme.SizeNameCaptionText:        11,
}

func rgb(r, g, b uint8) color.NRGBA     { return color.NRGBA{R: r, G: g, B: b, A: 0xFF} }
func rgba(r, g, b, a uint8) color.NRGBA { return color.NRGBA{R: r, G: g, B: b, A: a} }

// DBAccentTint returns DBAccent at ~20% opacity, for chip backgrounds.
func DBAccentTint(engine string, v fyne.ThemeVariant) color.Color {
	c, ok := DBAccent(engine, v).(color.NRGBA)
	if !ok {
		return DBAccent(engine, v)
	}
	c.A = 0x33
	return c
}

// DBAccent returns the accent color for a database engine, per spec §2.3.
// The `engine` string should match the kvstore.EngineKind values
// ("pebble", "badger", "leveldb"). Unknown engines fall back to placeholder.
func DBAccent(engine string, v fyne.ThemeVariant) color.Color {
	dark := v == fynetheme.VariantDark
	switch engine {
	case "pebble":
		if dark {
			return rgb(0x38, 0xBD, 0xF8)
		}
		return rgb(0x0E, 0xA5, 0xE9)
	case "badger":
		if dark {
			return rgb(0xFB, 0xBF, 0x24)
		}
		return rgb(0xF5, 0x9E, 0x0B)
	case "leveldb":
		if dark {
			return rgb(0x34, 0xD3, 0x99)
		}
		return rgb(0x10, 0xB9, 0x81)
	}
	return rgb(0x94, 0xA3, 0xB8)
}

var lightPalette = map[fyne.ThemeColorName]color.Color{
	fynetheme.ColorNameBackground:          rgb(0xFA, 0xFA, 0xFA),
	fynetheme.ColorNameForeground:          rgb(0x0F, 0x17, 0x2A),
	fynetheme.ColorNamePlaceHolder:         rgb(0x64, 0x74, 0x8B),
	fynetheme.ColorNameDisabled:            rgb(0xCB, 0xD5, 0xE1),
	fynetheme.ColorNameInputBackground:     rgb(0xFF, 0xFF, 0xFF),
	fynetheme.ColorNameInputBorder:         rgb(0xE2, 0xE8, 0xF0),
	fynetheme.ColorNameSeparator:           rgb(0xE2, 0xE8, 0xF0),
	fynetheme.ColorNameHeaderBackground:    rgb(0xF3, 0xF4, 0xF6),
	fynetheme.ColorNameMenuBackground:      rgb(0xFF, 0xFF, 0xFF),
	fynetheme.ColorNameOverlayBackground:   rgb(0xFF, 0xFF, 0xFF),
	fynetheme.ColorNamePrimary:             rgb(0x4F, 0x46, 0xE5),
	fynetheme.ColorNameForegroundOnPrimary: rgb(0xFF, 0xFF, 0xFF),
	fynetheme.ColorNameFocus:               rgb(0x4F, 0x46, 0xE5),
	fynetheme.ColorNameSelection:           rgb(0xEE, 0xF2, 0xFF),
	fynetheme.ColorNameHover:               rgba(0x00, 0x00, 0x00, 0x0A),
	fynetheme.ColorNamePressed:             rgba(0x00, 0x00, 0x00, 0x14),
	fynetheme.ColorNameSuccess:             rgb(0x10, 0xB9, 0x81),
	fynetheme.ColorNameForegroundOnSuccess: rgb(0xFF, 0xFF, 0xFF),
	fynetheme.ColorNameWarning:             rgb(0xF5, 0x9E, 0x0B),
	fynetheme.ColorNameForegroundOnWarning: rgb(0xFF, 0xFF, 0xFF),
	fynetheme.ColorNameError:               rgb(0xEF, 0x44, 0x44),
	fynetheme.ColorNameForegroundOnError:   rgb(0xFF, 0xFF, 0xFF),
	fynetheme.ColorNameScrollBar:           rgb(0x94, 0xA3, 0xB8),
	fynetheme.ColorNameShadow:              rgba(0x0F, 0x17, 0x2A, 0x14),
}

var darkPalette = map[fyne.ThemeColorName]color.Color{
	fynetheme.ColorNameBackground:          rgb(0x0B, 0x12, 0x20),
	fynetheme.ColorNameForeground:          rgb(0xF1, 0xF5, 0xF9),
	fynetheme.ColorNamePlaceHolder:         rgb(0x94, 0xA3, 0xB8),
	fynetheme.ColorNameDisabled:            rgb(0x47, 0x55, 0x69),
	fynetheme.ColorNameInputBackground:     rgb(0x11, 0x18, 0x27),
	fynetheme.ColorNameInputBorder:         rgb(0x1F, 0x29, 0x37),
	fynetheme.ColorNameSeparator:           rgb(0x1F, 0x29, 0x37),
	fynetheme.ColorNameHeaderBackground:    rgb(0x1F, 0x29, 0x37),
	fynetheme.ColorNameMenuBackground:      rgb(0x11, 0x18, 0x27),
	fynetheme.ColorNameOverlayBackground:   rgb(0x11, 0x18, 0x27),
	fynetheme.ColorNamePrimary:             rgb(0x81, 0x8C, 0xF8),
	fynetheme.ColorNameForegroundOnPrimary: rgb(0x0B, 0x12, 0x20),
	fynetheme.ColorNameFocus:               rgb(0x81, 0x8C, 0xF8),
	fynetheme.ColorNameSelection:           rgb(0x1E, 0x1B, 0x4B),
	fynetheme.ColorNameHover:               rgba(0xFF, 0xFF, 0xFF, 0x0F),
	fynetheme.ColorNamePressed:             rgba(0xFF, 0xFF, 0xFF, 0x1A),
	fynetheme.ColorNameSuccess:             rgb(0x34, 0xD3, 0x99),
	fynetheme.ColorNameForegroundOnSuccess: rgb(0x0B, 0x12, 0x20),
	fynetheme.ColorNameWarning:             rgb(0xFB, 0xBF, 0x24),
	fynetheme.ColorNameForegroundOnWarning: rgb(0x0B, 0x12, 0x20),
	fynetheme.ColorNameError:               rgb(0xF8, 0x71, 0x71),
	fynetheme.ColorNameForegroundOnError:   rgb(0xFF, 0xFF, 0xFF),
	fynetheme.ColorNameScrollBar:           rgb(0x47, 0x55, 0x69),
	fynetheme.ColorNameShadow:              rgba(0x00, 0x00, 0x00, 0x66),
}
