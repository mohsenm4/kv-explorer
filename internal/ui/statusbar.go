package ui

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
	"github.com/mohsenm4/kv-explorer/internal/i18n"
	apptheme "github.com/mohsenm4/kv-explorer/internal/ui/theme"
)

// welcomeStatusBar shows "No database open" with the theme toggle on the right.
func welcomeStatusBar(v fyne.ThemeVariant, onToggle func()) fyne.CanvasObject {
	muted := themeColor(v, fynetheme.ColorNamePlaceHolder)
	msg := caption(i18n.T("status.noDatabase"), muted)
	return statusBarShell(v, container.NewPadded(msg), themeToggle(v, onToggle))
}

// mainStatusBar shows engine dot + name + key count + size + path on the
// left, theme toggle on the right.
func mainStatusBar(v fyne.ThemeVariant, sess *app.Session, onToggle func()) fyne.CanvasObject {
	muted := themeColor(v, fynetheme.ColorNamePlaceHolder)
	accent := apptheme.DBAccent(string(sess.Engine), v)

	dot := engineDot(string(sess.Engine), v, 11)

	name := canvas.NewText(strings.ToUpper(string(sess.Engine)), accent)
	name.TextSize = 11
	name.TextStyle = fyne.TextStyle{Bold: true}

	keys := caption(i18n.Tf("status.keyCount", map[string]any{"Count": thousands(sess.KeyCount)}), muted)
	size := caption(humanSize(sess.SizeBytes), muted)
	path := caption(displayPath(sess.Path), muted)

	left := container.NewHBox(
		dot,
		name,
		caption("  ", muted),
		keys,
		caption("  ", muted),
		size,
		caption("  ", muted),
		path,
	)

	return statusBarShell(v, container.NewPadded(left), themeToggle(v, onToggle))
}

func statusBarShell(v fyne.ThemeVariant, leftContent, rightContent fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(themeColor(v, fynetheme.ColorNameHeaderBackground))
	row := container.NewBorder(nil, nil, leftContent, container.NewPadded(rightContent), layout.NewSpacer())
	return container.NewStack(bg, row)
}

func themeToggle(v fyne.ThemeVariant, onToggle func()) fyne.CanvasObject {
	label := i18n.T("status.theme.light")
	if v == fynetheme.VariantDark {
		label = i18n.T("status.theme.dark")
	}
	b := widget.NewButton(label, onToggle)
	b.Importance = widget.MediumImportance
	return b
}

func caption(text string, c color.Color) fyne.CanvasObject {
	t := canvas.NewText(text, c)
	t.TextSize = 11
	return t
}

func themeColor(v fyne.ThemeVariant, name fyne.ThemeColorName) color.Color {
	return fyne.CurrentApp().Settings().Theme().Color(name, v)
}

func humanSize(b int64) string {
	switch {
	case b < 1024:
		return fmt.Sprintf("%d B", b)
	case b < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(b)/1024)
	case b < 1024*1024*1024:
		return fmt.Sprintf("%.1f MB", float64(b)/(1024*1024))
	default:
		return fmt.Sprintf("%.1f GB", float64(b)/(1024*1024*1024))
	}
}

func thousands(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var b strings.Builder
	pre := len(s) % 3
	if pre > 0 {
		b.WriteString(s[:pre])
		if len(s) > pre {
			b.WriteByte(',')
		}
	}
	for i := pre; i < len(s); i += 3 {
		b.WriteString(s[i : i+3])
		if i+3 < len(s) {
			b.WriteByte(',')
		}
	}
	return b.String()
}

func displayPath(p string) string {
	home, err := os.UserHomeDir()
	if err == nil && strings.HasPrefix(p, home) {
		return "~" + strings.TrimPrefix(p, home)
	}
	return p
}
