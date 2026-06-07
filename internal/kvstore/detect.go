package kvstore

import (
	"os"
	"strings"
)

// DetectEngine returns the engine that most likely created the directory; ok=false means "use whatever the user picked".
func DetectEngine(path string) (kind EngineKind, ok bool) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return "", false
	}

	var hasKeyRegistry, hasLDB, hasManifest, hasCurrent, hasVlog, hasSST, hasPebbleOptions bool
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		switch {
		case name == "KEYREGISTRY":
			hasKeyRegistry = true
		case strings.HasPrefix(name, "OPTIONS-"):
			hasPebbleOptions = true
		case strings.HasSuffix(name, ".vlog"):
			hasVlog = true
		case strings.HasSuffix(name, ".ldb"):
			hasLDB = true
		case strings.HasSuffix(name, ".sst"):
			hasSST = true
		case strings.HasPrefix(name, "MANIFEST-"):
			hasManifest = true
		case name == "CURRENT":
			hasCurrent = true
		}
	}

	switch {
	case hasKeyRegistry || hasVlog:
		return EngineBadger, true
	case hasPebbleOptions:
		return EnginePebble, true
	case hasLDB || (hasCurrent && hasManifest && !hasSST):
		return EngineLevelDB, true
	}
	return "", false
}
