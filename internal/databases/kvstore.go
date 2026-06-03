// Package databases defines the shared KVStore interface that every database
// adapter (pebble, badger, leveldb) must implement.
package databases

// KVStore is the unified interface every database adapter implements.
type KVStore interface {
	Open(path string) error
	Close() error

	Get(key []byte) ([]byte, error)
	Set(key, value []byte) error
	Delete(key []byte) error

	Iterate(prefix []byte, fn func(key, value []byte) bool) error

	Stats() Stats
}

// Stats holds runtime statistics for a KVStore.
type Stats struct {
	KeyCount  uint64
	SizeBytes uint64
}
