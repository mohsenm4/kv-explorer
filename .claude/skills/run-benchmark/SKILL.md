---
name: Run DB Benchmark
description: Run a comparative benchmark across PebbleDB, BadgerDB, and LevelDB for Set/Get/Delete/Iterate operations
examples:
  - "/run-benchmark"
  - "/run-benchmark --keys=1000000"
---

# Database Benchmark Runner

This skill compares the performance of the three supported KV-Studio databases.

## Steps

1. Verify that benchmark packages exist in `internal/databases/<x>/bench_test.go`.
2. Run the benchmarks:

   ```bash
   go test -bench=. -benchmem -benchtime=5s ./internal/databases/...
   ```

3. Collect the output and generate a comparison report at `docs/benchmarks/YYYY-MM-DD.md`.
4. Record `ns/op`, `B/op`, and `allocs/op` for each operation.

## Comparison Metrics

- **Set throughput** — keys per second for sequential writes
- **Get latency** — average single-key read time
- **Iterate scan rate** — keys per second for range scans
- **Memory overhead** — allocations per operation

## Final Output

A Markdown report containing a comparison table and a usage recommendation
per scenario (read-heavy, write-heavy, mixed).
