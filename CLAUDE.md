# KV-Studio

ابزار گرافیکی دسکتاپ برای مدیریت و بررسی پایگاه‌داده‌های کلید-مقدار (Key-Value Stores).
این پروژه بازنویسی کامل و مدرن «KV-Toolbox» با کمک هوش مصنوعی است.

## دیتابیس‌های پشتیبانی‌شده

- **PebbleDB** — موتور سطح پایین مبتنی بر LSM-Tree (مشتق‌شده از RocksDB)
- **BadgerDB** — KV store پراستفاده در پروژه‌های Go-native
- **LevelDB** — کلاسیک، ساده، سبک

## معماری

پروژه از الگوی استاندارد Go با تفکیک `cmd/` و `internal/` پیروی می‌کند:

```
kv-studio/
├── cmd/kvstudio/           # نقطه‌ی ورود برنامه (main package)
├── internal/
│   ├── databases/          # آداپتورهای هر دیتابیس
│   │   ├── pebble/         # پیاده‌سازی برای PebbleDB
│   │   ├── badger/         # پیاده‌سازی برای BadgerDB
│   │   └── leveldb/        # پیاده‌سازی برای LevelDB
│   ├── ui/                 # لایه‌ی رابط کاربری (Fyne)
│   │   ├── mainwindow/     # پنجره‌ی اصلی
│   │   ├── components/     # ویجت‌های قابل استفاده‌ی مجدد
│   │   └── theme/          # تم و استایل
│   ├── logic/              # منطق کسب‌وکار (filter, search, ...)
│   ├── config/             # خواندن/نوشتن تنظیمات
│   └── utils/              # توابع کمکی عمومی
├── docs/                   # مستندات فنی پروژه
└── .claude/                # پیکربندی Claude Code (skills, agents, settings)
```

## اصل کلیدی معماری: یک Interface برای همه‌ی دیتابیس‌ها

هر دیتابیس باید interface مشترکی به نام `KVStore` را پیاده‌سازی کند:

```go
type KVStore interface {
    Open(path string) error
    Close() error
    Get(key []byte) ([]byte, error)
    Set(key, value []byte) error
    Delete(key []byte) error
    Iterate(prefix []byte, fn func(key, value []byte) bool) error
    Stats() Stats
}
```

این الگو باعث می‌شود لایه‌ی UI و logic از پیاده‌سازی هر دیتابیس کاملاً مستقل باشد.

## دستورات اصلی

| دستور | توضیح |
|---|---|
| `go build ./cmd/kvstudio` | بیلد باینری |
| `go run ./cmd/kvstudio` | اجرای برنامه در حالت توسعه |
| `go test ./...` | اجرای کل تست‌ها |
| `go vet ./...` | بررسی استاتیک کد |
| `gofmt -w .` | فرمت‌بندی |

## استانداردهای کدنویسی

- **زبان**: Go 1.22+
- **رابط گرافیکی**: Fyne v2
- **ساختار خطاها**: استفاده از `errors.Is`/`errors.As` و wrapping با `fmt.Errorf("...: %w", err)`
- **نام‌گذاری**: UpperCamelCase برای exported، lowerCamelCase برای داخلی
- **تست**: هر پکیج `internal/...` باید فایل `*_test.go` همراه داشته باشد
- **کامنت**: فقط جایی که «چرا» غیرواضح است — نه برای توضیح «چه».
- **هیچ `panic` در مسیر اصلی** — همه‌ی خطاها باید propagate شوند تا UI تصمیم بگیرد.

## قراردادهای پروژه

1. هر آداپتور دیتابیس در پکیج جداگانه‌ی خود (`internal/databases/<name>`) قرار می‌گیرد.
2. UI نباید مستقیماً به packageهای دیتابیس import داشته باشد — همیشه از طریق interface.
3. تنظیمات کاربر در `~/.kvstudio/config.json` ذخیره می‌شود.
4. لاگ‌ها در `~/.kvstudio/logs/` با rotation روزانه نوشته می‌شوند.
5. هیچ secret یا path خاص ماشین در ریپو commit نمی‌شود.

## وضعیت فعلی

این پروژه در فاز **bootstrap** قرار دارد. ساختار پایه ایجاد شده و توسعه‌ی modules اصلی هنوز شروع نشده است. اولین گام، پیاده‌سازی interface و سه آداپتور دیتابیس است.

## برای Claude

- پیش از تغییرات معماری بزرگ، **حالت plan** را فعال کن.
- پیش از هر commit، اطمینان حاصل کن `go vet` و `go test ./...` با موفقیت می‌گذرند.
- وقتی فایلی در `internal/databases/<x>/` تغییر می‌کند، آداپتورهای دیگر را هم بررسی کن که interface هنوز یکدست است.
- برای کارهای تخصصی، از subagentهای تعریف‌شده در `.claude/agents/` استفاده کن.
