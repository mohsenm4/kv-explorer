// stress fills a Pebble store with a heavy, varied dataset for load testing
// the UI. Usage: go run ./cmd/stress <path> [--count N] [--large] [--huge]
package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mohsenm4/kv-explorer/internal/kvstore"
	"github.com/mohsenm4/kv-explorer/internal/kvstore/pebble"
)

func main() {
	count := flag.Int("count", 50_000, "number of bulk JSON keys")
	large := flag.Bool("large", true, "include 1KB–10MB single-value keys")
	huge := flag.Bool("huge", false, "also write a 100MB single value (slow)")
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("usage: stress <pebble-path> [--count N] [--large=false] [--huge]")
	}
	path := flag.Arg(0)

	s, err := pebble.Open(path, kvstore.OpenOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	set := func(k string, v []byte) {
		if err := s.Set([]byte(k), v); err != nil {
			log.Fatalf("set %q: %v", k, err)
		}
	}
	setRaw := func(k []byte, v []byte) {
		if err := s.Set(k, v); err != nil {
			log.Fatalf("set raw: %v", err)
		}
	}

	start := time.Now()
	roles := []string{"admin", "user", "guest", "viewer", "moderator", "auditor"}

	log.Printf("writing %d bulk JSON keys under bulk/users/…", *count)
	for i := 0; i < *count; i++ {
		v, _ := json.Marshal(map[string]any{
			"id":     i,
			"name":   fmt.Sprintf("user_%d", i),
			"email":  fmt.Sprintf("u%d@example.com", i),
			"role":   roles[i%len(roles)],
			"score":  i * 7 % 10_000,
			"active": i%2 == 0,
			"tags":   []string{"alpha", "beta", roles[i%len(roles)]},
		})
		set(fmt.Sprintf("bulk/users/%08d", i), v)
		if i > 0 && i%10_000 == 0 {
			rate := float64(i) / time.Since(start).Seconds()
			log.Printf("  %d / %d (%.0f keys/s)", i, *count, rate)
		}
	}

	log.Printf("writing prefix-tree fan-out (deep/aaa/bbb/leaf)")
	for i := 0; i < 50; i++ {
		for j := 0; j < 50; j++ {
			set(fmt.Sprintf("deep/%03d/%03d/leaf", i, j), []byte(fmt.Sprintf("leaf-%d-%d", i, j)))
		}
	}

	log.Printf("writing pattern-tagged values (UUID, URL, email)")
	for i := 0; i < 100; i++ {
		set(fmt.Sprintf("patterns/uuid/%03d", i),
			[]byte(fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", i, i*3, i*5, i*7, i*11)))
		set(fmt.Sprintf("patterns/url/%03d", i),
			[]byte(fmt.Sprintf("https://api.example.com/v1/users/%d", i)))
		set(fmt.Sprintf("patterns/email/%03d", i),
			[]byte(fmt.Sprintf("user.%d+kv@subdomain.example.com", i)))
	}

	log.Printf("writing unix timestamps in s / ms / ns")
	now := time.Now().Unix()
	for i := 0; i < 30; i++ {
		set(fmt.Sprintf("timestamps/s/%02d", i), []byte(fmt.Sprintf("%d", now-int64(i)*86400)))
		set(fmt.Sprintf("timestamps/ms/%02d", i), []byte(fmt.Sprintf("%d", (now-int64(i)*86400)*1000)))
		set(fmt.Sprintf("timestamps/ns/%02d", i), []byte(fmt.Sprintf("%d", (now-int64(i)*86400)*1_000_000_000)))
	}

	log.Printf("writing edge cases (long key, binary key, unicode, control bytes, empty)")
	set("edge/long-key/"+strings.Repeat("a", 800), []byte("value-with-very-long-key"))
	set("edge/empty", []byte{})
	set("edge/unicode/中文/日本語/Ελληνικά", []byte("multi-script key"))
	set("edge/special chars (space, parens) [brackets]!", []byte("special key chars"))
	set("edge/json/empty-object", []byte(`{}`))
	set("edge/json/empty-array", []byte(`[]`))
	set("edge/json/single-key", []byte(`{"x":1}`))
	set("edge/json/deeply-nested", []byte(buildDeepJSON(40)))
	set("edge/json/wide", []byte(buildWideJSON(500)))
	set("edge/json/array-of-1000", []byte(buildJSONArray(1000)))
	set("edge/json/invalid", []byte(`{not valid json,`))

	// Binary key — DetectContent should route this to hex view.
	binKey := make([]byte, 32)
	_, _ = rand.Read(binKey)
	setRaw(binKey, []byte("value under random binary key"))

	// Binary value with non-text bytes — Auto picks Hex.
	binValue := make([]byte, 4096)
	_, _ = rand.Read(binValue)
	set("edge/binary/random", binValue)

	if *large {
		log.Printf("writing single large values (1KB → 10MB)")
		for _, sz := range []int{1 << 10, 10 << 10, 100 << 10, 1 << 20, 5 << 20, 10 << 20} {
			v := make([]byte, sz)
			_, _ = rand.Read(v)
			set(fmt.Sprintf("large/%s-bin", humanLabel(sz)), v)

			textRepeat := strings.Repeat("Lorem ipsum dolor sit amet, consectetur adipiscing elit. ", sz/57+1)
			set(fmt.Sprintf("large/%s-text", humanLabel(sz)), []byte(textRepeat[:sz]))
		}
	}

	if *huge {
		log.Printf("writing 100MB single value (this can take a minute)")
		v := make([]byte, 100<<20)
		_, _ = rand.Read(v)
		set("huge/100MB-bin", v)
	}

	log.Printf("done in %v — open the DB at %s with engine=Pebble", time.Since(start), path)
	if info, err := dirSize(path); err == nil {
		log.Printf("on-disk size: %s", info)
	}
}

func buildDeepJSON(depth int) string {
	var b strings.Builder
	for i := 0; i < depth; i++ {
		b.WriteString(`{"level":`)
	}
	b.WriteString(`"bottom"`)
	for i := 0; i < depth; i++ {
		b.WriteString(`}`)
	}
	return b.String()
}

func buildWideJSON(fields int) string {
	m := make(map[string]any, fields)
	for i := 0; i < fields; i++ {
		m[fmt.Sprintf("field_%04d", i)] = i
	}
	out, _ := json.Marshal(m)
	return string(out)
}

func buildJSONArray(n int) string {
	items := make([]map[string]any, n)
	for i := range items {
		items[i] = map[string]any{"i": i, "name": fmt.Sprintf("item_%d", i)}
	}
	out, _ := json.Marshal(items)
	return string(out)
}

func humanLabel(n int) string {
	switch {
	case n >= 1<<20:
		return fmt.Sprintf("%dMB", n>>20)
	case n >= 1<<10:
		return fmt.Sprintf("%dKB", n>>10)
	default:
		return fmt.Sprintf("%dB", n)
	}
}

func dirSize(path string) (string, error) {
	var total int64
	err := walkSize(path, &total)
	if err != nil {
		return "", err
	}
	switch {
	case total >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(total)/(1<<30)), nil
	case total >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(total)/(1<<20)), nil
	case total >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(total)/(1<<10)), nil
	default:
		return fmt.Sprintf("%d B", total), nil
	}
}

func walkSize(root string, total *int64) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	for _, e := range entries {
		full := root + "/" + e.Name()
		if e.IsDir() {
			if err := walkSize(full, total); err != nil {
				return err
			}
			continue
		}
		info, err := e.Info()
		if err != nil {
			return err
		}
		*total += info.Size()
	}
	return nil
}
