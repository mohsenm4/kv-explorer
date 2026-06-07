package ui

import (
	"archive/zip"
	"bytes"
	"testing"
)

func makeZip(t *testing.T, entries map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for name, body := range entries {
		f, err := w.Create(name)
		if err != nil {
			t.Fatalf("zip create %s: %v", name, err)
		}
		if _, err := f.Write([]byte(body)); err != nil {
			t.Fatalf("zip write %s: %v", name, err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("zip close: %v", err)
	}
	return buf.Bytes()
}

func TestDetectContent_ZipContainers(t *testing.T) {
	cases := []struct {
		name    string
		entries map[string]string
		want    string
	}{
		{"xlsx", map[string]string{"[Content_Types].xml": "<x/>", "xl/workbook.xml": "<x/>"},
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
		{"docx", map[string]string{"[Content_Types].xml": "<x/>", "word/document.xml": "<x/>"},
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
		{"pptx", map[string]string{"[Content_Types].xml": "<x/>", "ppt/presentation.xml": "<x/>"},
			"application/vnd.openxmlformats-officedocument.presentationml.presentation"},
		{"epub", map[string]string{"mimetype": "application/epub+zip", "OEBPS/content.opf": "<x/>"},
			"application/epub+zip"},
		{"plain zip", map[string]string{"readme.txt": "hi"}, "application/zip"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, mime := DetectContent(makeZip(t, c.entries))
			if mime != c.want {
				t.Errorf("DetectContent(%s) mime = %q, want %q", c.name, mime, c.want)
			}
		})
	}
}

func TestDetectContent_Text(t *testing.T) {
	cases := [][]byte{
		[]byte("hello world"),
		[]byte(`{"foo":"bar"}`),
		nil,
		[]byte("line1\nline2\ttabbed"),
	}
	for _, c := range cases {
		if k, _ := DetectContent(c); k != KindText {
			t.Errorf("DetectContent(%q) kind = %v, want KindText", c, k)
		}
	}
}

func TestDetectContent_Image(t *testing.T) {
	pngHead := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0, 0, 0, 0}
	jpgHead := []byte{0xff, 0xd8, 0xff, 0xe0, 0, 0, 0, 0}
	gifHead := []byte("GIF89a")
	for _, c := range [][]byte{pngHead, jpgHead, gifHead} {
		k, mime := DetectContent(c)
		if k != KindImage {
			t.Errorf("DetectContent(%v) kind = %v, want KindImage (mime=%s)", c[:4], k, mime)
		}
	}
}

func TestDetectContent_Binary(t *testing.T) {
	b := []byte{0x00, 0x01, 0x02, 0xff, 0xfe, 0x7f}
	if k, _ := DetectContent(b); k != KindBinary {
		t.Errorf("got %v, want KindBinary", k)
	}
}
