package app

import (
	"os"
	"path/filepath"

	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

// Session represents one open database. Multiple sessions can coexist
// (one per UI tab once Step 17 lands).
type Session struct {
	Engine    kvstore.EngineKind
	Path      string
	Store     kvstore.Store
	KeyCount  int
	SizeBytes int64

	entries []kvstore.Entry
}

func OpenSession(kind kvstore.EngineKind, path string, opts kvstore.OpenOptions) (*Session, error) {
	store, err := OpenStore(kind, path, opts)
	if err != nil {
		return nil, err
	}
	s := &Session{Engine: kind, Path: path, Store: store}
	if _, err := s.reloadEntries(); err != nil {
		_ = store.Close()
		return nil, err
	}
	if size, err := dirSize(path); err == nil {
		s.SizeBytes = size
	}
	return s, nil
}

// Entries returns all key/value pairs in the store, cached after the first
// call. Use Refresh to invalidate the cache after a write.
func (s *Session) Entries() ([]kvstore.Entry, error) {
	if s.entries != nil {
		return s.entries, nil
	}
	return s.reloadEntries()
}

// Refresh re-iterates the store and updates the cached entries and counts.
func (s *Session) Refresh() error {
	_, err := s.reloadEntries()
	return err
}

func (s *Session) reloadEntries() ([]kvstore.Entry, error) {
	it, err := s.Store.Iter(nil)
	if err != nil {
		return nil, err
	}
	defer it.Close()
	var out []kvstore.Entry
	for it.Next() {
		out = append(out, it.Entry())
	}
	s.entries = out
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
