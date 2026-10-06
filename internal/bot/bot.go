package bot

import (
	"context"
	"fmt"
	"log/slog"
	"tg-reminder/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	registration *service.RegistrationService
}

func New(token string, debug bool, registration *service.RegistrationService) (*Bot, error) {
	bot1, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("подключение к telegram: %w", err)
	}
	if debug {
		bot1.Debug = true
	}

	slog.Info("Authorized on account", "username", bot1.Self.UserName)

	return &Bot{bot: bot1, registration: registration}, nil
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			b.handleUpdate(ctx, update)

		case <-ctx.Done():
			b.bot.StopReceivingUpdates()
			return nil
		}

	}
}
