package handlers

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *HabitBot) HandleFSMHabit(update *tgbotapi.Update, done *bool) {
	state := h.FSM(*update).Current()

	switch state {
	case "gettingHabitName":
		habitName := update.Message.Text
		msgText := fmt.Sprintf("Ваша привычка: %s", habitName)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		_, _ = h.Bot.Send(msg)
		_ = h.FSM(*update).Event(context.TODO(), "start")
	}
}
