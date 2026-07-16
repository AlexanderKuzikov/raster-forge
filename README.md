# raster-forge

[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](go.mod)
[![CI](https://github.com/AlexanderKuzikov/raster-forge/actions/workflows/ci.yml/badge.svg)](https://github.com/AlexanderKuzikov/raster-forge/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Architecture%20Skeleton-yellow)](CONTEXT.md)

**v0.1.0** — архитектурный каркас. Pipeline в разработке. [→ CONTEXT.md](CONTEXT.md) · [→ SRS.md](SRS.md)

Высокопроизводительный движок нормализации и растеризации документов с генерацией многоуровневой пирамиды разрешений для VLM-инференса.

Принимает PDF, PNG, JPG → растеризует → строит каскад DPI (300→75) → кодирует WebP.

---

## Мотивация

VLM-модели (DocuMind, Luminar) работают с изображениями, а не с PDF. raster-forge:
1. Растеризует PDF/изображения в чистые grayscale-страницы
2. Строит пирамиду разрешений (75–300 DPI)
3. Позволяет VLM выбирать уровень: 75 DPI для layout analysis, 300 DPI для детального OCR

---

## Быстрый старт

```bash
go build -o raster-forge .
./raster-forge -input ./testdata -output ./out
```

### Требования

- Go 1.24+ (toolchain go1.26.0, автоматически подтянет нужную версию)
- (Опционально) `cwebp` из libwebp для WebP-кодирования

### Флаги

| Флаг | Назначение |
|---|---|
| `-input <path>` | Входная директория (обязательный) |
| `-output <path>` | Выходная директория (по умолч. `output_YYYYMMDD_HHMMSS`) |
| `-config <path>` | Путь к YAML-конфигу (по умолч. `config.yaml`) |
| `-version` | Показать версию |
| `-json` | Вывести JSON-отчёт на stdout для интеграции |

---

## Архитектура

```
input/ → discovery → normalize → rasterize → pyramid → webp → output/
```

Подробнее: [SRS.md](SRS.md).

---

## Текущее состояние

| Компонент | Статус |
|---|---|
| CLI (флаги, config, -help, -json) | ✅ Работает |
| Config (YAML, валидация, дефолты) | ✅ Работает |
| Graceful shutdown (SIGINT/SIGTERM) | ✅ Работает |
| Smoke-тесты (6 тестов) | ✅ Работает |
| CI (GitHub Actions, 3 OS × 2 Go) | ✅ Работает |
| Документация (CONTEXT, SRS, DECISIONS, CODE_REVIEW) | ✅ Готова |
| Discovery (сканирование) | 🔄 В разработке |
| Normalize (нормализация) | 🔄 В разработке |
| Rasterize (PDF → image) | 🔄 В разработке |
| Pyramid (downsample) | 🔄 В разработке |
| WebP encoding | 🔄 В разработке |
| JSON-режим (вывод отчёта) | 🔄 В разработке |

---

## Документация

- [CONTEXT.md](CONTEXT.md) — актуальный статус, карта кода, backlog
- [CODE_REVIEW.md](CODE_REVIEW.md) — Code Review #1 (2026-07-16)
- [SRS.md](SRS.md) — Software Requirements Specification
- [DECISIONS.md](DECISIONS.md) — Архитектурные решения
- `config.yaml` — Конфигурация по умолчанию

---

## Roadmap

- [ ] **P0:** Minimal pipeline: PNG → resize → WebP
- [ ] **P0:** JSON-режим для DocuMind
- [ ] **P1:** PDF-поддержка (pdfcpu)
- [ ] **P1:** Пирамида разрешений
- [x] **P2:** CI (GitHub Actions)
- [x] **P1:** Graceful shutdown
- [ ] **P2:** Параллельная обработка
- [ ] **P2:** Тесты (smoke, unit)

---

## Лицензия

Apache License 2.0 — см. [LICENSE](LICENSE).

## Автор

Alexander Kuzikov · [GitHub](https://github.com/AlexanderKuzikov)
