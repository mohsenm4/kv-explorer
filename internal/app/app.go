package app

import (
	"fmt"

	"github.com/mohsenm4/kv-explorer/internal/kvstore"
	"github.com/mohsenm4/kv-explorer/internal/kvstore/badger"
	"github.com/mohsenm4/kv-explorer/internal/kvstore/leveldb"
	"github.com/mohsenm4/kv-explorer/internal/kvstore/pebble"
)

// Only place concrete adapter packages are imported; UI talks to kvstore types.
func OpenStore(kind kvstore.EngineKind, path string, opts kvstore.OpenOptions) (kvstore.Store, error) {
	switch kind {
	case kvstore.EnginePebble:
		return pebble.Open(path, opts)
	case kvstore.EngineBadger:
		return badger.Open(path, opts)
	case kvstore.EngineLevelDB:
		return leveldb.Open(path, opts)
	default:
		return nil, fmt.Errorf("unknown engine %q", kind)
	}
}

func CountKeys(s kvstore.Store) (int, error) {
	it, err := s.Iter(nil)
	if err != nil {
		return 0, err
	}
	defer it.Close()
	n := 0
	for it.Next() {
		n++
	}
	return n, nil
}
