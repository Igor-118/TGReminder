package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"tg-reminder/internal/bot"
	"tg-reminder/internal/config"
	"tg-reminder/internal/repository"
	"tg-reminder/internal/service"
	"tg-reminder/pkg/logger"
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

	lg := logger.NewLogger(&cfg.Logger)
	slog.SetDefault(lg)
	for _, w := range cfg.Warnings {
		slog.Warn(w)
	}

	pool, err := pgxpool.New(ctx, cfg.DB.DBString())
	if err != nil {
		return nil, fmt.Errorf("ошибка с пулом: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("соединиться к бд не прошло: %w", err)
	}
	slog.Info("Соединение прошло успешно")
	userRepo := repository.NewUserRepository(pool)
	registration := service.NewRegistrationService(userRepo)
	jopaBota, err := bot.New(cfg.BotToken, cfg.Debug, registration)
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
