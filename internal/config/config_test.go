package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"tg-reminder/internal/config"
)

func TestLoad_NoEnvFile(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("BOT_TOKEN", "test-token")
	if _, err := config.Load(); err != nil {
		t.Fatalf("missing .env must not be an error, got %v", err)
	}
}

func TestLoad_BroeknEnvFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("Bad Line\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)

	t.Setenv("BOT_TOKEN", "test-token")
	_, err := config.Load()
	if err == nil {
		t.Fatalf("broken .env file should produce error")
	}
}

func TestLoad_EmptyBoatToken(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("BOT_TOKEN", "")
	_, err := config.Load()
	if err == nil {
		t.Fatalf("Missing Bot token should produce error")
	}
}

func TestLoad_EnvVars(t *testing.T) {
	t.Chdir(t.TempDir())
	tests := []struct {
		name    string
		key     string
		value   string
		warning string
	}{
		{
			name:    "Invalid debug value",
			key:     "DEBUG",
			value:   "123",
			warning: fmt.Sprintf("неизвестный DEBUG=%q, используем false", "123"),
		},
		{
			name:    "Invalid log level",
			key:     "LOG_LEVEL",
			value:   "123",
			warning: fmt.Sprintf("неизвестный LOG_LEVEL=%q, используем INFO", "123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("BOT_TOKEN", "test-token")
			t.Setenv(tt.key, tt.value)
			cfg, err := config.Load()
			if err != nil {
				t.Fatalf("Should not produce error for optional params %v", err)
			}
			if !slices.Contains(cfg.Warnings, tt.warning) {
				t.Errorf("Not produces expected warnings, expected=%v, got=%v", tt.warning, cfg.Warnings)
			}
		})
	}
}

func TestDbString(t *testing.T) {
	t.Chdir(t.TempDir())

	t.Setenv("BOT_TOKEN", "test-token")
	t.Setenv("POSTGRES_HOST", "test_host")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_USER", "test_user")
	t.Setenv("POSTGRES_PASSWORD", "test_password")
	t.Setenv("POSTGRES_DB", "test_db")

	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}

	dbString := cfg.DB.DBString()
	testString := "postgres://test_user:test_password@test_host:5432/test_db"
	if dbString != testString {
		t.Fatalf("DBString method should return %s, actually got %s", testString, dbString)
	}
}
