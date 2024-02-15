package keyboards

import (
	"HabitsBot/internal/callbackdata"
	"HabitsBot/internal/models"
	"HabitsBot/internal/pagination"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func ReminderListKeyboard(
	habitID int, reminders []models.Timestamp, rh pagination.ReminderPagination, pageHabit string,
) tgbotapi.InlineKeyboardMarkup {
	var reminderButtons [][]tgbotapi.InlineKeyboardButton

	habitIDStr := strconv.Itoa(habitID)

	for _, reminder := range reminders {
		reminderButtons = append(reminderButtons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(reminder.Time, "reminder__"+strconv.Itoa(reminder.ID)),
		))
	}

	cbd := callbackdata.NewCallBackData("reminder_page")
	cbd.Add(habitIDStr).Add(pageHabit)

	if rh.Page == 1 && rh.ExistsMore {
		reminderButtons = append(reminderButtons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➡️", cbd.Add(rh.NextPage()).String()),
		))
	} else if rh.ExistsMore {
		reminderButtons = append(reminderButtons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️", cbd.Add(rh.PrevPage()).String()),
			tgbotapi.NewInlineKeyboardButtonData("➡️", cbd.Add(rh.NextPage()).String()),
		))
	} else if rh.Page != 1 {
		reminderButtons = append(reminderButtons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️", cbd.Add(rh.PrevPage()).String()),
		))
	}

	reminderButtons = append(reminderButtons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"Назад",
			"habit__"+habitIDStr+"__"+pageHabit,
		),
	))

	return tgbotapi.NewInlineKeyboardMarkup(reminderButtons...)
}
