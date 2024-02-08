package filters

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func IsCallBackDataContinue(update *tgbotapi.Update) bool {
	return update.CallbackQuery.Data == "continue"
}

func IsCallBackDataTimeContinue(update *tgbotapi.Update) bool {
	return update.CallbackQuery.Data == "time__continue"
}

func IsCallBackDataStartWithTime(update *tgbotapi.Update) bool {
	return strings.HasPrefix(update.CallbackQuery.Data, "time__")
}

func IsCallBackDataStartWithDay(update *tgbotapi.Update) bool {
	return strings.HasPrefix(update.CallbackQuery.Data, "day__")
}

func IsCallBackDataCancelHabit(update *tgbotapi.Update) bool {
	return strings.HasPrefix(update.CallbackQuery.Data, "cancel_habit__")
}

func IsCreateHabitText(update *tgbotapi.Update) bool {
	return update.Message.Text == "Создать привычку  ✨"
}

func IsCancelCreateHabitText(update *tgbotapi.Update) bool {
	return update.Message.Text == "Отмена  ❌"
}
