package keyboards

import (
	"HabitsBot/internal/models"
	"HabitsBot/internal/pagination"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func HabitsListKeyboard(habits []models.Habit, hp pagination.HabitPagination) tgbotapi.InlineKeyboardMarkup {
	var habitButtons [][]tgbotapi.InlineKeyboardButton

	pageStr := strconv.Itoa(hp.Page)

	for i, habit := range habits {
		i := strconv.Itoa(i + 1 + (hp.Page-1)*5)
		habitButtons = append(habitButtons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				i+". "+habit.Title, "habit__"+strconv.Itoa(habit.ID)+"__"+pageStr),
		))
	}

	if hp.Page == 1 && hp.ExistsMore {
		habitButtons = append(habitButtons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➡️", "habits__page__"+strconv.Itoa(hp.Page+1)),
		))
	} else if hp.ExistsMore {
		habitButtons = append(habitButtons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️", "habits__page__"+strconv.Itoa(hp.Page-1)),
			tgbotapi.NewInlineKeyboardButtonData("➡️", "habits__page__"+strconv.Itoa(hp.Page+1)),
		))
	} else if hp.Page != 1 {
		habitButtons = append(habitButtons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️", "habits__page__"+strconv.Itoa(hp.Page-1)),
		))
	}

	return tgbotapi.NewInlineKeyboardMarkup(habitButtons...)
}
