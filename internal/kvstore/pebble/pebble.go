package pebble

import (
	"fmt"

	"github.com/cockroachdb/pebble"

	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

type Store struct {
	db *pebble.DB
}

func Open(path string, opts kvstore.OpenOptions) (*Store, error) {
	db, err := pebble.Open(path, &pebble.Options{ReadOnly: opts.ReadOnly})
	if err != nil {
		return nil, fmt.Errorf("pebble open %q: %w", path, err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Get(key []byte) ([]byte, error) {
	val, closer, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	out := make([]byte, len(val))
	copy(out, val)
	return out, nil
}

func (s *Store) Set(key, value []byte) error {
	return s.db.Set(key, value, pebble.Sync)
}

func (s *Store) Delete(key []byte) error {
	return s.db.Delete(key, pebble.Sync)
}

func (s *Store) Iter(prefix []byte) (kvstore.Iterator, error) {
	o := &pebble.IterOptions{}
	if len(prefix) > 0 {
		o.LowerBound = prefix
		o.UpperBound = upperBound(prefix)
	}
	it, err := s.db.NewIter(o)
	if err != nil {
		return nil, err
	}
	return &iterator{it: it, first: true}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

type iterator struct {
	it    *pebble.Iterator
	first bool
}

func (i *iterator) Next() bool {
	if i.first {
		i.first = false
		return i.it.First()
	}
	return i.it.Next()
}

func (i *iterator) Entry() kvstore.Entry {
	k := i.it.Key()
	v := i.it.Value()
	out := kvstore.Entry{Key: make([]byte, len(k)), Value: make([]byte, len(v))}
	copy(out.Key, k)
	copy(out.Value, v)
	return out
}

func (i *iterator) Close() error {
	return i.it.Close()
}

// upperBound returns the smallest key strictly greater than every key with
// the given prefix, so pebble's range iterator stops at the prefix boundary.
func upperBound(prefix []byte) []byte {
	end := make([]byte, len(prefix))
	copy(end, prefix)
	for i := len(end) - 1; i >= 0; i-- {
		end[i]++
		if end[i] != 0 {
			return end
		}
	}
	return nil
}
