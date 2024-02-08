package handlers

import (
	"HabitsBot/internal/messages"
	"HabitsBot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sort"
)

type ByID []models.Habit

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

// Habits COMMAND = /my_habits
func (h *HabitBot) Habits(u *tgbotapi.Update) {
	userID := u.Message.From.ID

	habits, err := h.HabitStorage.GetAll(userID, 0)

	sort.Sort(ByID(habits))

	if err != nil {
		msgError := tgbotapi.NewMessage(userID, messages.ErrorGetHabits)
		h.Bot.Send(msgError)
	}

	msg := tgbotapi.NewMessage(userID, messages.HabitListMsg(habits))
	msg.ParseMode = "HTML"
	h.Bot.Send(msg)
}
