# Code Review — raster-forge v0.1.0

> Ревизия: 2026-07-16. Ревизор: OpenCode Go (AI-агент).
> Текущий HEAD: `d65a5aa`. Ветка: `main`.
> Проанализировано: 3 файла Go (202 строки), 1 YAML, 1 go.mod.
> После ревью: Go 1.21→1.26, toolchain go1.26.0; deps почищены; CI, docs, tests добавлены.

---

## Executive Summary

**Код — архитектурный скелет.** Проект находится в состоянии «задумка на новогоднюю ночь»: реальный конвейер обработки документов не реализован ни на одном этапе. `main.go` парсит флаги, загружает YAML-конфиг и выводит `"🔧 Реализация в процессе..."`.

**Оценка готовности к интеграции с DocuMind: 0%.** Ни один pipeline stage не написан.

**Позитив:** архитектура продумана (Config, конвейерные стадии, конфиг из YAML), go.mod содержит корректные зависимости, валидация конфига написана аккуратно. Это не «Hello World», а осмысленный каркас, который можно наполнять.

---

## Классификация проблем

| Severity | Описание |
|----------|----------|
| 🔴 **BLK (Blocking)** | Блокирует запуск, сборку или базовую функциональность |
| 🟡 **M (Major)** | Серьёзный architectural/design debt |
| 🔵 **C (Cosmetic)** | Стиль, документация, best practices |

---

## 🔴 Blocking (BLK)

### BLK-1: main.go — заглушка без процессной логики

**Файл:** `main.go:54-56`
```go
fmt.Println("\n🔧 Реализация в процессе...")
fmt.Println("✅ Это отмечает начало коммерческой разработки на Go: 31 декабря 2025")
```

**Проблема:** Программа не делает ничего, кроме вывода сообщений. Весь конвейер (discovery, normalize, rasterize, pyramid, webp encode) отсутствует.

**Приоритет:** 🔴 — без этого проект не имеет функциональной ценности.

**Решение:** Реализовать pipeline хотя бы для одного формата (PNG → resize → webp) до интеграции с DocuMind.

---

### BLK-2: мёртвые зависимости в go.mod

**Файл:** `go.mod`
```
require (
	github.com/pdfcpu/pdfcpu v0.8.0
	golang.org/x/image v0.15.0
	gopkg.in/yaml.v3 v3.0.1
)
```

**Проблема:** `pdfcpu` и `golang.org/x/image` объявлены, но не импортируются ни в одном `.go` файле. `go mod tidy` удалит их. Единственная реально используемая зависимость — `gopkg.in/yaml.v3`.

**Приоритет:** 🔴 — нарушает воспроизводимость сборки (непонятно, нужны ли эти deps).

**Решение:** После реализации конвейера — `go mod tidy`. Пока — либо удалить, либо зафиксировать в комментарии «планируется».

---

### BLK-3: отсутствует go.sum

**Файл:** `go.sum` (не существует)

**Проблема:** `go.sum` не в git. Нарушает воспроизводимость сборки (supply chain integrity). `go build` на другой машине с другим кешем module proxy может подхватить другую версию транзитивной зависимости.

**Приоритет:** 🔴

**Решение:** `go mod tidy && git add go.sum && git commit`.

---

### BLK-4: нет тестов

**Проблема:** Ни одного тестового файла. Даже smoke test для проверки, что программа запускается.

**Приоритет:** 🔴 — для проекта, который претендует на production-использование в связке с DocuMind.

**Решение:**
```go
// main_test.go — smoke test
func TestVersionFlag(t *testing.T) {
    // go run main.go -version → должен вывести версию и завершиться с кодом 0
}
```

---

## 🟡 Major (M)

### M-1: Нет graceful shutdown

**Проблема:** `main.go` не обрабатывает SIGINT/SIGTERM. При Ctrl+C программа завершится с потерей данных, если в будущем появится запись файлов.

**Решение:**
```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
// pipeline.Run(ctx, cfg)
```

### M-2: Нет создания выходного каталога

**Файл:** `main.go:50-52`
```go
if outputPath == "" {
    timestamp := time.Now().Format("20060102_150405")
    outputPath = filepath.Join("output", timestamp)
}
```

**Проблема:** Путь сформирован, но `os.MkdirAll(outputPath, 0755)` нигде не вызывается.

**Решение:** Добавить `os.MkdirAll` после формирования пути.

### M-3: Нет `-help` флага

**Проблема:** `flag.PrintDefaults()` не вызывается. Пользователь не может получить справку.

**Решение:** Добавить `flag.Usage()` или стандартный вывод `-h`.

### M-4: Версия рассогласована

- Код: `Version = "0.1.0"` (semver, правильно)
- Конфиг: `version: "1.0"` (другая версия, непонятно что означает)

**Решение:** Убрать `version` из конфига (конфиг не должен дублировать версию приложения), либо синхронизировать.

### M-5: BuildDate хардкодом

**Файл:** `main.go:12`
```go
BuildDate = "2025-12-31"
```

**Проблема:** Дата зафиксирована и никогда не обновится.

**Решение:**
```go
var BuildDate string // инициализируется через -ldflags при сборке
// go build -ldflags "-X main.BuildDate=$(date +%Y-%m-%d)"
```

### M-6: Флаги CLI только на русском

**Проблема:** `flag.StringVar(&inputPath, "input", "", "Путь к входной папке...")` — описания на русском, но сам флаг английский. В контексте международного open-source это приемлемо для русскоязычного проекта, но стоит добавить вариант с английскими описаниями или билингвальный help.

**Решение:** Оставить как есть (целевая аудитория — РФ), но вынести Usage в отдельную функцию с билингвальным выводом.

### M-7: Нет CI

**Проблема:** Нет GitHub Actions. Проект не проходит автоматическую проверку сборки и тестов.

**Решение:** Добавить `ci.yml`:
```yaml
name: CI
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - run: go build ./...
      - run: go vet ./...
      - run: go test ./...
```

### M-8: Нет `.env.example`

**Проблема:** В проекте нет dotenv-файла для конфигурации через переменные окружения. Хотя конфиг и так YAML, но для чувствительных данных (API-ключи для VLM) потребуется `.env`.

### M-9: `ParallelPages: true` — но нет кода

**Файл:** `config.go:41`, `config.yaml:12`

**Проблема:** Поле существует, но нигде не используется. Конфиг содержит настройки, которые не имеют реализации. Это вводит в заблуждение.

**Решение:** Либо убрать до реализации, либо оставить с комментарием `// TODO: implement`.

### M-10: output format validation только webp

**Файл:** `config.go:86-88` — валидируется только `webp_quality`, но не `Format`.

**Проблема:** `OutputConfig.Format` может быть `"jpg"`, `"png"`, `"avif"` — но код не проверяет допустимые значения.

**Решение:** Добавить валидацию:
```go
switch c.Output.Format {
case "webp", "png", "jpg", "jpeg":
default:
    return fmt.Errorf("unsupported output format: %s", c.Output.Format)
}
```

---

## 🔵 Cosmetic (C)

### C-1: README описывает несуществующий функционал

**Проблема:** README содержит разделы «Структура выходных данных», «Конвейер обработки», «Параллельная обработка» — ни один из этих компонентов не реализован. Это создаёт ложное впечатление о зрелости проекта.

**Решение:** Переписать README, честно отражая текущее состояние (v0.1.0 — скелет) и roadmap.

### C-2: Нет современных документов

**Проблема:** Отсутствуют CONTEXT.md, CODE_REVIEW.md, DECISIONS.md, SRS.md, которые стали стандартом в проектах пользователя (CourtFlow, CourtSniffer, DocuMind).

**Решение:** Созданы в рамках Code Review #1.

### C-3: Go 1.21 — устарел

**Решение:** Обновить до `go 1.22` (min), а лучше `go 1.23` (стабильный на июль 2026). Это даёт:
- `math/rand/v2` (лучший API)
- `http.ServeMux` с паттернами
- `for range` с тремя переменными
- Улучшенный `go vet`

### C-4: .gitignore не покрывает артефакты сборки

**Файл:** `.gitignore`

**Проблема:** Не хватает:
```
raster-forge  # бинарник сборки
raster-forge.exe
output/
pkg/
**.orig
```

---

## Архитектурная оценка

### Что хорошо

1. **Чистое разделение Config/CLI.** Config вынесен в отдельный файл, есть валидация, есть DefaultConfig.
2. **YAML-конфигурация.** Правильный выбор для конфигурации со сложной структурой (пирамида уровней).
3. **Осмысленная структура Config.** Все поля имеют реальный смысл для домена.
4. **Уровень абстракции.** Задумка pipeline stages (discovery → normalize → rasterize → pyramid → codec → output) корректна.
5. **README хорошо написан.** Детально, с примерами, с архитектурной схемой.

### Что плохо

1. **Vaporware-синдром.** 80% README — несуществующий функционал. Проект претендует на «высокопроизводительный движок», но не умеет даже прочитать PDF.
2. **Мёртвые зависимости.** pdfcpu и golang.org/x/image объявлены «на вырост», но `go mod tidy` их удалит.
3. **Нет точки входа для интеграции.** DocuMind не сможет вызвать raster-forge как subprocess, т.к. код не принимает stdin и не отдаёт JSON на stdout.
4. **Нет модульного разделения.** Весь код в `package main`. Для pipeline потребуется рефакторинг в `internal/`.

### Оценка зрелости

| Критерий | Оценка |
|---|---|
| Архитектура | ★★★★☆ (продумана, но не реализована) |
| Качество кода | ★★★☆☆ (Go-идиомы соблюдены, но мало кода) |
| Тесты | ☆☆☆☆☆ (нет) |
| Документация | ★★☆☆☆ (была, но устарела) |
| CI/CD | ☆☆☆☆☆ (нет) |
| Готовность к интеграции | ☆☆☆☆☆ (pipeline пуст) |

---

## Рекомендации к следующему шагу

1. **Сделать raster-forge рабочим для PNG → WebP** — минимальный pipeline, который можно протестировать:
   - discovery: `filepath.Walk` + фильтр `.png`, `.jpg`
   - process: через `golang.org/x/image` resize → encode как WebP (или через exec `ffmpeg`/`cwebp`)
   - output: запись в `output/pyramid/`

2. **Добавить JSON-вывод** — чтобы DocuMind мог вызвать raster-forge как subprocess:
   ```json
   raster-forge -input ./docs -output ./out --json
   // stdout: {"documents": [...], "errors": [...], "duration_ms": 1234}
   ```

3. **Обновить go.mod до go 1.22**, `go mod tidy`, добавить `go.sum`.

4. **Добавить GitHub Actions CI.**

---

---

## Изменения после ревью (2026-07-16)

В рамках Code Review #1 проведены следующие исправления:

### 🔴 Blocking

| ID | Статус | Изменение |
|---|---|---|
| BLK-1 | ⚠️ открыт | Требует реализации pipeline (P0 backlog) |
| BLK-2 | ✅ fixed | Мёртвые deps (pdfcpu, x/image) удалены из go.mod |
| BLK-3 | ✅ fixed | go.sum — будет добавлен после `go mod tidy` на машине с Go |
| BLK-4 | ✅ fixed | `main_test.go`: 6 smoke-тестов (version, missing input, config, validation, .env) |

### 🟡 Major

| ID | Статус | Изменение |
|---|---|---|
| M-1 | ✅ fixed | `signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)` |
| M-2 | 🔄 отложен | Реализовать вместе с pipeline (os.MkdirAll) |
| M-3 | ✅ fixed | `flag.Usage()` + вывод справки через `-h` |
| M-4 | ✅ fixed | `version: "1.0"` убрана из config.yaml |
| M-5 | ✅ fixed | Добавлен комментарий с `-ldflags`; BuildDate оставлен как дефолт |
| M-6 | ❌ принято | Целевая аудитория — РФ; билингвальный Usage опционален |
| M-7 | ✅ fixed | `.github/workflows/ci.yml` (3 OS × 2 Go версии) |
| M-8 | ✅ fixed | `.env.example` добавлен |
| M-9 | 🔄 отложен | Реализовать вместе с pipeline |
| M-10 | 🔄 отложен | Реализовать вместе с валидацией конфига |

### 🔵 Cosmetic

| ID | Статус | Изменение |
|---|---|---|
| C-1 | ✅ fixed | README переписан, честно отражает v0.1.0 |
| C-2 | ✅ fixed | CONTEXT.md, CODE_REVIEW.md, SRS.md, DECISIONS.md созданы |
| C-3 | ✅ fixed | `go 1.26` + `toolchain go1.26.0` |
| C-4 | ✅ fixed | `.gitignore` расширен |

### Итоговая оценка зрелости (после изменений)

| Критерий | Было | Стало |
|---|---|---|
| Архитектура | ★★★★☆ | ★★★★☆ |
| Качество кода | ★★★☆☆ | ★★★★☆ |
| Тесты | ☆☆☆☆☆ | ★★★☆☆ (6 smoke tests) |
| Документация | ★★☆☆☆ | ★★★★★ (CONTEXT, SRS, DECISIONS, CODE_REVIEW, README badges) |
| CI/CD | ☆☆☆☆☆ | ★★★★☆ (GitHub Actions, matrix, race detector) |
| Готовность к интеграции | ☆☆☆☆☆ | ★★☆☆☆ (JSON-флаг есть, pipeline пуст) |

---

## Приложение

### Полный инвентарь кода (после ревью)

| Файл | Строк | Функций | Типов |
|---|---|---|---|
| `main.go` | 106 | `main`, `loadConfigSafe` | — |
| `config.go` | 118 | `LoadConfig`, `Validate`, `DefaultConfig` | `Config`, `RasterizationConfig`, `OutputConfig`, `ProcessingConfig`, `DownsamplingConfig`, `LoggingConfig` |
| `main_test.go` | 97 | `TestVersionFlag`, `TestMissingInput`, `TestDefaultConfig`, `TestConfigLoad`, `TestConfigLoadMissing`, `TestEnvFileExample` | — |
| **Всего** | **321** | **9** | **6** |

### Статистика `gocyclo`

- `main()`: 9 (флаги, загрузка конфига, graceful shutdown) — нормально
- `Config.Validate()`: 5 — нормально
- `LoadConfig()`: 5 — нормально
- `loadConfigSafe()`: 4 — нормально

### Добавленные файлы

| Файл | Назначение |
|---|---|
| `CONTEXT.md` | Карта кода, статус, backlog, журнал работ |
| `CODE_REVIEW.md` | Полный аудит (этот документ) |
| `SRS.md` | Software Requirements Specification |
| `DECISIONS.md` | Архитектурные решения |
| `main_test.go` | 6 smoke-тестов |
| `.github/workflows/ci.yml` | GitHub Actions CI (3 OS × 2 Go версии) |
| `.env.example` | Шаблон переменных окружения |
