package config

import (
	"testing"
	"time"
)

func TestAddRecent_DedupAndCap(t *testing.T) {
	c := Config{}
	for i := 0; i < 12; i++ {
		c.AddRecent("/path/"+string(rune('a'+i)), "pebble")
	}
	if got := len(c.Recents); got != maxRecents {
		t.Errorf("len = %d, want %d", got, maxRecents)
	}

	// Re-adding an existing path moves it to the top, no duplicate.
	c.AddRecent("/path/d", "pebble")
	if len(c.Recents) != maxRecents {
		t.Errorf("len after re-add = %d, want %d", len(c.Recents), maxRecents)
	}
	if c.Recents[0].Path != "/path/d" {
		t.Errorf("front = %q, want /path/d", c.Recents[0].Path)
	}
}

func TestAddRecent_Timestamp(t *testing.T) {
	c := Config{}
	before := time.Now()
	c.AddRecent("/p", "leveldb")
	if c.Recents[0].OpenedAt.Before(before) {
		t.Errorf("timestamp %v before call time %v", c.Recents[0].OpenedAt, before)
	}
}
