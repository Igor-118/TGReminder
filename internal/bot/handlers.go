package bot

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {

	if update.Message == nil { // If we got a message
		return
	}
	slog.Debug("Got an update message", "user", update.Message.From.UserName, "text", update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	_, err := b.bot.Send(msg)
	if err != nil {
		slog.Error("Error while sending message", "error", err)
	}
}
