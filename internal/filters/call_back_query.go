package filters

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func IsCallbackQueryEmpty(update *tgbotapi.Update) bool {
	return update.CallbackQuery == nil
}

func IsCallbackQuery(update *tgbotapi.Update) bool {
	return update.CallbackQuery != nil
}
