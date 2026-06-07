package app

import (
	"strings"
	"testing"
)

func TestPatternTag(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"uuid lower", "550e8400-e29b-41d4-a716-446655440000", "[UUID] 550e8400-e29b-41d4-a716-446655440000"},
		{"uuid upper", "550E8400-E29B-41D4-A716-446655440000", "[UUID] 550E8400-E29B-41D4-A716-446655440000"},
		{"uuid trimmed", "  550e8400-e29b-41d4-a716-446655440000  ", "[UUID] 550e8400-e29b-41d4-a716-446655440000"},
		{"url https", "https://example.com/api/users/42", "[URL] https://example.com/api/users/42"},
		{"url http", "http://example.com", "[URL] http://example.com"},
		{"email", "ali@example.com", "[Email] ali@example.com"},
		{"email plus", "ali+kv@sub.example.co", "[Email] ali+kv@sub.example.co"},
		{"plain text", "hello world", ""},
		{"not quite uuid (length)", "550e8400-e29b-41d4-a716-44665544", ""},
		{"not quite email", "@bad", ""},
		{"empty", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := patternTag([]byte(c.in))
			if got != c.want {
				t.Errorf("patternTag(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestPatternTag_LongURLTruncated(t *testing.T) {
	url := "https://example.com/" + strings.Repeat("a", 200)
	got := patternTag([]byte(url))
	if !strings.HasPrefix(got, "[URL] ") || !strings.HasSuffix(got, "…") {
		t.Errorf("expected truncated URL with [URL] prefix, got %q", got)
	}
}

func TestMakePreview_UsesPatternTag(t *testing.T) {
	got := makePreview([]byte("550e8400-e29b-41d4-a716-446655440000"))
	if !strings.HasPrefix(got, "[UUID] ") {
		t.Errorf("makePreview(uuid) = %q, want [UUID] prefix", got)
	}
}
