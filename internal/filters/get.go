package filters

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func IsPage(update *tgbotapi.Update) bool {
	return strings.HasPrefix(update.CallbackQuery.Data, "habits__page__")
}
