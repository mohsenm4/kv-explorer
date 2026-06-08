package app

import (
	"strings"
	"testing"
)

func TestPatternKind(t *testing.T) {
	cases := []struct {
		name        string
		in          string
		wantKind    string
		wantPreview string
	}{
		{"uuid lower", "550e8400-e29b-41d4-a716-446655440000", "UUID", "550e8400-e29b-41d4-a716-446655440000"},
		{"uuid upper", "550E8400-E29B-41D4-A716-446655440000", "UUID", "550E8400-E29B-41D4-A716-446655440000"},
		{"uuid trimmed", "  550e8400-e29b-41d4-a716-446655440000  ", "UUID", "550e8400-e29b-41d4-a716-446655440000"},
		{"url https", "https://example.com/api/users/42", "URL", "https://example.com/api/users/42"},
		{"url http", "http://example.com", "URL", "http://example.com"},
		{"email", "ali@example.com", "Email", "ali@example.com"},
		{"email plus", "ali+kv@sub.example.co", "Email", "ali+kv@sub.example.co"},
		{"plain text", "hello world", "", ""},
		{"not quite uuid (length)", "550e8400-e29b-41d4-a716-44665544", "", ""},
		{"not quite email", "@bad", "", ""},
		{"empty", "", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			kind, preview := patternKind([]byte(c.in))
			if kind != c.wantKind || preview != c.wantPreview {
				t.Errorf("patternKind(%q) = (%q, %q), want (%q, %q)", c.in, kind, preview, c.wantKind, c.wantPreview)
			}
		})
	}
}

func TestPatternKind_LongURLTruncated(t *testing.T) {
	url := "https://example.com/" + strings.Repeat("a", 200)
	kind, preview := patternKind([]byte(url))
	if kind != "URL" {
		t.Errorf("expected URL kind, got %q", kind)
	}
	if !strings.HasSuffix(preview, "…") {
		t.Errorf("expected truncated preview, got %q", preview)
	}
}

func TestMakePreview_UsesPatternKind(t *testing.T) {
	preview, kind := makePreview([]byte("550e8400-e29b-41d4-a716-446655440000"))
	if kind != "UUID" {
		t.Errorf("makePreview(uuid) kind = %q, want UUID", kind)
	}
	if strings.HasPrefix(preview, "[") {
		t.Errorf("preview should not embed a bracketed prefix, got %q", preview)
	}
}
