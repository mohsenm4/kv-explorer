# KV-Studio

ابزار گرافیکی دسکتاپ برای مدیریت، بررسی و مقایسه‌ی پایگاه‌داده‌های کلید-مقدار.
بازنویسی کامل و مدرن KV-Toolbox با کمک هوش مصنوعی.

## دیتابیس‌های پشتیبانی‌شده

- PebbleDB
- BadgerDB
- LevelDB

## نیازمندی‌ها

- Go 1.22 یا بالاتر
- CGO فعال (برای LevelDB)
- macOS / Linux / Windows

## شروع سریع

```bash
git clone <repo-url>
cd kv-studio
go mod tidy
go run ./cmd/kvstudio
```

## ساختار پروژه

ساختار کامل و قراردادها در [CLAUDE.md](./CLAUDE.md) توضیح داده شده است.

## توسعه

- اجرای تست‌ها: `go test ./...`
- بیلد: `go build ./cmd/kvstudio`
- بیلد چندسکویی: مراجعه به skill `/build-cross-platform`

## مجوز

TBD
