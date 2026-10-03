package config

import (
	//"context"
	//"log"
	"fmt"
	"os"

	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	//"github.com/jackc/pgx/v5"
	// "tgtest/repo"

	//"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Config struct {
	DB DB
	Telegram
	Debug bool
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
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Лоад не сработал, ошибочка %w", err)
	}
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("пустой бот токен")
	}

	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")

	cfg := Config{
		Telegram: Telegram{
			BotToken: botToken,
		},
		DB: DB{
			PostgresUser:     postgresUser,
			PostgresPassword: postgresPassword,
			PostgresDB:       postgresDB,
			PostgreHost:      postgresHost,
			PostgresPort:     postgresPort,
		},
	}

	return &cfg, nil
}
