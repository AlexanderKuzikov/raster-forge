package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config - основная структура конфигурации
type Config struct {
	Version        string              `yaml:"version"`
	Rasterization  RasterizationConfig `yaml:"rasterization"`
	Output         OutputConfig        `yaml:"output"`
	Processing     ProcessingConfig    `yaml:"processing"`
	Downsampling   DownsamplingConfig  `yaml:"downsampling"`
	Logging        LoggingConfig       `yaml:"logging"`
}

// RasterizationConfig - настройки растеризации
type RasterizationConfig struct {
	BaseDPI        int   `yaml:"base_dpi"`
	PyramidLevels  []int `yaml:"pyramid_levels"`
}

// OutputConfig - настройки выходных файлов
type OutputConfig struct {
	Format            string `yaml:"format"`
	WebPQuality       int    `yaml:"webp_quality"`
	PDFCompression    bool   `yaml:"pdf_compression"`
	PreserveMetadata  bool   `yaml:"preserve_metadata"`
}

// ProcessingConfig - настройки обработки
type ProcessingConfig struct {
	ParallelPages  bool `yaml:"parallel_pages"`
	MaxGoroutines  int  `yaml:"max_goroutines"`
	MemoryLimitMB  int  `yaml:"memory_limit_mb"`
}

// DownsamplingConfig - настройки downsampling
type DownsamplingConfig struct {
	Algorithm string `yaml:"algorithm"`
}

// LoggingConfig - настройки логирования
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

// LoadConfig - загрузка конфигурации из YAML файла
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфиг: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("не удалось распарсить YAML: %w", err)
	}

	// Валидация конфига
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("некорректный конфиг: %w", err)
	}

	return &cfg, nil
}

// Validate - проверка корректности конфигурации
func (c *Config) Validate() error {
	if c.Rasterization.BaseDPI < 50 || c.Rasterization.BaseDPI > 600 {
		return fmt.Errorf("base_dpi должен быть между 50 и 600, получено: %d", c.Rasterization.BaseDPI)
	}

	if len(c.Rasterization.PyramidLevels) == 0 {
		return fmt.Errorf("pyramid_levels не может быть пустым")
	}

	if c.Output.WebPQuality < 1 || c.Output.WebPQuality > 100 {
		return fmt.Errorf("webp_quality должен быть между 1 и 100, получено: %d", c.Output.WebPQuality)
	}

	return nil
}

// DefaultConfig - возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		Version: "1.0",
		Rasterization: RasterizationConfig{
			BaseDPI:       300,
			PyramidLevels: []int{300, 250, 200, 150, 100, 75},
		},
		Output: OutputConfig{
			Format:           "webp",
			WebPQuality:      90,
			PDFCompression:   true,
			PreserveMetadata: false,
		},
		Processing: ProcessingConfig{
			ParallelPages: true,
			MaxGoroutines: 0,
			MemoryLimitMB: 2048,
		},
		Downsampling: DownsamplingConfig{
			Algorithm: "lanczos3",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
			Output: "stdout",
		},
	}
}
