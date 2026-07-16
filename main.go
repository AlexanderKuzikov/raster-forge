package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// Информация о версии — обновляется при сборке:
//
//	go build -ldflags "-X main.Version=0.1.0 -X main.BuildDate=$(date +%Y-%m-%d)" -o raster-forge .
const (
	Version   = "0.1.0"
	BuildDate = "2025-12-31" // default, override via -ldflags
)

func main() {
	var (
		inputPath   string
		outputPath  string
		configPath  string
		showVersion bool
		jsonMode    bool
	)

	// Флаги CLI
	flag.StringVar(&inputPath, "input", "", "Путь к входной папке с документами (обязательный)")
	flag.StringVar(&outputPath, "output", "", "Путь к выходной папке (по умолч. output_YYYYMMDD_HHMMSS)")
	flag.StringVar(&configPath, "config", "config.yaml", "Путь к файлу конфигурации (по умолч. config.yaml)")
	flag.BoolVar(&showVersion, "version", false, "Показать информацию о версии")
	flag.BoolVar(&jsonMode, "json", false, "Вывести JSON-отчёт в stdout (для интеграции с DocuMind)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Использование: raster-forge -input <path> [опции]\n\n")
		fmt.Fprintf(os.Stderr, "Высокопроизводительный движок нормализации и растеризации документов.\n")
		fmt.Fprintf(os.Stderr, "Принимает PDF/PNG/JPG → растеризует → строит пирамиду DPI → кодирует WebP.\n\n")
		fmt.Fprintf(os.Stderr, "Обязательные флаги:\n")
		fmt.Fprintf(os.Stderr, "  -input <path>\tПуть к входной папке с документами\n\n")
		fmt.Fprintf(os.Stderr, "Опции:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nПримеры:\n")
		fmt.Fprintf(os.Stderr, "  raster-forge -input ./documents -output ./out\n")
		fmt.Fprintf(os.Stderr, "  raster-forge -input ./docs -config custom.yaml --json\n")
		fmt.Fprintf(os.Stderr, "  raster-forge -version\n")
	}
	flag.Parse()

	if showVersion {
		fmt.Printf("raster-forge v%s (built: %s)\n", Version, BuildDate)
		os.Exit(0)
	}

	if inputPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Загрузка конфигурации
	cfg, err := loadConfigSafe(configPath)
	if err != nil {
		log.Fatalf("Ошибка конфигурации: %v", err)
	}

	// Формирование выходного пути
	if outputPath == "" {
		timestamp := time.Now().Format("20060102_150405")
		outputPath = filepath.Join("output", timestamp)
	}

	fmt.Printf("raster-forge v%s — pipeline запуск\n\n", Version)
	fmt.Printf("📂 Вход:  %s\n", inputPath)
	fmt.Printf("📂 Выход: %s\n", outputPath)
	fmt.Printf("⚙️  DPI:  %d (уровни пирамиды: %v)\n", cfg.Rasterization.BaseDPI, cfg.Rasterization.PyramidLevels)
	fmt.Printf("🖼️  Формат: %s (качество: %d)\n", cfg.Output.Format, cfg.Output.WebPQuality)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// TODO: запуск pipeline.Run(ctx, cfg, inputPath, outputPath)
	fmt.Println("\n🔧 Реализация pipeline в процессе...")
	fmt.Println("   См. SRS.md и CONTEXT.md для плана разработки.")

	// Ожидание сигнала или завершения
	<-ctx.Done()
	fmt.Println("\n⏹ Завершение работы...")
}

// loadConfigSafe загружает конфиг из файла или возвращает дефолтный.
func loadConfigSafe(path string) (*Config, error) {
	if _, err := os.Stat(path); err != nil {
		log.Printf("ℹ️  Конфиг %s не найден, используется конфигурация по умолчанию", path)
		return DefaultConfig(), nil
	}
	cfg, err := LoadConfig(path)
	if err != nil {
		return nil, fmt.Errorf("загрузка %s: %w", path, err)
	}
	log.Printf("✅ Конфиг загружен из %s", path)
	return cfg, nil
}
