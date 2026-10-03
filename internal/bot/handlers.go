package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	//"updates"
)

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {

	if update.Message == nil { // If we got a message
		return
	}
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	_, err := b.bot.Send(msg)
	if err != nil {
		log.Printf("send message: %v", err)
	}
}
