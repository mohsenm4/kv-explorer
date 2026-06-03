---
name: Run DB Benchmark
description: اجرای benchmark مقایسه‌ای بین سه دیتابیس PebbleDB، BadgerDB و LevelDB روی عملیات Set/Get/Delete/Iterate
examples:
  - "/run-benchmark"
  - "/run-benchmark --keys=1000000"
---

# اجرای Benchmark دیتابیس‌ها

این skill برای مقایسه‌ی سرعت سه دیتابیس پشتیبانی‌شده در KV-Studio استفاده می‌شود.

## مراحل اجرا

1. اطمینان از وجود پکیج‌های benchmark در `internal/databases/<x>/bench_test.go`
2. اجرای benchmarkها:
   ```bash
   go test -bench=. -benchmem -benchtime=5s ./internal/databases/...
   ```
3. جمع‌آوری خروجی و تولید جدول مقایسه‌ای در `docs/benchmarks/YYYY-MM-DD.md`
4. ثبت ns/op، B/op و allocs/op برای هر عملیات

## معیارهای مقایسه

- **Set throughput** — کلید/ثانیه برای نوشتن متوالی
- **Get latency** — میانگین زمان خواندن یک کلید
- **Iterate scan rate** — کلید/ثانیه برای پیمایش range
- **Memory overhead** — تخصیص حافظه به ازای هر عملیات

## خروجی نهایی

گزارش به‌صورت Markdown با نمودار جدول مقایسه‌ای و توصیه‌ی استفاده برای هر سناریو (read-heavy, write-heavy, mixed).
