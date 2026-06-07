package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

// KeyMeta is the lightweight per-row info the UI keeps in memory: the key,
// its value size, and a short preview string. The preview is computed once
// at session load and reused so the table doesn't refetch values on every
// cell render.
type KeyMeta struct {
	Key     string
	Size    int
	Preview string
}

// Session represents one open database. Multiple sessions can coexist
// (one per UI tab once Step 17 lands).
type Session struct {
	Engine    kvstore.EngineKind
	Path      string
	Store     kvstore.Store
	KeyCount  int
	SizeBytes int64

	keys []KeyMeta
}

func OpenSession(kind kvstore.EngineKind, path string, opts kvstore.OpenOptions) (*Session, error) {
	store, err := OpenStore(kind, path, opts)
	if err != nil {
		return nil, err
	}
	s := &Session{Engine: kind, Path: path, Store: store}
	if _, err := s.reloadKeys(); err != nil {
		_ = store.Close()
		return nil, err
	}
	if size, err := dirSize(path); err == nil {
		s.SizeBytes = size
	}
	return s, nil
}

// Keys returns the cached list of (key, size, preview) tuples sorted by
// key. The cache is populated on session open and invalidated by Refresh.
func (s *Session) Keys() ([]KeyMeta, error) {
	if s.keys != nil {
		return s.keys, nil
	}
	return s.reloadKeys()
}

// Value fetches the value bytes for a key directly from the store. Values
// are never cached at the session level — callers can wrap their own
// cache if hot lookups warrant it.
func (s *Session) Value(key []byte) ([]byte, error) {
	return s.Store.Get(key)
}

// Refresh re-iterates the store, refreshing the key cache and counts.
func (s *Session) Refresh() error {
	_, err := s.reloadKeys()
	return err
}

func (s *Session) reloadKeys() ([]KeyMeta, error) {
	it, err := s.Store.Iter(nil)
	if err != nil {
		return nil, err
	}
	defer it.Close()
	var out []KeyMeta
	for it.Next() {
		e := it.Entry()
		out = append(out, KeyMeta{
			Key:     string(e.Key),
			Size:    len(e.Value),
			Preview: makePreview(e.Value),
		})
	}
	s.keys = out
	s.KeyCount = len(out)
	return out, nil
}

func (s *Session) Close() error {
	if s == nil || s.Store == nil {
		return nil
	}
	return s.Store.Close()
}

// makePreview returns a short one-line preview for the table cell.
// JSON gets a field-summary like `{id, name, +3}` or `[12 items]` so
// users can recognise the shape without reading raw braces. Binary content
// is summarised by MIME and size; other text gets its first ~120 chars
// with whitespace collapsed.
func makePreview(v []byte) string {
	if len(v) == 0 {
		return ""
	}
	mime := http.DetectContentType(v)
	isText := strings.HasPrefix(mime, "text/") || mime == "application/json"
	if !isText {
		if utf8.Valid(v) && !hasControlBytes(v) {
			isText = true
		}
	}
	if !isText {
		return fmt.Sprintf("[%s · %s]", mime, humanSize(int64(len(v))))
	}

	if jp := jsonPreview(v); jp != "" {
		return jp
	}

	const max = 120
	limit := max * 4
	if len(v) < limit {
		limit = len(v)
	}
	s := string(v[:limit])
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' || s[i] == '\r' || s[i] == '\t' {
			s = s[:i] + " " + s[i+1:]
		}
	}
	if len(s) > max {
		return s[:max] + "…"
	}
	return s
}

// jsonPreview returns a structural one-liner for JSON values:
//   - object  → `{id, name, email, +3}` (first 3 keys, +N for the rest)
//   - array   → `[12 items]` for non-empty, `[]` for empty
//   - empty object → `{}`
//
// Keys are listed in their on-wire order via json.Decoder so the preview
// matches what the user sees in the editor. Returns "" if v isn't JSON,
// letting the caller fall back to the raw-text preview.
func jsonPreview(v []byte) string {
	trimmed := bytes.TrimSpace(v)
	if len(trimmed) == 0 {
		return ""
	}
	switch trimmed[0] {
	case '{':
		keys, err := topLevelObjectKeys(trimmed)
		if err != nil {
			return ""
		}
		if len(keys) == 0 {
			return "{}"
		}
		const show = 3
		head := keys
		extra := 0
		if len(keys) > show {
			head = keys[:show]
			extra = len(keys) - show
		}
		if extra > 0 {
			return fmt.Sprintf("{%s, +%d}", strings.Join(head, ", "), extra)
		}
		return fmt.Sprintf("{%s}", strings.Join(head, ", "))
	case '[':
		n, err := topLevelArrayLen(trimmed)
		if err != nil {
			return ""
		}
		if n == 0 {
			return "[]"
		}
		if n == 1 {
			return "[1 item]"
		}
		return fmt.Sprintf("[%d items]", n)
	}
	return ""
}

func topLevelObjectKeys(v []byte) ([]string, error) {
	dec := json.NewDecoder(bytes.NewReader(v))
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if d, ok := tok.(json.Delim); !ok || d != '{' {
		return nil, fmt.Errorf("not an object")
	}
	var keys []string
	for dec.More() {
		kt, err := dec.Token()
		if err != nil {
			return nil, err
		}
		k, ok := kt.(string)
		if !ok {
			return nil, fmt.Errorf("non-string key")
		}
		keys = append(keys, k)
		// Skip the value (could be any JSON).
		var skip json.RawMessage
		if err := dec.Decode(&skip); err != nil {
			return nil, err
		}
	}
	return keys, nil
}

func topLevelArrayLen(v []byte) (int, error) {
	dec := json.NewDecoder(bytes.NewReader(v))
	tok, err := dec.Token()
	if err != nil {
		return 0, err
	}
	if d, ok := tok.(json.Delim); !ok || d != '[' {
		return 0, fmt.Errorf("not an array")
	}
	n := 0
	for dec.More() {
		var skip json.RawMessage
		if err := dec.Decode(&skip); err != nil {
			return 0, err
		}
		n++
	}
	return n, nil
}

func hasControlBytes(v []byte) bool {
	for _, b := range v {
		if b < 0x09 || (b > 0x0d && b < 0x20) {
			return true
		}
	}
	return false
}

func humanSize(b int64) string {
	switch {
	case b < 1024:
		return fmt.Sprintf("%d B", b)
	case b < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(b)/1024)
	case b < 1024*1024*1024:
		return fmt.Sprintf("%.1f MB", float64(b)/(1024*1024))
	default:
		return fmt.Sprintf("%.1f GB", float64(b)/(1024*1024*1024))
	}
}

func dirSize(path string) (int64, error) {
	var total int64
	err := filepath.WalkDir(path, func(_ string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		total += info.Size()
		return nil
	})
	return total, err
}
