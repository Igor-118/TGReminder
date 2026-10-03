package bot

import (
	"context"
	"fmt"
	"log"
	// "os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//"fmt"
)

type Bot struct {
	bot tgbotapi.BotAPI
}

func New(token string) (*Bot, error) {
	bot1, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("Ошибка в получении токена %w", err)
	}

	// bot1.Debug = true

	log.Printf("Authorized on account %s", bot1.Self.UserName)

	return &Bot{bot: *bot1}, nil
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
