// seed populates a Pebble store with synthetic data shaped like the
// wireframes in docs/design/spec.md (users/, sessions/, logs/, cache/),
// plus a few binary samples for the editor's image/hex modes.
// Usage: go run ./cmd/seed <pebble-path>
package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/mohsenm4/kv-explorer/internal/kvstore"
	"github.com/mohsenm4/kv-explorer/internal/kvstore/pebble"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: seed <pebble-path>")
	}
	path := os.Args[1]

	s, err := pebble.Open(path, kvstore.OpenOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	put := func(k, v string) {
		if err := s.Set([]byte(k), []byte(v)); err != nil {
			log.Fatal(err)
		}
	}

	names := []string{"ali", "sara", "reza", "mina", "navid", "leyla", "amir", "hassan", "zahra", "mahdi", "narges", "ahmad", "fatemeh", "hossein", "neda"}
	roles := []string{"admin", "user", "guest", "viewer"}

	for i := 1; i <= 120; i++ {
		v, _ := json.Marshal(map[string]any{
			"name":   names[i%len(names)],
			"age":    20 + i%50,
			"role":   roles[i%len(roles)],
			"active": i%3 != 0,
		})
		put(fmt.Sprintf("users/%04d", i), string(v))
	}

	sessionIDs := []string{"abc123", "xyz789", "def456", "qrs012", "mno345", "jkl678", "pqr901", "stu234", "vwx567", "yza890",
		"bcd111", "efg222", "hij333", "klm444", "nop555", "qrs666", "tuv777", "wxy888", "zab999", "cde000",
		"fgh1aa", "ijk2bb", "lmn3cc", "opq4dd", "rst5ee", "uvw6ff", "xyz7gg", "abc8hh", "def9ii", "ghi0jj"}
	for i, id := range sessionIDs {
		v, _ := json.Marshal(map[string]any{
			"user_id": fmt.Sprintf("users/%04d", (i%120)+1),
			"ts":      1717459200 + int64(i*3600),
			"ttl":     86400,
		})
		put("sessions/"+id, string(v))
	}

	logMsgs := []string{"app startup", "user login", "query executed", "cache miss", "cache hit", "user logout", "config reloaded", "checkpoint", "compaction started", "compaction finished"}
	for i := 1; i <= 30; i++ {
		put(fmt.Sprintf("logs/2026-06-03/%04d", i),
			fmt.Sprintf(`{"level":"info","msg":"%s","ts":%d}`, logMsgs[i%len(logMsgs)], 1717459200+i*60))
	}

	put("cache/v1/users", `["ali","sara","reza","mina"]`)
	put("cache/v1/stats", `{"hits":1234,"misses":56,"ratio":0.957}`)
	put("cache/v1/config", `{"theme":"light","lang":"en"}`)
	put("cache/v1/version", "1.4.2")

	put("meta/created_at", "2026-06-03T12:14:00Z")
	put("meta/schema_version", "3")
	put("meta/last_compaction", "2026-06-02T22:00:00Z")

	// Binary samples for the editor's image / hex modes (Step 10c).
	put("media/red.png", string(makePNG(96, 96, color.RGBA{R: 220, G: 50, B: 50, A: 255})))
	put("media/blue.png", string(makePNG(96, 96, color.RGBA{R: 50, G: 80, B: 220, A: 255})))
	put("media/wide.png", string(makePNG(192, 64, color.RGBA{R: 50, G: 180, B: 80, A: 255})))

	randBytes := make([]byte, 256)
	_, _ = rand.Read(randBytes)
	put("media/random.bin", string(randBytes))

	// Fake MP3 header so DetectContent sees audio/mpeg.
	mp3 := append([]byte{0xff, 0xfb, 0x90, 0x00}, randBytes[:200]...)
	put("media/song.mp3", string(mp3))

	fmt.Printf("seeded %s (192 keys: 120 users, 30 sessions, 30 logs, 4 cache, 3 meta, 5 media)\n", path)
}

func makePNG(w, h int, c color.RGBA) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Fatalf("encode png: %v", err)
	}
	return buf.Bytes()
}
