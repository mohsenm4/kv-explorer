package ui

import (
	"archive/zip"
	"bytes"
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
	case mime == "application/zip":
		// ZIP is a container — xlsx/docx/pptx/odt/epub/jar all sniff as
		// application/zip. Peek inside to return the real type.
		return KindBinary, refineZipMIME(v)
	case utf8.Valid(v) && !hasControlBytes(v):
		return KindText, "text/plain"
	default:
		return KindBinary, mime
	}
}

// refineZipMIME inspects a ZIP archive's entries and returns a more specific
// MIME type for known ZIP-based container formats (OOXML, ODF, EPUB, JAR).
// Returns "application/zip" if the archive is plain or unreadable.
func refineZipMIME(v []byte) string {
	r, err := zip.NewReader(bytes.NewReader(v), int64(len(v)))
	if err != nil {
		return "application/zip"
	}
	// EPUB / ODF put a "mimetype" file with the exact MIME as its content.
	for _, f := range r.File {
		if f.Name != "mimetype" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			break
		}
		buf := make([]byte, 128)
		n, _ := rc.Read(buf)
		rc.Close()
		s := strings.TrimSpace(string(buf[:n]))
		if s != "" {
			return s
		}
		break
	}
	// OOXML (Office) and JAR — recognized by directory prefixes.
	for _, f := range r.File {
		switch {
		case strings.HasPrefix(f.Name, "xl/"):
			return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		case strings.HasPrefix(f.Name, "word/"):
			return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		case strings.HasPrefix(f.Name, "ppt/"):
			return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
		case f.Name == "META-INF/MANIFEST.MF":
			return "application/java-archive"
		}
	}
	return "application/zip"
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
