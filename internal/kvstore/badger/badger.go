package badger

import (
	"fmt"

	badgerdb "github.com/dgraph-io/badger/v4"

	"github.com/mohsenm4/kv-explorer/internal/kvstore"
)

type Store struct {
	db *badgerdb.DB
}

func Open(path string, opts kvstore.OpenOptions) (*Store, error) {
	// Badger's default value-log file size is 2 GiB which the OS reports as
	// allocated even when sparse. For an explorer UI 64 MiB segments are
	// plenty and keep the on-disk size readout sensible.
	o := badgerdb.DefaultOptions(path).
		WithReadOnly(opts.ReadOnly).
		WithLoggingLevel(badgerdb.WARNING).
		WithValueLogFileSize(64 << 20)
	db, err := badgerdb.Open(o)
	if err != nil {
		return nil, fmt.Errorf("badger open %q: %w", path, err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Get(key []byte) ([]byte, error) {
	var out []byte
	err := s.db.View(func(txn *badgerdb.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(v []byte) error {
			out = make([]byte, len(v))
			copy(out, v)
			return nil
		})
	})
	return out, err
}

func (s *Store) Set(key, value []byte) error {
	return s.db.Update(func(txn *badgerdb.Txn) error {
		return txn.Set(key, value)
	})
}

func (s *Store) Delete(key []byte) error {
	return s.db.Update(func(txn *badgerdb.Txn) error {
		return txn.Delete(key)
	})
}

func (s *Store) Iter(prefix []byte) (kvstore.Iterator, error) {
	txn := s.db.NewTransaction(false)
	it := txn.NewIterator(badgerdb.DefaultIteratorOptions)
	return &iterator{txn: txn, it: it, prefix: prefix, first: true}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

type iterator struct {
	txn    *badgerdb.Txn
	it     *badgerdb.Iterator
	prefix []byte
	first  bool
}

func (i *iterator) Next() bool {
	if i.first {
		i.first = false
		if len(i.prefix) > 0 {
			i.it.Seek(i.prefix)
		} else {
			i.it.Rewind()
		}
	} else {
		i.it.Next()
	}
	if !i.it.Valid() {
		return false
	}
	if len(i.prefix) > 0 && !i.it.ValidForPrefix(i.prefix) {
		return false
	}
	return true
}

func (i *iterator) Entry() kvstore.Entry {
	item := i.it.Item()
	key := item.KeyCopy(nil)
	var val []byte
	_ = item.Value(func(b []byte) error {
		val = make([]byte, len(b))
		copy(val, b)
		return nil
	})
	return kvstore.Entry{Key: key, Value: val}
}

func (i *iterator) Close() error {
	i.it.Close()
	i.txn.Discard()
	return nil
}
