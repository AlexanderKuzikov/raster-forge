package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestVersionFlag проверяет, что `-version` выводит версию и завершается с кодом 0.
func TestVersionFlag(t *testing.T) {
	// Собираем бинарник
	bin := filepath.Join(t.TempDir(), "raster-forge")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	// Запускаем с -version
	cmd = exec.Command(bin, "-version")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got: %v\n%s", err, out)
	}

	if len(out) == 0 {
		t.Fatal("expected version output, got empty")
	}
}

// TestMissingInput проверяет, что без -input программа завершается с ошибкой.
func TestMissingInput(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "raster-forge")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	cmd = exec.Command(bin)
	err = cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing -input")
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
	// Используем config.yaml из корня проекта
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

// TestEnvFileExample проверяет, что .env.example существует (непустой).
func TestEnvFileExample(t *testing.T) {
	data, err := os.ReadFile(".env.example")
	if err != nil {
		t.Fatal(".env.example not found")
	}
	if len(data) == 0 {
		t.Fatal(".env.example is empty")
	}
}
