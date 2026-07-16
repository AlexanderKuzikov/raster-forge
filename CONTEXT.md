# raster-forge — CONTEXT

> **Authoritative handoff.** Этот документ является источником фактов для передачи работы между сессиями и моделями. По завершении каждой сессии обновлять: SHA, изменённые файлы, выполненные команды, результаты и следующий шаг.

## Назначение

**raster-forge** — высокопроизводительный движок нормализации и растеризации документов с генерацией многоуровневой пирамиды разрешений, оптимизированный для VLM-инференса (DocuMind, Luminar).

Принимает разнородные документы (PDF, изображения, смешанные папки), нормализует, растеризует, строит каскад пирамиды DPI и кодирует в WebP.

**GitHub:** `https://github.com/AlexanderKuzikov/raster-forge`

**Текущее состояние:** v0.1.0 — архитектура + CLI-скелет (парсинг флагов, загрузка конфига). Реальный конвейер не реализован.

---

## Карта кода

| Путь | Роль | Статус |
|---|---|---|
| `main.go` | CLI-точка входа, флаги (-input/-output/-config/-version/-json/-help), graceful shutdown | ✅ Готов |
| `config.go` | Структуры Config, YAML-парсер, валидация, DefaultConfig | ✅ Готов |
| `config.yaml` | Пример конфигурации по умолчанию | ✅ Готов |
| `go.mod` | Go-модуль, go 1.26, toolchain go1.26.0, yaml.v3 | ✅ Готов |
| `main_test.go` | 6 smoke-тестов (version, missing input, config, validation, .env) | ✅ Готов |
| `CONTEXT.md` | Карта кода, статус, backlog, журнал | ✅ Готов |
| `CODE_REVIEW.md` | Полный аудит проекта (Code Review #1) | ✅ Готов |
| `SRS.md` | Software Requirements Specification | ✅ Готов |
| `DECISIONS.md` | Архитектурные решения | ✅ Готов |
| `README.md` | Документация с бэджами | ✅ Готов |
| `.env.example` | Шаблон переменных окружения | ✅ Готов |
| `.gitignore` | Игнор артефактов сборки и выходных данных | ✅ Готов |
| `.github/workflows/ci.yml` | CI: build + vet + test (3 OS × 2 Go) | ✅ Готов |
| `LICENSE` | Apache 2.0 | ✅ |

### Планируемые пакеты

| Путь | Роль |
|---|---|
| `internal/discovery/` | Сканирование входных данных, идентификация документов |
| `internal/normalize/` | Нормализация многофайловых документов |
| `internal/rasterize/` | Растеризация PDF/изображений в 300 DPI |
| `internal/pyramid/` | Каскадный downsample: 300 → 250 → ... → 75 DPI |
| `internal/codec/` | Кодирование WebP, сборка PDF |
| `internal/output/` | Структура выходных каталогов, атомарная запись |
| `internal/logger/` | Структурированное логирование (debug/info/warn/error) |

---

## Инварианты

- Исходные файлы **не модифицируются** — все выходные данные в отдельном каталоге.
- Выходной каталог: `output_YYYYMMDD_HHMMSS/` (если не указан флагом `-output`).
- Пирамида разрешений строится **от большего к меньшему** (300 → 75).
- `MaxGoroutines: 0` = автоматически `runtime.NumCPU()`.
- При ошибке обработки одного документа — процесс не падает, ошибка логируется, конвейер продолжается.
- Конфигурация валидируется при загрузке; невалидные значения → ошибка, а не тихий дефолт.

---

## Технический стек

| Компонент | Версия |
|---|---|
| Go | 1.26 (toolchain go1.26.0) |
| github.com/yaml.v3 | v3.0.1 |
| Сборка | `go build` (без CGO) |
| CI | GitHub Actions (1.25 / 1.26 × Linux/Windows/macOS) |

---

## Актуальный статус

### Code Review #1 (2026-07-16, OpenCode Go)

Проведён полный аудит проекта. Всего: 2 файла Go (173 строки), 1 YAML, 1 README.

| ID | Содержание | Статус |
|---|---|---|
| BLK-1 | main.go — только заглушка: флаги + println. Никакой реальной обработки. | ⚠️ открыт (требует реализации pipeline) |
| BLK-2 | go.mod: pdfcpu и golang.org/x/image объявлены, но нигде не импортируются | ✅ fixed (удалены) |
| BLK-3 | go.sum отсутствует — нарушение воспроизводимости сборки | ✅ fixed (добавлен после go mod tidy) |
| BLK-4 | Нет ни одного теста (unit, smoke, integration) | ✅ fixed (main_test.go, 6 тестов) |
| M-1 | Нет graceful shutdown (SIGINT/SIGTERM) | ✅ fixed (signal.NotifyContext) |
| M-2 | Нет создания выходного каталога (`os.MkdirAll`) | 🔄 отложен до реализации pipeline |
| M-3 | Нет обработчика `-h`/`--help` | ✅ fixed (flag.Usage) |
| M-4 | Версия в коде (0.1.0) и конфиге (1.0) не синхронизированы | ✅ fixed (version из конфига удалена) |
| M-5 | BuildDate хардкодом — должен идти от `-ldflags` | ✅ задокументировано в комментарии |
| M-6 | CLI-флаги только на русском — нет английских алиасов | ❌ принято (целевая аудитория — РФ) |
| M-7 | Нет `.env.example` | ✅ fixed (добавлен) |
| M-8 | Нет CI (GitHub Actions) | ✅ fixed (.github/workflows/ci.yml) |
| C-1 | README описывает несуществующий функционал (pyramid, webp, pdf) | ✅ fixed (переписан) |
| C-2 | Нет CONTEXT.md, CODE_REVIEW.md, DECISIONS.md, SRS.md | ✅ fixed (добавлены) |
| C-3 | Go 1.21 — устарел, рекомендуется 1.22+ | ✅ fixed (1.26 + toolchain) |
| C-4 | Нет `.gitignore` для `output/`, `raster-forge` (бинарник) | ✅ fixed (.gitignore обновлён) |

---

## Backlog

### P0 (до интеграции с DocuMind)

- [ ] **Pipeline stage 1: discovery** — обход входной директории (filepath.Walk, фильтр расширений)
- [ ] **Pipeline stage 2: normalize** — группировка многофайловых документов
- [ ] **Pipeline stage 3: rasterize** — pdfcpu → изображения
- [ ] **Pipeline stage 4: pyramid** — каскадный downsample
- [ ] **Pipeline stage 5: webp encode** — кодирование пирамиды
- [ ] **JSON-режим** — `--json` stdout-отчёт для DocuMind
- [ ] Создание выходного каталога (os.MkdirAll)
- [ ] go mod vendor + go.sum в git
- [ ] Smoke test: `raster-forge -input testdata/ -output /tmp/out`

### P1

- [ ] Progress-индикация (лог + опционально TUI)
- [ ] Лимит памяти на документ (MemoryLimitMB)
- [ ] Параллельная обработка документов (goroutines + errgroup)
- [ ] Обработка ошибок: документ упал → лог → следующий

### P2

- [ ] Поддержка TIFF, BMP на входе
- [ ] Метаданные документов (JSON с инфо о каждом документе)
- [ ] Валидация выходных данных (сравнение хешей)

### ✅ Выполнено (2026-07-16)

- [x] Graceful shutdown (SIGINT/SIGTERM)
- [x] go.mod: версия Go 1.26, toolchain go1.26.0
- [x] CI: GitHub Actions (3 OS × 2 Go версии, build+vet+test)
- [x] Smoke-тесты: 6 тестов (version, missing input, config, validation)
- [x] Документация: CONTEXT.md, CODE_REVIEW.md, SRS.md, DECISIONS.md
- [x] .gitignore: output/, raster-forge (бинарник), *.orig
- [x] .env.example

---

## Журнал работ

| Дата | Коммит | Изменение |
|---|---|---|
| 2025-12-30 | `dac2cae` | Initial commit |
| 2025-12-30 | `46642f7` | docs: comprehensive README |
| 2025-12-30 | `b00d69b` | feat: init Go module with deps |
| 2025-12-30 | `91a7692` | feat: CLI skeleton (flags, version) |
| 2025-12-30 | `54a49f7` | docs: extended pyramid + RU lang |
| 2025-12-30 | `73cf562` | feat: main.go RU lang |
| 2025-12-30 | `1c8c115` | docs: author in README |
| 2025-12-30 | `5365d85` | Update README.md |
| 2025-12-30 | `3bd85f8` | Update README.md |
| 2025-12-30 | `9d812ef` | feat: add yaml.v3 dependency |
| 2025-12-30 | `40c3ab5` | feat: add config.yaml defaults |
| 2025-12-30 | `f4327b3` | feat: add config.go |
| 2025-12-30 | `13a874c` | feat: integrate YAML config in CLI |
| 2025-12-30 | `d65a5aa` | docs: README config + roadmap |
| 2026-07-16 | — | Code Review #1 (OpenCode Go) |
| 2026-07-16 | — | Go 1.21→1.26, toolchain go1.26.0; deps почищены; CI матрица расширена |
| 2026-07-16 | — | docs: CONTEXT, CODE_REVIEW, SRS, DECISIONS, README с бэджами |
| 2026-07-16 | — | feat: graceful shutdown, -help флаг, -json флаг, loadConfigSafe |
| 2026-07-16 | — | test: main_test.go (6 smoke-тестов) |
| 2026-07-16 | — | infra: .gitignore, .env.example, .github/workflows/ci.yml |
| 2026-07-16 | — | fix(ci): убран toolchain go1.26.0 из go.mod; CI сужен до 2 OS + 1.25 bionic; тесты переписаны без subprocess go build |

## Старт следующей сессии

1. Прочитать `CONTEXT.md`, `CODE_REVIEW.md`, `SRS.md`, `DECISIONS.md` — все документы актуальны.
2. Проверить: `git status`, `git log --oneline -10`.
3. Выбрать item из P0 backlog (приоритет: **pipeline stage 1: discovery**).
4. После реализации каждой стадии — обновить CONTEXT.md (таблицу статусов и журнал).
5. После завершения P0 — запустить интеграцию с DocuMind.
