package app

import (
	"strings"
	"testing"
)

func TestMakePreview_JSONObject(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty object", `{}`, "{}"},
		{"single field", `{"id":42}`, "{id}"},
		{"three fields", `{"id":1,"name":"a","email":"a@b"}`, "{id, name, email}"},
		{"more than three", `{"a":1,"b":2,"c":3,"d":4,"e":5}`, "{a, b, c, +2}"},
		{"nested values irrelevant", `{"user":{"x":1},"items":[1,2,3]}`, "{user, items}"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			preview, kind := makePreview([]byte(c.in))
			if preview != c.want {
				t.Errorf("makePreview(%q) preview = %q, want %q", c.in, preview, c.want)
			}
			if kind != "JSON" {
				t.Errorf("makePreview(%q) kind = %q, want JSON", c.in, kind)
			}
		})
	}
}

func TestMakePreview_JSONArray(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{`[]`, "[]"},
		{`[1]`, "[1 item]"},
		{`[1,2,3,4,5]`, "[5 items]"},
		{`[{"a":1},{"b":2}]`, "[2 items]"},
	}
	for _, c := range cases {
		preview, kind := makePreview([]byte(c.in))
		if preview != c.want {
			t.Errorf("makePreview(%q) preview = %q, want %q", c.in, preview, c.want)
		}
		if kind != "JSON" {
			t.Errorf("makePreview(%q) kind = %q, want JSON", c.in, kind)
		}
	}
}

func TestMakePreview_NonJSONFallback(t *testing.T) {
	// JSON-ish but invalid: must fall back to raw-text preview, not "".
	preview, kind := makePreview([]byte("{not json at all"))
	if preview == "" {
		t.Errorf("makePreview fell through to empty for non-JSON text")
	}
	if kind != "TXT" {
		t.Errorf("makePreview kind = %q, want TXT", kind)
	}
}

func TestMakePreview_PlainText(t *testing.T) {
	preview, kind := makePreview([]byte("hello world"))
	if preview != "hello world" {
		t.Errorf("makePreview(plain) preview = %q, want %q", preview, "hello world")
	}
	if kind != "TXT" {
		t.Errorf("makePreview(plain) kind = %q, want TXT", kind)
	}
}

func TestMakePreview_Binary(t *testing.T) {
	preview, kind := makePreview([]byte{0x00, 0x01, 0x02, 0xff})
	if preview == "" {
		t.Errorf("makePreview(binary) preview empty")
	}
	if kind != "BIN" {
		t.Errorf("makePreview(binary) kind = %q, want BIN", kind)
	}
	if strings.HasPrefix(preview, "[") {
		t.Errorf("preview should not embed a bracketed prefix, got %q", preview)
	}
}

func TestMakePreview_Image(t *testing.T) {
	pngHead := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0, 0, 0, 0}
	preview, kind := makePreview(pngHead)
	if kind != "IMG" {
		t.Errorf("makePreview(png) kind = %q, want IMG", kind)
	}
	if !strings.Contains(preview, "image/") {
		t.Errorf("makePreview(png) preview = %q, want image mime", preview)
	}
}

func TestMakePreview_Empty(t *testing.T) {
	preview, kind := makePreview(nil)
	if preview != "" || kind != "" {
		t.Errorf("makePreview(nil) = (%q, %q), want both empty", preview, kind)
	}
}
