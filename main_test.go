package main

import (
	"os"
	"testing"
)

// TestVersionConst проверяет, что версия задана.
func TestVersionConst(t *testing.T) {
	if Version == "" {
		t.Fatal("Version is empty")
	}
}

// TestDefaultConfig проверяет, что DefaultConfig возвращает валидные значения.
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("DefaultConfig validation failed: %v", err)
	}
	if cfg.Rasterization.BaseDPI != 300 {
		t.Fatalf("expected BaseDPI=300, got %d", cfg.Rasterization.BaseDPI)
	}
	if len(cfg.Rasterization.PyramidLevels) == 0 {
		t.Fatal("PyramidLevels is empty")
	}
	if cfg.Output.Format != "webp" {
		t.Fatalf("expected Format=webp, got %s", cfg.Output.Format)
	}
}

// TestConfigLoad проверяет загрузку существующего YAML-файла.
func TestConfigLoad(t *testing.T) {
	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadConfig returned nil")
	}
}

// TestConfigLoadMissing проверяет, что загрузка несуществующего файла возвращает ошибку.
func TestConfigLoadMissing(t *testing.T) {
	_, err := LoadConfig("nonexistent.yaml")
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

// TestConfigValidate проверяет валидацию невалидного конфига.
func TestConfigValidate(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Rasterization.BaseDPI = 1000
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for BaseDPI=1000")
	}
	cfg.Rasterization.BaseDPI = 300
	cfg.Rasterization.PyramidLevels = []int{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty PyramidLevels")
	}
	cfg.Rasterization.PyramidLevels = []int{300}
	cfg.Output.WebPQuality = 200
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for WebPQuality=200")
	}
}

// TestConfigValidateWebPQualityEdge проверяет граничные значения качества.
func TestConfigValidateWebPQualityEdge(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Output.WebPQuality = 1
	if err := cfg.Validate(); err != nil {
		t.Fatalf("WebPQuality=1 should be valid: %v", err)
	}
	cfg.Output.WebPQuality = 100
	if err := cfg.Validate(); err != nil {
		t.Fatalf("WebPQuality=100 should be valid: %v", err)
	}
}

// TestEnvFileExample проверяет, что .env.example существует и непустой.
func TestEnvFileExample(t *testing.T) {
	data, err := os.ReadFile(".env.example")
	if err != nil {
		t.Fatal(".env.example not found")
	}
	if len(data) == 0 {
		t.Fatal(".env.example is empty")
	}
}

// TestConfigYamlExists проверяет, что config.yaml существует.
func TestConfigYamlExists(t *testing.T) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		t.Fatal("config.yaml not found")
	}
	if len(data) == 0 {
		t.Fatal("config.yaml is empty")
	}
}

// TestLicenseExists проверяет, что LICENSE существует.
func TestLicenseExists(t *testing.T) {
	data, err := os.ReadFile("LICENSE")
	if err != nil {
		t.Fatal("LICENSE not found")
	}
	if len(data) == 0 {
		t.Fatal("LICENSE is empty")
	}
}

// TestGitIgnoreExists проверяет, что .gitignore существует.
func TestGitIgnoreExists(t *testing.T) {
	data, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatal(".gitignore not found")
	}
	if len(data) == 0 {
		t.Fatal(".gitignore is empty")
	}
}
