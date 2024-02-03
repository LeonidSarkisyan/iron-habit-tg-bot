package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *HabitBot) AnswerCallbackQuery(update *tgbotapi.Update) {
	callbackConfig := tgbotapi.CallbackConfig{
		CallbackQueryID: update.CallbackQuery.ID,
	}
	_, _ = h.Bot.Send(callbackConfig)
}
