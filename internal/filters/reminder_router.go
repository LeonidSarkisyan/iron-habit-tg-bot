package filters

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func IsReminderListCallBackData(update *tgbotapi.Update) bool {
	return strings.HasPrefix(update.CallbackQuery.Data, "reminder_list__")
}

func IsReminderPage(update *tgbotapi.Update) bool {
	return strings.HasPrefix(update.CallbackQuery.Data, "reminder_page")
}
