package pebble_test

import (
	"bytes"
	"errors"
	"testing"

	pebbledb "github.com/cockroachdb/pebble"

	"github.com/mohsenm4/kv-explorer/internal/kvstore"
	"github.com/mohsenm4/kv-explorer/internal/kvstore/pebble"
)

func TestStore_RoundTrip(t *testing.T) {
	s, err := pebble.Open(t.TempDir(), kvstore.OpenOptions{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer s.Close()

	if err := s.Set([]byte("k1"), []byte("v1")); err != nil {
		t.Fatalf("set: %v", err)
	}
	if err := s.Set([]byte("k2"), []byte("v2")); err != nil {
		t.Fatalf("set: %v", err)
	}

	got, err := s.Get([]byte("k1"))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !bytes.Equal(got, []byte("v1")) {
		t.Errorf("get k1 = %q, want v1", got)
	}

	if err := s.Delete([]byte("k1")); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := s.Get([]byte("k1")); !errors.Is(err, pebbledb.ErrNotFound) {
		t.Errorf("get after delete = %v, want ErrNotFound", err)
	}
}

func TestStore_PrefixIter(t *testing.T) {
	s, err := pebble.Open(t.TempDir(), kvstore.OpenOptions{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer s.Close()

	for _, kv := range []struct{ k, v string }{
		{"a/1", "v1"},
		{"a/2", "v2"},
		{"b/1", "v3"},
	} {
		if err := s.Set([]byte(kv.k), []byte(kv.v)); err != nil {
			t.Fatalf("set %q: %v", kv.k, err)
		}
	}

	it, err := s.Iter([]byte("a/"))
	if err != nil {
		t.Fatalf("iter: %v", err)
	}
	defer it.Close()

	var keys []string
	for it.Next() {
		keys = append(keys, string(it.Entry().Key))
	}
	if len(keys) != 2 || keys[0] != "a/1" || keys[1] != "a/2" {
		t.Errorf("prefix iter = %v, want [a/1 a/2]", keys)
	}
}
