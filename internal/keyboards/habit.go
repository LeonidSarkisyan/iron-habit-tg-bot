package keyboards

import (
	"HabitsBot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func CancelHabitKeyboard(habit models.Habit) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"Отменить выполнение привычки  ❌", "cancel_habit__"+strconv.Itoa(habit.ID),
		),
	))
	return keyboard
}
