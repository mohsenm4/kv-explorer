---
name: Generate Test Data
description: Populate any of the three databases with realistic synthetic key/value datasets
examples:
  - "/generate-test-data"
  - "/generate-test-data --count=100000 --db=badger"
---

# Test Data Generator

## Purpose

Produce reproducible datasets used for UI testing, performance work, and filter behavior validation.

## Data Patterns

1. **Sequential keys** — `user:0001` through `user:NNNN`
2. **Random keys** — UUIDv4 keys with random JSON values
3. **Hierarchical keys** — `org/<id>/user/<id>/profile`
4. **Large values** — binary blobs from 1KB to 1MB

## Execution

Use the built-in generator:

```bash
go run ./cmd/kvstudio --generate-test-data \
    --db=<pebble|badger|leveldb> \
    --pattern=<sequential|random|hierarchical|large> \
    --count=<N> \
    --output=./testdata/<db>-<pattern>.db
```

## Notes

- Test data lives under `testdata/` (ignored by git).
- Keep the random seed fixed so results stay reproducible across runs.
