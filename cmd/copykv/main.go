// copykv copies all keys/values from one kvstore engine to another.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mohsenm4/kv-explorer/internal/kvstore"
	"github.com/mohsenm4/kv-explorer/internal/kvstore/badger"
	"github.com/mohsenm4/kv-explorer/internal/kvstore/leveldb"
	"github.com/mohsenm4/kv-explorer/internal/kvstore/pebble"
)

func main() {
	srcEngine := flag.String("src-engine", "", "leveldb | badger | pebble")
	srcPath := flag.String("src-path", "", "source database path")
	dstEngine := flag.String("dst-engine", "", "leveldb | badger | pebble")
	dstPath := flag.String("dst-path", "", "destination database path")
	flag.Parse()

	if *srcEngine == "" || *srcPath == "" || *dstEngine == "" || *dstPath == "" {
		log.Fatal("usage: copykv --src-engine=<e> --src-path=<p> --dst-engine=<e> --dst-path=<p>")
	}

	src, err := openStore(*srcEngine, *srcPath, false)
	if err != nil {
		log.Fatalf("open source: %v", err)
	}
	defer src.Close()

	if err := os.MkdirAll(*dstPath, 0o755); err != nil {
		log.Fatalf("mkdir dst: %v", err)
	}
	dst, err := openStore(*dstEngine, *dstPath, false)
	if err != nil {
		log.Fatalf("open destination: %v", err)
	}
	defer dst.Close()

	it, err := src.Iter(nil)
	if err != nil {
		log.Fatalf("iter: %v", err)
	}
	defer it.Close()

	var copied int
	for it.Next() {
		e := it.Entry()
		k := make([]byte, len(e.Key))
		copy(k, e.Key)
		v := make([]byte, len(e.Value))
		copy(v, e.Value)
		if err := dst.Set(k, v); err != nil {
			log.Fatalf("set %q: %v", k, err)
		}
		copied++
	}
	fmt.Printf("copied %d keys from %s(%s) to %s(%s)\n", copied, *srcEngine, *srcPath, *dstEngine, *dstPath)
}

func openStore(engine, path string, readOnly bool) (kvstore.Store, error) {
	opts := kvstore.OpenOptions{ReadOnly: readOnly}
	switch engine {
	case "leveldb":
		return leveldb.Open(path, opts)
	case "badger":
		return badger.Open(path, opts)
	case "pebble":
		return pebble.Open(path, opts)
	default:
		return nil, fmt.Errorf("unknown engine %q", engine)
	}
}
