package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-reminder/internal/service"
)

func optionsKeyboard(opts []service.Option) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for _, opt := range opts {
		button := tgbotapi.NewInlineKeyboardButtonData(opt.Label, opt.Value)
		row = append(row, button)

		if len(row) == 3 {
			rows = append(rows, row)
			row = nil
		}
	}

	if len(row) > 0 {
		rows = append(rows, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
