package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"tg-reminder/internal/bot"
	"tg-reminder/internal/config"
)

type App struct {
	pool *pgxpool.Pool
	bot  *bot.Bot
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("загрузка конфига не прошла: %w", err)
	}

	pool, err := pgxpool.New(ctx, cfg.DB.DBString())
	if err != nil {
		return nil, fmt.Errorf("ошибка с пулом: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("соединиться к бд не прошло: %w", err)
	}
	log.Println("Соединение прошло успешно")

	jopaBota, err := bot.New(cfg.BotToken)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("не получилось создать бота: %w", err)
	}

	return &App{
		pool: pool,
		bot:  jopaBota,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	err := a.bot.Run(ctx)
	if err != nil {
		return fmt.Errorf("бот не запустился %w", err)
	}
	return err
}

func (a *App) Close() {
	a.pool.Close()
}
