# testdata — тестовые файлы для raster-forge

## Структура

```
testdata/
├── single-page/       # Одностраничные документы
│   ├── document1.pdf
│   ├── scan1.png
│   └── scan2.jpg
├── multi-page/        # Многостраничный PDF
│   └── contract.pdf
├── multi-file/        # Многофайловый документ
│   ├── page001.png
│   ├── page002.png
│   └── page003.png
└── mixed/             # Смешанный документ
    ├── cover.pdf
    ├── page002.png
    └── page003.jpg
```

## Генерация

Тестовые файлы создаются скриптом `testdata/generate.go` (TODO).
Для ручного теста достаточно одного PNG 100x100:

```bash
go run golang.org/x/tools/cmd/stringer 2>/dev/null
# или просто создайте вручную:
#   go run testdata/gen.go
```
