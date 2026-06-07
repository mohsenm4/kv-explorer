package ui

import (
	"strconv"
	"strings"
	"time"

	"github.com/mohsenm4/kv-explorer/internal/i18n"
)

// timestampHint returns "Unix timestamp: <iso8601 UTC>" when v is a bare integer
// in a unit (s/ms/µs/ns) that lands in the 2001–2100 range. Empty string means
// "this isn't a timestamp" so the caller can skip rendering.
func timestampHint(v []byte) string {
	s := strings.TrimSpace(string(v))
	if s == "" {
		return ""
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil || n <= 0 {
		return ""
	}

	t, ok := decodeUnixTimestamp(n, len(s))
	if !ok {
		return ""
	}
	return i18n.Tf("editor.timestampHint", map[string]any{
		"Time": t.UTC().Format("2006-01-02 15:04:05 MST"),
	})
}

// Decode is driven by digit count so 1717804800 is read as seconds and the same
// number with three zeros tacked on is read as milliseconds — matches how these
// values usually appear in KV stores. The sanity range filters out random ints.
func decodeUnixTimestamp(n int64, digits int) (time.Time, bool) {
	var t time.Time
	switch {
	case digits == 10:
		t = time.Unix(n, 0)
	case digits == 13:
		t = time.UnixMilli(n)
	case digits == 16:
		t = time.UnixMicro(n)
	case digits == 19:
		t = time.Unix(0, n)
	default:
		return time.Time{}, false
	}
	if t.Year() < 2001 || t.Year() > 2100 {
		return time.Time{}, false
	}
	return t, true
}
