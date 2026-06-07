package ui

import (
	"strings"
	"testing"

	"github.com/mohsenm4/kv-explorer/internal/i18n"
)

func TestTimestampHint(t *testing.T) {
	i18n.Init("en")
	cases := []struct {
		name    string
		in      string
		wantSub string
	}{
		{"seconds 2024", "1717804800", "2024-"},
		{"milliseconds 2024", "1717804800000", "2024-"},
		{"microseconds 2024", "1717804800000000", "2024-"},
		{"nanoseconds 2024", "1717804800000000000", "2024-"},
		{"empty", "", ""},
		{"not numeric", "hello", ""},
		{"too short", "12345", ""},
		{"out of range tiny", "1", ""},
		{"out of range giant", "99999999999999999999", ""},
		{"with whitespace", "  1717804800  ", "2024-"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := timestampHint([]byte(c.in))
			if c.wantSub == "" {
				if got != "" {
					t.Errorf("timestampHint(%q) = %q, want empty", c.in, got)
				}
				return
			}
			if !strings.Contains(got, c.wantSub) {
				t.Errorf("timestampHint(%q) = %q, want substring %q", c.in, got, c.wantSub)
			}
		})
	}
}
