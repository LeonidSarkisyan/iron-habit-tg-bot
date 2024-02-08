package filters

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func IsCallbackQueryEmpty(update *tgbotapi.Update) bool {
	return update.CallbackQuery == nil
}

func IsCallbackQuery(update *tgbotapi.Update) bool {
	return update.CallbackQuery != nil
}

func IsCallBackDataAddReminder(update *tgbotapi.Update) bool {
	return strings.HasPrefix(update.CallbackQuery.Data, "add_reminder__")
}
