// Package i18n loads message catalogs for the supported UI languages and
// exposes T/Tf helpers that the UI layer calls. Catalogs live in JSON
// files under locales/ and are embedded into the binary.
//
// First-run language picks come from the system locale. The user can
// override the pick in Settings → Appearance; that choice is persisted in
// config.Config.Language.
package i18n

import (
	"embed"
	"encoding/json"
	"strings"
	"sync"

	"github.com/jeandeaual/go-locale"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localesFS embed.FS

// LangChoice is one option the user can pick in the Language dropdown.
type LangChoice struct {
	Code  string // BCP-47 tag, or "" for "system default"
	Label string // shown in the dropdown, in its own script
}

// SystemTag is the sentinel for "follow the OS locale". Stored in
// config as Language="" so an empty config also gets system-pick behavior.
const SystemTag = ""

var (
	mu       sync.RWMutex
	bundle   *goi18n.Bundle
	loc      *goi18n.Localizer
	chosen   string // raw choice ("" = system); used to round-trip the Settings dropdown
	resolved string // tag actually being served (e.g. "en", "fr")
)

// Available returns the language list shown in Settings. The first entry
// is the system-default sentinel; the rest are supported catalogs.
func Available() []LangChoice {
	return []LangChoice{
		{Code: SystemTag, Label: "System default"},
		{Code: "en", Label: "English"},
		{Code: "es", Label: "Español"},
		{Code: "de", Label: "Deutsch"},
		{Code: "fr", Label: "Français"},
		{Code: "zh-Hans", Label: "中文 (简体)"},
		{Code: "ja", Label: "日本語"},
	}
}

// Init loads every embedded catalog into a single bundle and selects the
// initial language. Called once at app start; SetLanguage handles later
// changes.
func Init(preferred string) {
	mu.Lock()
	defer mu.Unlock()

	bundle = goi18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	entries, err := localesFS.ReadDir("locales")
	if err == nil {
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			_, _ = bundle.LoadMessageFileFS(localesFS, "locales/"+e.Name())
		}
	}

	chosen = preferred
	resolved = resolveLanguage(preferred)
	loc = goi18n.NewLocalizer(bundle, resolved, "en")
}

// SetLanguage switches the active catalog. Pass SystemTag to fall back to
// the OS locale. The new tag is returned (the resolved one, not the raw
// choice) so callers can persist it.
func SetLanguage(preferred string) string {
	mu.Lock()
	defer mu.Unlock()
	if bundle == nil {
		return resolved
	}
	chosen = preferred
	resolved = resolveLanguage(preferred)
	loc = goi18n.NewLocalizer(bundle, resolved, "en")
	return resolved
}

// Chosen returns the raw preference ("" for system default, or a tag).
// Used by Settings to round-trip the dropdown selection.
func Chosen() string {
	mu.RLock()
	defer mu.RUnlock()
	return chosen
}

// Current returns the tag currently serving translations — never "".
func Current() string {
	mu.RLock()
	defer mu.RUnlock()
	return resolved
}

// T returns the translated message for id. Falls back to id itself if the
// catalog has no entry — that way missing keys are visible during dev.
func T(id string) string {
	mu.RLock()
	l := loc
	mu.RUnlock()
	if l == nil {
		return id
	}
	s, err := l.Localize(&goi18n.LocalizeConfig{MessageID: id})
	if err != nil {
		return id
	}
	return s
}

// Tf is T with template-data substitution. Use go-i18n template syntax in
// catalogs: `{{.Key}}`, `{{.Size}}`, etc.
func Tf(id string, data map[string]any) string {
	mu.RLock()
	l := loc
	mu.RUnlock()
	if l == nil {
		return id
	}
	s, err := l.Localize(&goi18n.LocalizeConfig{MessageID: id, TemplateData: data})
	if err != nil {
		return id
	}
	return s
}

// resolveLanguage maps a raw preference to a supported tag. Empty / "auto"
// reads the OS locale; an unknown tag falls back to English.
func resolveLanguage(preferred string) string {
	if preferred != SystemTag && preferred != "auto" {
		if supportedTag(preferred) != "" {
			return supportedTag(preferred)
		}
	}
	sysTag := ""
	if s, err := locale.GetLocale(); err == nil {
		sysTag = s
	}
	if m := supportedTag(sysTag); m != "" {
		return m
	}
	return "en"
}

// supportedTag matches a raw locale string (e.g. "fr_FR.UTF-8", "zh-CN",
// "ja_JP") against the catalogs we ship. Returns the canonical tag we use
// internally, or "" if no match.
func supportedTag(raw string) string {
	if raw == "" {
		return ""
	}
	// Normalize: strip charset, swap underscore for dash, lowercase the
	// language subtag.
	s := raw
	if i := strings.IndexByte(s, '.'); i > 0 {
		s = s[:i]
	}
	s = strings.ReplaceAll(s, "_", "-")
	lower := strings.ToLower(s)

	// Direct exact matches first.
	for _, c := range Available() {
		if c.Code == "" {
			continue
		}
		if strings.EqualFold(c.Code, s) {
			return c.Code
		}
	}

	// Chinese: any zh-* with Hans/CN/SG/MY maps to zh-Hans; Hant/TW/HK
	// would map to zh-Hant if we add it later. For now anything Chinese
	// goes to Simplified.
	if strings.HasPrefix(lower, "zh") {
		return "zh-Hans"
	}

	// Otherwise compare the language subtag (first 2 chars).
	prim := lower
	if i := strings.IndexByte(lower, '-'); i > 0 {
		prim = lower[:i]
	}
	switch prim {
	case "en":
		return "en"
	case "es":
		return "es"
	case "de":
		return "de"
	case "fr":
		return "fr"
	case "ja":
		return "ja"
	}
	return ""
}
