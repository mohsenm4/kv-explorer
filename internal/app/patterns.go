package app

import (
	"regexp"
	"strings"
)

// Strict, anchored: only tag whole-string matches so substrings inside a larger
// blob don't get mis-labelled (a UUID inside JSON has already been handled).
var (
	reUUID  = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	reURL   = regexp.MustCompile(`(?i)^https?://[^\s]+$`)
	reEmail = regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
)

// patternTag prefixes the preview when the whole value matches a common form
// like a UUID or URL. The value itself is kept truncated to keep table rows
// scannable.
func patternTag(v []byte) string {
	s := strings.TrimSpace(string(v))
	if s == "" || len(s) > 2048 {
		return ""
	}
	switch {
	case reUUID.MatchString(s):
		return "[UUID] " + s
	case reURL.MatchString(s):
		return "[URL] " + truncate(s, 100)
	case reEmail.MatchString(s):
		return "[Email] " + s
	}
	return ""
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
