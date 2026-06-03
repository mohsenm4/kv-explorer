package ui

import (
	"net/http"
	"strings"
	"unicode/utf8"
)

// ContentKind is the broad category we use to pick an editor body.
type ContentKind int

const (
	KindText ContentKind = iota
	KindImage
	KindBinary
)

// DetectContent inspects value bytes and returns its kind plus a best-effort
// MIME type. It uses http.DetectContentType for magic-byte sniffing and
// falls back to a UTF-8 + control-byte check to catch generic text.
func DetectContent(v []byte) (ContentKind, string) {
	if len(v) == 0 {
		return KindText, "text/plain"
	}
	mime := http.DetectContentType(v)
	switch {
	case strings.HasPrefix(mime, "image/"):
		return KindImage, mime
	case strings.HasPrefix(mime, "text/"):
		return KindText, mime
	case mime == "application/json":
		return KindText, mime
	case utf8.Valid(v) && !hasControlBytes(v):
		return KindText, "text/plain"
	default:
		return KindBinary, mime
	}
}

// hasControlBytes reports whether the slice contains any C0 control byte
// other than the standard whitespace ones (tab, LF, CR).
func hasControlBytes(v []byte) bool {
	for _, b := range v {
		if b < 0x09 || (b > 0x0d && b < 0x20) {
			return true
		}
	}
	return false
}
