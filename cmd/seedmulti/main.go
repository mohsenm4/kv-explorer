// seedmulti populates a LevelDB or BadgerDB store with a diverse mix of
// value types — light/heavy text, JSON, XML, CSV, YAML, Markdown, PNG/JPEG
// images, fake MP3/WAV audio, a real minimal XLSX, and random binary blobs.
// It exists to give the kv-explorer UI varied content to render.
//
// Usage:
//
//	go run ./cmd/seedmulti --engine=leveldb --path=/path/to/db
//	go run ./cmd/seedmulti --engine=badger  --path=/path/to/db
package main

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"strings"

	"github.com/mohsenm4/kv-explorer/internal/kvstore"
	"github.com/mohsenm4/kv-explorer/internal/kvstore/badger"
	"github.com/mohsenm4/kv-explorer/internal/kvstore/leveldb"
)

func main() {
	engine := flag.String("engine", "", "leveldb | badger")
	path := flag.String("path", "", "filesystem path for the database")
	flag.Parse()

	if *engine == "" || *path == "" {
		log.Fatal("usage: seedmulti --engine=<leveldb|badger> --path=<dir>")
	}

	if err := os.MkdirAll(*path, 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}

	store, err := openStore(*engine, *path)
	if err != nil {
		log.Fatalf("open %s at %s: %v", *engine, *path, err)
	}
	defer store.Close()

	count := seed(store)
	fmt.Printf("seeded %d keys into %s at %s\n", count, *engine, *path)
}

func openStore(engine, path string) (kvstore.Store, error) {
	switch engine {
	case "leveldb":
		return leveldb.Open(path, kvstore.OpenOptions{})
	case "badger":
		return badger.Open(path, kvstore.OpenOptions{})
	default:
		return nil, fmt.Errorf("unknown engine %q", engine)
	}
}

func seed(s kvstore.Store) int {
	n := 0
	put := func(k string, v []byte) {
		if err := s.Set([]byte(k), v); err != nil {
			log.Fatalf("set %s: %v", k, err)
		}
		n++
	}

	// --- Light text -----------------------------------------------------
	put("text/light/greeting/en", []byte("Hello, world!"))
	put("text/light/greeting/fa", []byte("سلام دنیا"))
	put("text/light/greeting/ja", []byte("こんにちは世界"))
	put("text/light/quote", []byte("Simplicity is the soul of efficiency. — Austin Freeman"))
	put("text/light/tagline", []byte("kv-explorer: inspect any key-value store."))

	// --- Heavy text -----------------------------------------------------
	put("text/heavy/lorem/short", []byte(loremParagraph(3)))
	put("text/heavy/lorem/medium", []byte(loremParagraph(20)))
	put("text/heavy/lorem/long", []byte(loremParagraph(120)))
	put("text/heavy/article/release-notes", []byte(releaseNotes()))

	// --- Markdown -------------------------------------------------------
	put("text/markdown/readme", []byte(readmeMarkdown()))
	put("text/markdown/changelog", []byte("# Changelog\n\n## v1.0.0\n- Initial release.\n- Added LevelDB adapter.\n- Added BadgerDB adapter.\n"))

	// --- JSON -----------------------------------------------------------
	cfg, _ := json.MarshalIndent(map[string]any{
		"app":     "kv-explorer",
		"version": "1.0.0",
		"theme":   "dark",
		"window":  map[string]int{"width": 1280, "height": 800},
		"recent":  []string{"~/db/users.pebble", "~/db/sessions.badger"},
	}, "", "  ")
	put("json/config", cfg)

	for i, u := range users() {
		v, _ := json.Marshal(u)
		put(fmt.Sprintf("json/users/%04d", i+1), v)
	}

	stats, _ := json.MarshalIndent(map[string]any{
		"hits": 12345, "misses": 67, "ratio": 0.994,
		"by_engine": map[string]int{"pebble": 8000, "badger": 3000, "leveldb": 1345},
	}, "", "  ")
	put("json/cache/stats", stats)

	// --- CSV ------------------------------------------------------------
	put("csv/sales/q1", []byte(csvSales()))
	put("csv/users", []byte("id,name,email,active\n1,Ali,ali@example.com,true\n2,Sara,sara@example.com,true\n3,Reza,reza@example.com,false\n"))

	// --- XML / YAML -----------------------------------------------------
	put("xml/feed", []byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>kv-explorer release feed</title>
    <link>https://example.com/kv-explorer</link>
    <item><title>v1.0 ships</title><description>Initial public release.</description></item>
    <item><title>v0.9 beta</title><description>Beta for early testers.</description></item>
  </channel>
</rss>`))
	put("yaml/config", []byte("app: kv-explorer\nversion: 1.0.0\nlog:\n  level: info\n  rotate: daily\nengines:\n  - pebble\n  - badger\n  - leveldb\n"))

	// --- Images (PNG) ---------------------------------------------------
	put("image/png/red", makePNG(96, 96, color.RGBA{R: 220, G: 50, B: 50, A: 255}))
	put("image/png/blue", makePNG(96, 96, color.RGBA{R: 50, G: 80, B: 220, A: 255}))
	put("image/png/green/wide", makePNG(256, 96, color.RGBA{R: 60, G: 180, B: 90, A: 255}))
	put("image/png/gradient", makeGradientPNG(192, 128))

	// --- Images (JPEG) --------------------------------------------------
	put("image/jpeg/sunset", makeJPEG(160, 120))

	// --- Audio (fake but with real magic bytes) -------------------------
	put("audio/mp3/song", makeFakeMP3())
	put("audio/wav/beep", makeFakeWAV())

	// --- Excel ----------------------------------------------------------
	put("excel/sales-q1.xlsx", makeXLSX())

	// --- Binary blobs ---------------------------------------------------
	put("binary/random/1k", randomBytes(1024))
	put("binary/random/16k", randomBytes(16*1024))
	put("binary/random/128k", randomBytes(128*1024))
	put("binary/random/1m", randomBytes(1024*1024))
	put("binary/zeros/4k", make([]byte, 4096))

	// --- Hierarchical samples mirroring real apps -----------------------
	for i := 1; i <= 30; i++ {
		v, _ := json.Marshal(map[string]any{
			"user_id": fmt.Sprintf("u-%04d", i),
			"action":  []string{"login", "logout", "click", "purchase", "search"}[i%5],
			"ts":      1717459200 + i*60,
		})
		put(fmt.Sprintf("logs/2026-06-06/%04d", i), v)
	}
	for i := 1; i <= 20; i++ {
		v, _ := json.Marshal(map[string]any{"hits": i * 17, "miss": i % 4})
		put(fmt.Sprintf("metrics/2026-06-06/%02d", i), v)
	}

	// --- Meta -----------------------------------------------------------
	put("meta/created_at", []byte("2026-06-06T10:00:00Z"))
	put("meta/schema_version", []byte("1"))
	put("meta/source", []byte("seedmulti"))

	return n
}

// --- Helpers ----------------------------------------------------------------

func loremParagraph(sentences int) string {
	s := []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.",
		"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum.",
		"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia.",
		"Curabitur pretium tincidunt lacus, nulla gravida orci a odio.",
		"Pellentesque habitant morbi tristique senectus et netus et malesuada.",
	}
	var b strings.Builder
	for i := 0; i < sentences; i++ {
		b.WriteString(s[i%len(s)])
		b.WriteByte(' ')
	}
	return strings.TrimSpace(b.String())
}

func readmeMarkdown() string {
	return `# kv-explorer

A desktop GUI tool for inspecting key-value databases.

## Features
- Open **PebbleDB**, **BadgerDB**, and **LevelDB** stores
- Prefix tree + table view
- Inline editor with UTF-8 / Hex / JSON toggle
- Multi-tab database comparison

## Install
` + "```" + `bash
go install github.com/mohsenm4/kv-explorer/cmd/kvexplorer@latest
` + "```" + `

## Status
Early but usable. See ` + "`docs/design/`" + ` for the visual spec.
`
}

func releaseNotes() string {
	var b strings.Builder
	b.WriteString("kv-explorer Release Notes\n=========================\n\n")
	for v := 1; v <= 10; v++ {
		fmt.Fprintf(&b, "v0.%d.0\n------\n", v)
		fmt.Fprintf(&b, "- Improved %s rendering for large values.\n", []string{"hex", "json", "image", "text"}[v%4])
		fmt.Fprintf(&b, "- Fixed crash when opening read-only databases.\n")
		fmt.Fprintf(&b, "- Performance: %d%% faster prefix scans on cold cache.\n\n", 5+v*3)
	}
	return b.String()
}

func users() []map[string]any {
	names := []string{"Ali", "Sara", "Reza", "Mina", "Navid", "Leyla", "Amir", "Hassan", "Zahra", "Mahdi"}
	roles := []string{"admin", "user", "guest", "viewer"}
	out := make([]map[string]any, 0, 25)
	for i := 0; i < 25; i++ {
		out = append(out, map[string]any{
			"id":     fmt.Sprintf("u-%04d", i+1),
			"name":   names[i%len(names)],
			"role":   roles[i%len(roles)],
			"age":    20 + i%30,
			"active": i%3 != 0,
		})
	}
	return out
}

func csvSales() string {
	var b strings.Builder
	b.WriteString("month,product,units,revenue\n")
	products := []string{"Widget", "Gadget", "Sprocket", "Cog"}
	for m := 1; m <= 12; m++ {
		for _, p := range products {
			fmt.Fprintf(&b, "2026-%02d,%s,%d,%.2f\n", m, p, 100+m*7, 250.50+float64(m)*float64(p[0]))
		}
	}
	return b.String()
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

func makeGradientPNG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r := uint8(255 * x / w)
			g := uint8(255 * y / h)
			b := uint8(255 - r/2)
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Fatalf("encode png: %v", err)
	}
	return buf.Bytes()
}

func makeJPEG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r := uint8(255 * y / h)
			g := uint8(120 + 60*x/w)
			b := uint8(40)
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85}); err != nil {
		log.Fatalf("encode jpeg: %v", err)
	}
	return buf.Bytes()
}

// makeFakeMP3 returns bytes that begin with a real MPEG audio frame sync so
// http.DetectContentType reports audio/mpeg. The body is not playable audio.
func makeFakeMP3() []byte {
	header := []byte{0xff, 0xfb, 0x90, 0x00}
	return append(header, randomBytes(4096)...)
}

// makeFakeWAV returns a minimally-correct RIFF/WAVE header followed by silence.
func makeFakeWAV() []byte {
	const sampleRate = 8000
	const seconds = 1
	const numSamples = sampleRate * seconds
	const dataSize = numSamples * 2

	buf := &bytes.Buffer{}
	buf.WriteString("RIFF")
	writeU32LE(buf, uint32(36+dataSize))
	buf.WriteString("WAVE")
	buf.WriteString("fmt ")
	writeU32LE(buf, 16)           // subchunk size
	writeU16LE(buf, 1)            // PCM
	writeU16LE(buf, 1)            // channels
	writeU32LE(buf, sampleRate)   // sample rate
	writeU32LE(buf, sampleRate*2) // byte rate
	writeU16LE(buf, 2)            // block align
	writeU16LE(buf, 16)           // bits/sample
	buf.WriteString("data")
	writeU32LE(buf, uint32(dataSize))
	// Soft sine so the file isn't all zeros.
	for i := 0; i < numSamples; i++ {
		v := int16(3000 * math.Sin(2*math.Pi*440*float64(i)/sampleRate))
		writeI16LE(buf, v)
	}
	return buf.Bytes()
}

func writeU16LE(b *bytes.Buffer, v uint16) { b.Write([]byte{byte(v), byte(v >> 8)}) }
func writeI16LE(b *bytes.Buffer, v int16)  { writeU16LE(b, uint16(v)) }
func writeU32LE(b *bytes.Buffer, v uint32) {
	b.Write([]byte{byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24)})
}

// makeXLSX builds a minimal but spec-valid .xlsx (Office Open XML) workbook
// containing a single sheet with a small table. Excel and LibreOffice will
// open it cleanly.
func makeXLSX() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Default Extension="xml" ContentType="application/xml"/>
<Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>
<Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>
</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>
</Relationships>`,
		"xl/_rels/workbook.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>
</Relationships>`,
		"xl/workbook.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
<sheets><sheet name="Sales" sheetId="1" r:id="rId1"/></sheets>
</workbook>`,
		"xl/worksheets/sheet1.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
<sheetData>
<row r="1"><c r="A1" t="inlineStr"><is><t>Month</t></is></c><c r="B1" t="inlineStr"><is><t>Product</t></is></c><c r="C1"><v>0</v></c></row>
<row r="2"><c r="A2" t="inlineStr"><is><t>Jan</t></is></c><c r="B2" t="inlineStr"><is><t>Widget</t></is></c><c r="C2"><v>120</v></c></row>
<row r="3"><c r="A3" t="inlineStr"><is><t>Feb</t></is></c><c r="B3" t="inlineStr"><is><t>Gadget</t></is></c><c r="C3"><v>98</v></c></row>
<row r="4"><c r="A4" t="inlineStr"><is><t>Mar</t></is></c><c r="B4" t="inlineStr"><is><t>Sprocket</t></is></c><c r="C4"><v>143</v></c></row>
</sheetData>
</worksheet>`,
	}
	for name, body := range files {
		w, err := zw.Create(name)
		if err != nil {
			log.Fatalf("zip create %s: %v", name, err)
		}
		if _, err := w.Write([]byte(body)); err != nil {
			log.Fatalf("zip write %s: %v", name, err)
		}
	}
	if err := zw.Close(); err != nil {
		log.Fatalf("zip close: %v", err)
	}
	return buf.Bytes()
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("rand: %v", err)
	}
	return b
}
