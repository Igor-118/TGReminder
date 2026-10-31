package bot

import (
	"context"
	"log/slog"
	"tg-reminder/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.Message == nil { // If we got a message
		return
	}
	id := update.Message.From.ID
	text := update.Message.Text
	var reply service.Reply

	slog.Debug("Got an update message", "user", update.Message.From.UserName, "text", update.Message.Text)

	switch update.Message.Command() {
	case "list":
		reply = b.reminders.List(ctx, id)
	case "profile":
		reply = b.profile.Show(ctx, id)
	default:
		reply = b.registration.HandleInput(ctx, id, text)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	_, err := b.bot.Send(msg)
	if err != nil {
		slog.Error("Error while sending message", "error", err)
	}
}
