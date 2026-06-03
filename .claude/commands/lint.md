---
name: Lint
description: فرمت‌بندی + بررسی استاتیک کد
---

# Lint

```bash
gofmt -w . && go vet ./...
```

اگر `golangci-lint` نصب باشد، آن را هم اجرا کن:
```bash
golangci-lint run ./...
```

گزارش بده چه فایل‌هایی فرمت شدند و چه issuesی مانده است.
