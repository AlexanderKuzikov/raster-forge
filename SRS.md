# SRS — raster-forge

## 1. Назначение

Высокопроизводительный CLI-движок нормализации и растеризации документов с генерацией многоуровневой пирамиды разрешений для VLM-инференса. Вход: PDF, PNG, JPG (одиночные файлы или папки). Выход: растрированный PDF при 300 DPI + каскад WebP-изображений на уровнях 300/250/200/150/100/75 DPI.

### 1.1. Место в экосистеме

raster-forge подготавливает документы для VLM-моделей (DocuMind). Вместо того чтобы передавать оригинальный PDF (с текстовым слоем, шрифтами, сжатием), pipeline растеризует всё в чистые изображения и строит пирамиду — VLM может выбирать оптимальный уровень разрешения для своих задач.

---

## 2. Функциональные требования

### 2.1. Discovery (сканирование входа)

**Вход:** директория или список файлов.
**Процесс:**
- Рекурсивный обход (`filepath.Walk`)
- Фильтр расширений: `.pdf`, `.png`, `.jpg`, `.jpeg`, `.tiff`, `.bmp` (опционально)
- Многофайловый документ = директория, содержащая ≥1 файл документа (не является вложенной директорией с документами)
- Одиночный файл = документ из 1 страницы (изображение) или N страниц (PDF)

**Выход:** `[]Document` — список документов с типом, страницами, путями.

### 2.2. Normalize (нормализация)

**Процесс:**
- Определение типа документа (PDF / image / mixed)
- Сборка многофайлового документа: сортировка страниц по имени файла (natural sort)
- Если mixed-документ (PDF + изображения) — порядок определяется конфигурацией

**Выход:** `[]NormalizedDocument` — каждый документ имеет упорядоченный список страниц с путями к исходным файлам или страницам PDF.

### 2.3. Rasterize (растеризация)

**Зависимость:** `github.com/pdfcpu/pdfcpu` (для PDF → image conversion).

**Процесс:**
- Для PDF: рендеринг каждой страницы в grayscale (или RGB, по конфигу) при `BaseDPI`
- Для изображений: загрузка через `golang.org/x/image` или `image/png`/`image/jpeg`, интерпретация как 1 страница

**Выход:** `[]RasterizedDocument` — каждая страница — `image.Image` (или `*image.Gray`) при заданном DPI.

**Требования:**
- Использовать `context.Context` для отмены
- Параллельная обработка страниц (опционально, `ParallelPages`)

### 2.4. Pyramid (генерация пирамиды)

**Процесс:**
- Для каждого уровня `PyramidLevels` (отсортировать по убыванию: 300 → 250 → ... → 75):
  - Downsample исходного изображения при `BaseDPI` до целевого DPI
  - Алгоритм даунсамплинга: `lanczos3` (дефолт), `bicubic`, `bilinear`, `nearest`

**Выход:** `[]PyramidLevel` — для каждого документа: набор страниц × уровней.

**Требования:**
- Кэширование: если BaseDPI = 300, а уровень = 150, то downsample 300 → 150, а не оригинал → 150
- Ресайз через `golang.org/x/image/draw` (Lanczos, CatmullRom, etc.)

### 2.5. Codec (WebP-кодирование)

**Зависимость:** Внешний `cwebp` (из libwebp) или pure-Go библиотека (`golang.org/x/image/webp` — только декодирование, для кодирования нужна `github.com/kolesa-team/go-webp` или вызов `cwebp`).

**Процесс:**
- Сохранение каждого уровня пирамиды как WebP с заданным качеством

**Выход:** Файлы `.webp` в структуре каталогов (см. п. 2.6).

### 2.6. Output (выходная структура)

```
output_YYYYMMDD_HHMMSS/
├── pdfs/
│   ├── document1.pdf    # Растрирован @ 300 DPI, grayscale, без текстового слоя
│   └── document2.pdf
├── pyramid/
│   ├── document1/
│   │   ├── 300dpi/
│   │   │   └── page001.webp
│   │   ├── 250dpi/
│   │   │   └── page001.webp
│   │   ├── ...
│   │   └── 75dpi/
│   │       └── page001.webp
│   └── document2/
└── report.json            # Метаданные: документы, страницы, ошибки, время
```

---

## 3. Интерфейсы

### 3.1. CLI

```bash
raster-forge -input <path> [-output <path>] [-config <path>] [-version] [--json]
```

Флаги:
| Флаг | Тип | Дефолт | Описание |
|---|---|---|---|
| `-input` | string | (обязательный) | Путь к входной папке с документами |
| `-output` | string | `output_<timestamp>` | Путь к выходной папке |
| `-config` | string | `config.yaml` | Путь к файлу конфигурации YAML |
| `-version` | bool | false | Вывести версию |
| `-json` | bool | false | Вывести JSON-отчёт на stdout (для интеграции) |

### 3.2. JSON-режим (для интеграции с DocuMind)

При `-json`:
```json
{
  "version": "0.1.0",
  "input": "./docs",
  "output": "./out_20260716_191700",
  "duration_ms": 45230,
  "documents": [
    {
      "name": "contract.pdf",
      "pages": 5,
      "levels": [300, 250, 200, 150, 100, 75],
      "errors": null
    }
  ],
  "errors": [],
  "stats": {
    "total_pages": 12,
    "total_bytes": 45_000_000,
    "bytes_by_level": {
      "300": 15_000_000,
      "75": 800_000
    }
  }
}
```

---

## 4. Нефункциональные требования

### 4.1. Производительность
- Параллельная обработка документов: `MaxGoroutines` горутин
- Потребление памяти на документ: не более `MemoryLimitMB` (soft limit)
- Пирамида строится от большего к меньшему — используется кэш предыдущего уровня

### 4.2. Надёжность
- Ошибка одного документа не прерывает обработку остальных
- Атомарная запись: tmp-файл + rename
- Graceful shutdown по SIGINT/SIGTERM с завершением текущего документа

### 4.3. Конфигурация
- YAML-файл с валидацией всех полей
- Поддержка переменных окружения: `RASTER_FORGE_CONFIG`, `RASTER_FORGE_INPUT`, `RASTER_FORGE_OUTPUT`

### 4.4. Совместимость
- Go 1.24+ (рекомендуется 1.26)
- Windows 10/11, Linux (Ubuntu 22.04+)
- Для WebP-кодирования: внешний `cwebp` или pure-Go библиотека

---

## 5. Структура пакетов (рекомендуемая)

```
raster-forge/
├── main.go                 # CLI, флаги, инициализация
├── config.go               # Config, LoadConfig, Validate
├── config.yaml             # Пример конфига
├── internal/
│   ├── discovery/
│   │   └── discovery.go    # Walk, filter, classify
│   ├── normalize/
│   │   └── normalize.go    # Sort pages, group multi-file
│   ├── rasterize/
│   │   ├── rasterize.go    # Rasterizer interface
│   │   └── pdfcpu.go       # pdfcpu implementation
│   ├── pyramid/
│   │   └── pyramid.go      # Downsample chain
│   ├── codec/
│   │   └── webp.go         # WebP encoder
│   └── output/
│       └── output.go       # Write structure, atomic write, report.json
├── testdata/               # Тестовые файлы
│   ├── sample.pdf
│   └── sample.png
├── go.mod / go.sum
└── .github/workflows/ci.yml
```

---

## 6. Критерии приёмки (для каждой стадии)

| Стадия | Критерий |
|---|---|
| Discovery | `-input ./testdata` находит все PDF и изображения, игнорирует `.md`, `.txt` |
| Normalize | Папка с 3 PNG → 1 документ на 3 страницы в правильном порядке |
| Rasterize | PDF 5 страниц → 5 `image.Image` @ 300 DPI |
| Pyramid | 300 → 75 DPI, файлы `page001.webp` на каждом уровне |
| Output | `output_*/pyramid/doc/300dpi/page001.webp` существует и читается |
| JSON | `-json` выдаёт валидный JSON с документами и статистикой |
| Graceful | SIGTERM → процесс завершается в течение 5 сек без loss данных |
| CI | `go build ./...`, `go vet ./...`, `go test ./...` проходят |
