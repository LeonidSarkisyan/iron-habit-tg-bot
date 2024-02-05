package handlers

import (
	"HabitsBot/internal/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func (h *HabitBot) SendNotification(habit models.Habit) {
	msg := tgbotapi.NewMessage(habit.UserID, fmt.Sprintf("Напоминаю о вашей привычке: <b>%s!</b>", habit.Title))
	msg.ParseMode = "html"
	_, err := h.Bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	}
}

func (h *HabitBot) SendWarningBeforeNotification(habit models.Habit) {
	msg := tgbotapi.NewMessage(habit.UserID, fmt.Sprintf("Через 15 минут я напомню о вашей привычке: <b>%s!</b>", habit.Title))
	msg.ParseMode = "html"
	_, err := h.Bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	}
}
