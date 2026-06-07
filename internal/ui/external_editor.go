package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"github.com/fsnotify/fsnotify"
)

// Threshold above which we confirm before writing the value to a temp file: VSCode and the round-trip both struggle past this.
const externalSizeWarn = 50 << 20

type externalEditSession struct {
	path      string
	dir       string
	watcher   *fsnotify.Watcher
	closeOnce sync.Once
}

func (es *externalEditSession) Path() string { return es.path }

// Close stops the watcher and removes the temp dir. Safe to call multiple times.
func (es *externalEditSession) Close() {
	if es == nil {
		return
	}
	es.closeOnce.Do(func() {
		if es.watcher != nil {
			_ = es.watcher.Close()
		}
		if es.dir != "" {
			_ = os.RemoveAll(es.dir)
		}
	})
}

func startExternalEditSession(key, value []byte, onChange func()) (*externalEditSession, error) {
	dir, err := os.MkdirTemp("", "kvexplorer-")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	name := suggestedExportName(key, value)
	full := filepath.Join(dir, name)
	if err := os.WriteFile(full, value, 0o600); err != nil {
		_ = os.RemoveAll(dir)
		return nil, fmt.Errorf("write temp file: %w", err)
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		_ = os.RemoveAll(dir)
		return nil, fmt.Errorf("create watcher: %w", err)
	}
	// Watch the parent directory so atomic saves (write-temp-then-rename, used by VSCode) are seen.
	if err := w.Add(dir); err != nil {
		_ = w.Close()
		_ = os.RemoveAll(dir)
		return nil, fmt.Errorf("watch dir: %w", err)
	}

	go watchExternal(w, name, onChange)

	es := &externalEditSession{path: full, dir: dir, watcher: w}
	if err := launchVSCode(full); err != nil {
		return es, err
	}
	return es, nil
}

// Debounced: rapid events from a single save coalesce into one onChange call.
func watchExternal(w *fsnotify.Watcher, targetName string, onChange func()) {
	var pending atomic.Bool
	fire := func() {
		if pending.Swap(false) {
			fyne.Do(onChange)
		}
	}
	for event := range w.Events {
		if filepath.Base(event.Name) != targetName {
			continue
		}
		if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
			pending.Store(true)
			time.AfterFunc(150*time.Millisecond, fire)
		}
	}
}

func launchVSCode(file string) error {
	if runtime.GOOS == "darwin" {
		// Bundle ID is set by VSCode at install time, survives renames and any second install (e.g. "Visual Studio Code 2.app").
		for _, bundle := range []string{"com.microsoft.VSCode", "com.microsoft.VSCodeInsiders"} {
			if err := exec.Command("open", "-b", bundle, file).Run(); err == nil {
				return nil
			}
		}
		return exec.Command("open", "-a", "Visual Studio Code", file).Run()
	}
	if path, err := exec.LookPath("code"); err == nil {
		return exec.Command(path, file).Start()
	}
	return fmt.Errorf("VSCode 'code' command not found in PATH")
}
