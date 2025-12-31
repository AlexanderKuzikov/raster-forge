package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Информация о версии
const (
	Version   = "0.1.0"
	BuildDate = "2025-12-31"
)

func main() {
	var (
		inputPath    string
		outputPath   string
		configPath   string
		showVersion  bool
	)

	flag.StringVar(&inputPath, "input", "", "Путь к входной папке с документами")
	flag.StringVar(&outputPath, "output", "", "Путь к выходной папке для обработанных документов")
	flag.StringVar(&configPath, "config", "config.yaml", "Путь к файлу конфигурации (по умолчанию: config.yaml)")
	flag.BoolVar(&showVersion, "version", false, "Показать информацию о версии")
	flag.Parse()

	if showVersion {
		fmt.Printf("🚀 raster-forge v%s (собрано: %s)\n", Version, BuildDate)
		fmt.Println("Высокопроизводительный движок нормализации и растеризации документов")
		os.Exit(0)
	}

	if inputPath == "" {
		log.Fatal("Ошибка: требуется указать входной путь. Используйте флаг -input.")
	}

	// Загрузка конфигурации
	var cfg *Config
	var err error
	
	if _, statErr := os.Stat(configPath); statErr == nil {
		cfg, err = LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Ошибка загрузки конфигурации из %s: %v", configPath, err)
		}
		fmt.Printf("✅ Загружена конфигурация из: %s\n", configPath)
	} else {
		cfg = DefaultConfig()
		fmt.Println("ℹ️  Используется конфигурация по умолчанию")
	}

	// Формирование выходного пути
	if outputPath == "" {
		timestamp := time.Now().Format("20060102_150405")
		outputPath = filepath.Join("output", timestamp)
	}

	fmt.Printf("\n📂 Входная папка:  %s\n", inputPath)
	fmt.Printf("📂 Выходная папка: %s\n", outputPath)
	fmt.Printf("⚙️  Базовое разрешение: %d DPI\n", cfg.Rasterization.BaseDPI)
	fmt.Printf("📊 Уровни пирамиды: %v DPI\n", cfg.Rasterization.PyramidLevels)
	fmt.Printf("🖼️  Формат изображений: %s\n", cfg.Output.Format)
	fmt.Printf("🗜️  Качество WebP: %d\n", cfg.Output.WebPQuality)
	fmt.Printf("📐 Алгоритм даунсамплинга: %s\n", cfg.Downsampling.Algorithm)

	fmt.Println("\n🔧 Реализация в процессе...")
	fmt.Println("✅ Это отмечает начало коммерческой разработки на Go: 31 декабря 2025")
}
