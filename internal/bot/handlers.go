package bot

import (
	"context"
	"log/slog"
	"tg-reminder/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) sendReply(chatID int64, reply service.Reply) {
	msg := tgbotapi.NewMessage(chatID, reply.Text)

	if len(reply.Options) > 0 {
		msg.ReplyMarkup = optionsKeyboard(reply.Options)
	}

	_, err := b.bot.Send(msg)
	if err != nil {
		slog.Error("Error while sending message", "error", err)
	}
}

func (b *Bot) handleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	reply := b.registration.HandleInput(ctx, callback.From.ID, callback.Data)
	b.sendReply(callback.Message.Chat.ID, reply)

	_, err := b.bot.Request(tgbotapi.NewCallback(callback.ID, ""))
	if err != nil {
		slog.Error("Error while answering callback", "error", err)
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		b.handleCallback(ctx, update.CallbackQuery)
		return
	}
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
	case "remind":
		commandText := update.Message.CommandArguments()
		reply = b.reminders.Add(ctx, id, commandText)
	case "start":
		reply = b.registration.Start(ctx, id)
	default:
		reply = b.registration.HandleInput(ctx, id, text)
	}

	b.sendReply(update.Message.Chat.ID, reply)
}
