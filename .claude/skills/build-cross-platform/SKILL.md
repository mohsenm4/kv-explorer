---
name: Build Cross Platform
description: Build KV-Explorer binaries for macOS (Intel/ARM), Linux, and Windows in one pass
examples:
  - "/build-cross-platform"
  - "/build-cross-platform --release"
---

# Cross-Platform Build

## Prerequisites

- `fyne-cross` installed (handles CGO and GUI cross-compilation):

  ```bash
  go install github.com/fyne-io/fyne-cross@latest
  ```

## Steps

1. Clean the build directory:

   ```bash
   rm -rf ./build
   ```

2. Build for each platform:

   ```bash
   fyne-cross darwin  -arch=amd64,arm64 ./cmd/kvexplorer
   fyne-cross linux   -arch=amd64        ./cmd/kvexplorer
   fyne-cross windows -arch=amd64        ./cmd/kvexplorer
   ```

3. Output binaries land in `fyne-cross/bin/<platform>/`.

4. Generate checksums for every artifact:

   ```bash
   shasum -a 256 fyne-cross/bin/**/* > build/checksums.txt
   ```

## Notes

- Building macOS ARM from an Intel host requires Docker (`fyne-cross` uses it).
- The version string lives in `cmd/kvexplorer/main.go` (`var version`) and is
  injected at build time via `-ldflags "-X main.version=<tag>"`. The release
  workflow (`.github/workflows/release.yml`) sets it from the git tag — no
  source edit is needed for a tagged release.
