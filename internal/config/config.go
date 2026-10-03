package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"tg-reminder/pkg/logger"

	"github.com/joho/godotenv"
)

type Config struct {
	DB DB
	Telegram
	Debug    bool
	Logger   logger.Config
	Warnings []string
}

type DB struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgreHost      string
	PostgresPort     string
}

func (db *DB) DBString() string {
	//DATABASE_URL :="postgres://POSTGRES_USER:POSTGRES_PASSWORD@POSTGRES_HOST:POSTGRES_PORT/POSTGRES_DB"
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db.PostgresUser, db.PostgresPassword, db.PostgreHost, db.PostgresPort, db.PostgresDB)
}

type Telegram struct {
	BotToken string
}

func Load() (*Config, error) {
	cfg := &Config{}

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("пустой бот токен")
	}
	var debug bool
	if v := os.Getenv("DEBUG"); v != "" {
		d, err := strconv.ParseBool(v)
		if err != nil {
			cfg.Warnings = append(cfg.Warnings, fmt.Sprintf("неизвестный DEBUG=%q, используем false", v))
		}
		debug = d
	}

	var lvl slog.Level
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		if err := lvl.UnmarshalText([]byte(v)); err != nil {
			cfg.Warnings = append(cfg.Warnings, fmt.Sprintf("неизвестный LOG_LEVEL=%q, используем INFO", v))
		}
	}

	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")

	cfg.Telegram = Telegram{
		BotToken: botToken,
	}
	cfg.DB = DB{
		PostgresUser:     postgresUser,
		PostgresPassword: postgresPassword,
		PostgresDB:       postgresDB,
		PostgreHost:      postgresHost,
		PostgresPort:     postgresPort,
	}
	cfg.Debug = debug
	cfg.Logger = logger.Config{
		LogLevel: lvl,
	}

	return cfg, nil
}
