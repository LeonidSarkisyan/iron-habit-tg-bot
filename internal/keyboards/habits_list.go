package keyboards

import (
	"HabitsBot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func HabitsListKeyboard(habits []models.Habit, page int, existsMore bool) tgbotapi.InlineKeyboardMarkup {
	var habitButtons [][]tgbotapi.InlineKeyboardButton

	pageStr := strconv.Itoa(page)

	for i, habit := range habits {
		i := strconv.Itoa(i + 1 + (page-1)*5)
		habitButtons = append(habitButtons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				i+". "+habit.Title, "habit__"+strconv.Itoa(habit.ID)+"__"+pageStr),
		))
	}

	if page == 1 && existsMore {
		habitButtons = append(habitButtons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➡️", "habits__page__"+strconv.Itoa(page+1)),
		))
	} else if existsMore {
		habitButtons = append(habitButtons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️", "habits__page__"+strconv.Itoa(page-1)),
			tgbotapi.NewInlineKeyboardButtonData("➡️", "habits__page__"+strconv.Itoa(page+1)),
		))
	} else if page != 1 {
		habitButtons = append(habitButtons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️", "habits__page__"+strconv.Itoa(page-1)),
		))
	}

	return tgbotapi.NewInlineKeyboardMarkup(habitButtons...)
}
