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
}

func OpenSession(kind kvstore.EngineKind, path string, opts kvstore.OpenOptions) (*Session, error) {
	store, err := OpenStore(kind, path, opts)
	if err != nil {
		return nil, err
	}
	s := &Session{Engine: kind, Path: path, Store: store}
	if n, err := CountKeys(store); err == nil {
		s.KeyCount = n
	}
	if size, err := dirSize(path); err == nil {
		s.SizeBytes = size
	}
	return s, nil
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
