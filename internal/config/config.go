package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// Config is the persisted user state. It lives at
// ~/.kvexplorer/config.json. Missing fields fall back to zero values.
type Config struct {
	Theme        string   `json:"theme,omitempty"`    // "light" | "dark" | "system"
	Language     string   `json:"language,omitempty"` // BCP-47 tag; "" means follow OS locale
	WindowWidth  float32  `json:"window_width,omitempty"`
	WindowHeight float32  `json:"window_height,omitempty"`
	Recents      []Recent `json:"recents,omitempty"`
}

// Recent points at a previously opened database.
type Recent struct {
	Path     string    `json:"path"`
	Engine   string    `json:"engine"`
	OpenedAt time.Time `json:"opened_at"`
}

const (
	dirName     = ".kvexplorer"
	fileName    = "config.json"
	maxRecents  = 10
	defaultPath = "~/.kvexplorer/config.json"
)

// Path returns the absolute config path (~/.kvexplorer/config.json).
func Path() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return defaultPath
	}
	return filepath.Join(home, dirName, fileName)
}

// LogDir returns the log directory path (~/.kvexplorer/logs/).
func LogDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~/.kvexplorer/logs/"
	}
	return filepath.Join(home, dirName, "logs")
}

// Load reads the config from disk. A missing file is not an error — it
// returns an empty Config.
func Load() (Config, error) {
	var c Config
	data, err := os.ReadFile(Path())
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return c, nil
		}
		return c, err
	}
	if err := json.Unmarshal(data, &c); err != nil {
		return c, err
	}
	return c, nil
}

// Save writes the config to disk, creating the directory as needed.
func Save(c Config) error {
	path := Path()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// AddRecent prepends an entry, deduplicating by path and capping length.
func (c *Config) AddRecent(path, engine string) {
	now := time.Now()
	filtered := c.Recents[:0:0]
	for _, r := range c.Recents {
		if r.Path == path {
			continue
		}
		filtered = append(filtered, r)
	}
	c.Recents = append([]Recent{{Path: path, Engine: engine, OpenedAt: now}}, filtered...)
	if len(c.Recents) > maxRecents {
		c.Recents = c.Recents[:maxRecents]
	}
}
