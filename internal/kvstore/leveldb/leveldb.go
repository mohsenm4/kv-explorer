package leveldb

import (
	"fmt"

	ldb "github.com/syndtr/goleveldb/leveldb"
	ldbiter "github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

type Store struct {
	db *ldb.DB
}

func Open(path string, opts kvstore.OpenOptions) (*Store, error) {
	o := &opt.Options{ReadOnly: opts.ReadOnly}
	db, err := ldb.OpenFile(path, o)
	if err != nil {
		return nil, fmt.Errorf("leveldb open %q: %w", path, err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Get(key []byte) ([]byte, error) {
	v, err := s.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(v))
	copy(out, v)
	return out, nil
}

func (s *Store) Set(key, value []byte) error {
	return s.db.Put(key, value, nil)
}

func (s *Store) Delete(key []byte) error {
	return s.db.Delete(key, nil)
}

func (s *Store) Iter(prefix []byte) (kvstore.Iterator, error) {
	var slice *util.Range
	if len(prefix) > 0 {
		slice = util.BytesPrefix(prefix)
	}
	it := s.db.NewIterator(slice, nil)
	return &iterator{it: it, first: true}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

type iterator struct {
	it    ldbiter.Iterator
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
	i.it.Release()
	return i.it.Error()
}
