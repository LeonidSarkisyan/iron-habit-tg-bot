package keyboards

import (
	"HabitsBot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func MenuHabitKeyboard(habit models.Habit, page int) tgbotapi.InlineKeyboardMarkup {
	var habitButtons [][]tgbotapi.InlineKeyboardButton

	habitID := strconv.Itoa(habit.ID)
	pageStr := strconv.Itoa(page)

	habitButtons = append(habitButtons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Список напоминаний  🔔", "reminder_list__"+habitID+"__"+pageStr),
	))

	habitButtons = append(habitButtons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Изменить название  ✏️", "change_title__"+habitID),
	))

	habitButtons = append(habitButtons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Удалить привычку  ❌", "change_days__"+habitID),
	))

	habitButtons = append(habitButtons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад  ⬅️", "habits__page__"+pageStr),
	))

	return tgbotapi.NewInlineKeyboardMarkup(habitButtons...)
}
