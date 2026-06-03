---
name: Run All Tests
description: Run the full test suite together with go vet and report coverage
---

# Run All Tests

```bash
go vet ./... && go test -race -cover ./...
```

If anything fails, analyze the output and propose fixes.

The final report should include pass/fail counts, coverage percentage per package, and any `go vet` warnings.
