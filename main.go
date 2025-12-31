package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Version information
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

	flag.StringVar(&inputPath, "input", "", "Input folder path with documents")
	flag.StringVar(&outputPath, "output", "", "Output folder path for processed documents")
	flag.IntVar(&dpi, "dpi", 300, "Base DPI for rasterization (default: 300)")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.Parse()

	if showVersion {
		fmt.Printf("raster-forge v%s (built: %s)\n", Version, BuildDate)
		fmt.Println("High-performance document normalizer and rasterization engine")
		os.Exit(0)
	}

	if inputPath == "" {
		log.Fatal("Error: input path is required. Use -input flag.")
	}

	fmt.Printf("🔨 raster-forge v%s\n", Version)
	fmt.Printf("📂 Input:  %s\n", inputPath)
	fmt.Printf("📁 Output: %s\n", getOutputPath(outputPath))
	fmt.Printf("🎨 DPI:    %d\n\n", dpi)

	fmt.Println("⚠️  Implementation in progress...")
	fmt.Println("\n✅ This marks the beginning of commercial Go development: December 31, 2025")
}

func getOutputPath(path string) string {
	if path != "" {
		return path
	}
	return "output_<timestamp>"
}
