---
name: Generate Test Data
description: تولید مجموعه‌داده‌ی تستی برای پر کردن هر یک از سه دیتابیس با کلید/مقدارهای واقع‌گرایانه
examples:
  - "/generate-test-data"
  - "/generate-test-data --count=100000 --db=badger"
---

# تولید داده‌ی تست

## هدف

ساخت دیتاست‌های قابل تکرار برای تست UI، performance و رفتار filter.

## الگوهای داده

1. **Sequential keys** — `user:0001` تا `user:NNNN`
2. **Random keys** — UUIDv4 با مقدار JSON تصادفی
3. **Hierarchical keys** — `org/<id>/user/<id>/profile`
4. **Large values** — مقادیر باینری ۱KB تا ۱MB

## اجرا

اسکریپت داخلی برای تولید:
```bash
go run ./cmd/kvstudio --generate-test-data \
    --db=<pebble|badger|leveldb> \
    --pattern=<sequential|random|hierarchical|large> \
    --count=<N> \
    --output=./testdata/<db>-<pattern>.db
```

## نکته

- داده‌ی تست در `testdata/` ذخیره می‌شود (در `.gitignore` هست).
- seed تصادفی همیشه ثابت بمونه تا نتایج قابل تکرار باشن.
