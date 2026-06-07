package ui

import (
	"archive/zip"
	"bytes"
	"net/http"
	"strings"
	"unicode/utf8"
)

type ContentKind int

const (
	KindText ContentKind = iota
	KindImage
	KindBinary
)

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
		// xlsx/docx/pptx/odt/epub/jar all sniff as application/zip; peek inside for the real type.
		return KindBinary, refineZipMIME(v)
	case utf8.Valid(v) && !hasControlBytes(v):
		return KindText, "text/plain"
	default:
		return KindBinary, mime
	}
}

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

func hasControlBytes(v []byte) bool {
	for _, b := range v {
		if b < 0x09 || (b > 0x0d && b < 0x20) {
			return true
		}
	}
	return false
}
