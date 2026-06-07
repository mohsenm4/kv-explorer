package main

import "github.com/mohsenm4/kv-explorer/internal/ui"

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	ui.Run(version)
}
