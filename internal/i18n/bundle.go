// Package i18n loads embedded JSON catalogs and exposes T/Tf helpers.
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

type LangChoice struct {
	Code  string
	Label string
}

// SystemTag means follow the OS locale; stored as Language="" in config.
const SystemTag = ""

var (
	mu       sync.RWMutex
	bundle   *goi18n.Bundle
	loc      *goi18n.Localizer
	chosen   string
	resolved string
)

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

func Chosen() string {
	mu.RLock()
	defer mu.RUnlock()
	return chosen
}

func Current() string {
	mu.RLock()
	defer mu.RUnlock()
	return resolved
}

// T returns id itself when a key is missing, so dev sees the gap.
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

func supportedTag(raw string) string {
	if raw == "" {
		return ""
	}
	s := raw
	if i := strings.IndexByte(s, '.'); i > 0 {
		s = s[:i]
	}
	s = strings.ReplaceAll(s, "_", "-")
	lower := strings.ToLower(s)

	for _, c := range Available() {
		if c.Code == "" {
			continue
		}
		if strings.EqualFold(c.Code, s) {
			return c.Code
		}
	}

	// Any zh-* collapses to Simplified until we ship a Traditional catalog.
	if strings.HasPrefix(lower, "zh") {
		return "zh-Hans"
	}

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
