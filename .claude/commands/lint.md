---
name: Lint
description: Format code and run static analysis
---

# Lint

```bash
gofmt -w . && go vet ./...
```

If `golangci-lint` is installed, run it as well:

```bash
golangci-lint run ./...
```

Report which files were reformatted and any issues that remain.
