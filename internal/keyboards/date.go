package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	DoneEmoji = "✅"
)

func DaysPickerKeyboard(days []string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	daysOfWeek := []CallBackData{
		{Name: "Понедельник", Data: "1"},
		{Name: "Вторник", Data: "2"},
		{Name: "Среда", Data: "3"},
		{Name: "Четверг", Data: "4"},
		{Name: "Пятница", Data: "5"},
		{Name: "Суббота", Data: "6"},
		{Name: "Воскресенье", Data: "7"},
	}

	for _, day := range days {
		for i, dayOfWeek := range daysOfWeek {
			if day == dayOfWeek.Data {
				daysOfWeek[i].Name += "  " + DoneEmoji
				break
			}
		}
	}

	daysOfWeek = append(daysOfWeek, CallBackData{Name: "Продолжить  ⏩", Data: "continue"})

	for _, day := range daysOfWeek {
		rows = append(rows, createInlineKeyboardRow(day.Name, day.Data))
	}

	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func createInlineKeyboardRow(text, callbackData string) []tgbotapi.InlineKeyboardButton {
	btn := tgbotapi.NewInlineKeyboardButtonData(text, callbackData)
	return []tgbotapi.InlineKeyboardButton{btn}
}

func EditInlineKeyboard(newInlineKeyboard tgbotapi.InlineKeyboardMarkup, update *tgbotapi.Update) tgbotapi.EditMessageTextConfig {
	editMsg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		update.CallbackQuery.Message.Text,
	)

	editMsg.ReplyMarkup = &newInlineKeyboard

	return editMsg
}