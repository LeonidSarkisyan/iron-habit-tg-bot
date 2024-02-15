package filters

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func IsDetail(update *tgbotapi.Update) bool {
	return strings.HasPrefix(update.CallbackQuery.Data, "habit__")
}
