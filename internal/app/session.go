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

type KeyMeta struct {
	Key     string
	Size    int
	Preview string
}

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

func (s *Session) Keys() ([]KeyMeta, error) {
	if s.keys != nil {
		return s.keys, nil
	}
	return s.reloadKeys()
}

func (s *Session) Value(key []byte) ([]byte, error) {
	return s.Store.Get(key)
}

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

// Returns "" for non-JSON so the caller can fall back to the raw-text preview.
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
