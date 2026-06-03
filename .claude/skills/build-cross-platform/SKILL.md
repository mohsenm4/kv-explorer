---
name: Build Cross Platform
description: بیلد همزمان باینری‌های KV-Studio برای macOS (Intel/ARM)، Linux و Windows
examples:
  - "/build-cross-platform"
  - "/build-cross-platform --release"
---

# بیلد چندسکویی

## پیش‌نیاز

- نصب بودن `fyne-cross` (برای CGO و GUI cross-compile):
  ```bash
  go install github.com/fyne-io/fyne-cross@latest
  ```

## مراحل

1. پاک‌سازی دایرکتوری بیلد:
   ```bash
   rm -rf ./build
   ```

2. بیلد برای هر سکو:
   ```bash
   fyne-cross darwin -arch=amd64,arm64 ./cmd/kvstudio
   fyne-cross linux  -arch=amd64        ./cmd/kvstudio
   fyne-cross windows -arch=amd64        ./cmd/kvstudio
   ```

3. خروجی‌ها در `fyne-cross/bin/<platform>/` ذخیره می‌شوند.

4. تولید checksum برای هر باینری:
   ```bash
   shasum -a 256 fyne-cross/bin/**/* > build/checksums.txt
   ```

## نکته

- بیلد macOS ARM روی Intel سیستم نیاز به Docker دارد (fyne-cross استفاده می‌کند).
- پیش از انتشار، نسخه را در `internal/config/version.go` به‌روز کن.
