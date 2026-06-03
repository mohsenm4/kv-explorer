package kvstore

type EngineKind string

const (
	EnginePebble  EngineKind = "pebble"
	EngineBadger  EngineKind = "badger"
	EngineLevelDB EngineKind = "leveldb"
)

type Entry struct {
	Key   []byte
	Value []byte
}

type OpenOptions struct {
	ReadOnly bool
}

type Store interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte) error
	Delete(key []byte) error
	Iter(prefix []byte) (Iterator, error)
	Close() error
}

type Iterator interface {
	Next() bool
	Entry() Entry
	Close() error
}
