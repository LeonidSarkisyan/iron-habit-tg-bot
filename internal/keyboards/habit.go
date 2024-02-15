package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func BeforeCreateHabitReplyKeyboard() tgbotapi.ReplyKeyboardMarkup {
	k := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Создать привычку  ✨"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отмена  ❌"),
		),
	)

	k.OneTimeKeyboard = true

	return k
}

func CreateDayTimeInlineKeyboard(habitID int) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"Добавить напоминание  🔔", "add_reminder__"+strconv.Itoa(habitID),
		),
	))
	return keyboard
}
