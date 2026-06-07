package i18n

import "testing"

func TestT_FallsBackToIDWhenUninitialized(t *testing.T) {
	// Reset state from any prior test.
	loc = nil
	if got := T("nonexistent.key"); got != "nonexistent.key" {
		t.Errorf("T uninitialized = %q, want id passthrough", got)
	}
}

func TestInit_LoadsEnglishByDefault(t *testing.T) {
	Init("en")
	if got := T("toolbar.open"); got != "Open" {
		t.Errorf("T(toolbar.open) en = %q, want Open", got)
	}
}

func TestSetLanguage_SwitchesCatalog(t *testing.T) {
	Init("en")
	SetLanguage("de")
	if got := T("toolbar.open"); got != "Öffnen" {
		t.Errorf("T(toolbar.open) de = %q, want Öffnen", got)
	}
	SetLanguage("fr")
	if got := T("toolbar.open"); got != "Ouvrir" {
		t.Errorf("T(toolbar.open) fr = %q, want Ouvrir", got)
	}
	SetLanguage("ja")
	if got := T("toolbar.open"); got != "開く" {
		t.Errorf("T(toolbar.open) ja = %q, want 開く", got)
	}
}

func TestTf_TemplateData(t *testing.T) {
	Init("en")
	got := Tf("status.keyCount", map[string]any{"Count": "42"})
	if got != "42 keys" {
		t.Errorf("Tf(status.keyCount, 42) = %q, want %q", got, "42 keys")
	}
}

func TestT_MissingKeyReturnsID(t *testing.T) {
	Init("en")
	if got := T("does.not.exist"); got != "does.not.exist" {
		t.Errorf("T(missing) = %q, want id passthrough", got)
	}
}

func TestSupportedTag(t *testing.T) {
	cases := []struct {
		raw  string
		want string
	}{
		{"en_US.UTF-8", "en"},
		{"de_DE.UTF-8", "de"},
		{"es_ES", "es"},
		{"fr-CA", "fr"},
		{"ja_JP.UTF-8", "ja"},
		{"zh_CN.UTF-8", "zh-Hans"},
		{"zh-TW", "zh-Hans"}, // we only ship Hans for now
		{"ru_RU", ""},
		{"", ""},
		{"en", "en"},
	}
	for _, c := range cases {
		if got := supportedTag(c.raw); got != c.want {
			t.Errorf("supportedTag(%q) = %q, want %q", c.raw, got, c.want)
		}
	}
}

func TestResolveLanguage_UnknownFallsBackToEnglish(t *testing.T) {
	// Force a clearly-unsupported locale; result should be "en" (the
	// final fallback after system detection misses too).
	got := resolveLanguage("klingon")
	if got != "en" {
		// Allow that the OS happens to be one of our supported locales —
		// in that case resolveLanguage returns the system match, which is
		// still acceptable. Only flag truly unexpected returns.
		if supportedTag(got) == "" {
			t.Errorf("resolveLanguage(klingon) = %q, want en or a supported tag", got)
		}
	}
}
