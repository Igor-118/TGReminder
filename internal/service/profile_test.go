package service_test

import (
	"fmt"
	//"os"
	//"path/filepath"
	"slices"
	"testing"

	"tg-reminder/internal/config"
)

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
