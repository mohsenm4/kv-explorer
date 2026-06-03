package app

import (
	"fmt"

	"github.com/mohsenm4/kv-explorer/internal/kvstore"
	"github.com/mohsenm4/kv-explorer/internal/kvstore/badger"
	"github.com/mohsenm4/kv-explorer/internal/kvstore/pebble"
)

// OpenStore dispatches to the right engine adapter. This is the only place
// concrete adapter packages are imported — UI code talks to kvstore types only.
func OpenStore(kind kvstore.EngineKind, path string, opts kvstore.OpenOptions) (kvstore.Store, error) {
	switch kind {
	case kvstore.EnginePebble:
		return pebble.Open(path, opts)
	case kvstore.EngineBadger:
		return badger.Open(path, opts)
	case kvstore.EngineLevelDB:
		return nil, fmt.Errorf("engine %q not implemented yet", kind)
	default:
		return nil, fmt.Errorf("unknown engine %q", kind)
	}
}

// CountKeys iterates the entire store and returns the key count.
// This is a temporary helper used by Step 5 to verify the adapter; later
// the status bar will pull this from a Session.
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
