package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Информация о версии
const (
	Version = "0.1.0"
	BuildDate = "2025-12-31"
)

func main() {
	var (
		inputPath  string
		outputPath string
		dpi        int
		showVersion bool
	)

	flag.StringVar(&inputPath, "input", "", "Путь к входной папке с документами")
	flag.StringVar(&outputPath, "output", "", "Путь к выходной папке для обработанных документов")
	flag.IntVar(&dpi, "dpi", 300, "Базовое разрешение для растеризации (по умолчанию: 300)")
	flag.BoolVar(&showVersion, "version", false, "Показать информацию о версии")
	flag.Parse()

	if showVersion {
		fmt.Printf("raster-forge v%s (собрано: %s)\n", Version, BuildDate)
		fmt.Println("Высокопроизводительный движок нормализации и растеризации документов")
		os.Exit(0)
	}

	if inputPath == "" {
		log.Fatal("Ошибка: требуется указать входной путь. Используйте флаг -input.")
	}

	fmt.Printf("🔨 raster-forge v%s\n", Version)
	fmt.Printf("📂 Вход:   %s\n", inputPath)
	fmt.Printf("📁 Выход:  %s\n", getOutputPath(outputPath))
	fmt.Printf("🎨 DPI:     %d\n", dpi)
	fmt.Println("🔽 Пирамида: 75, 100, 150, 200, 250, 300 DPI\n")

	fmt.Println("⚠️  Реализация в процессе...")
	fmt.Println("\n✅ Это отмечает начало коммерческой разработки на Go: 31 декабря 2025")
}

func getOutputPath(path string) string {
	if path != "" {
		return path
	}
	return "output_<timestamp>"
}
