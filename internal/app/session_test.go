package app

import "testing"

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
			if got := makePreview([]byte(c.in)); got != c.want {
				t.Errorf("makePreview(%q) = %q, want %q", c.in, got, c.want)
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
		if got := makePreview([]byte(c.in)); got != c.want {
			t.Errorf("makePreview(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestMakePreview_NonJSONFallback(t *testing.T) {
	// Looks like JSON-ish text but isn't valid JSON: should fall back to
	// the raw-text preview rather than producing "".
	got := makePreview([]byte("{not json at all"))
	if got == "" {
		t.Errorf("makePreview fell through to empty for non-JSON text")
	}
}

func TestMakePreview_PlainText(t *testing.T) {
	got := makePreview([]byte("hello world"))
	if got != "hello world" {
		t.Errorf("makePreview(plain) = %q, want %q", got, "hello world")
	}
}

func TestMakePreview_Binary(t *testing.T) {
	got := makePreview([]byte{0x00, 0x01, 0x02, 0xff})
	if got == "" || got[0] != '[' {
		t.Errorf("makePreview(binary) = %q, want [mime · size]-style", got)
	}
}
