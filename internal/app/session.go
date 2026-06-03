package app

import (
	"os"
	"path/filepath"

	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

// KeyMeta is the lightweight per-row info the UI keeps in memory: just the
// key and its value size. Actual values are fetched on demand via Value().
type KeyMeta struct {
	Key  string
	Size int
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

// Keys returns the cached list of (key, size) pairs sorted by key. The
// cache is populated on session open and invalidated by Refresh.
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
		out = append(out, KeyMeta{Key: string(e.Key), Size: len(e.Value)})
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
