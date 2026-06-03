package ui

import "testing"

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
	// Arbitrary bytes with embedded NUL: not text, not image, not known.
	b := []byte{0x00, 0x01, 0x02, 0xff, 0xfe, 0x7f}
	if k, _ := DetectContent(b); k != KindBinary {
		t.Errorf("got %v, want KindBinary", k)
	}
}
