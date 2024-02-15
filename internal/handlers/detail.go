package handlers

import (
	"HabitsBot/internal/filters"
	"HabitsBot/internal/keyboards"
	"HabitsBot/internal/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

func NewHabitDetailRouter(habitBot *HabitBot) *Router {
	r := NewRouter(habitBot)

	r.CallBackQuery(habitBot.ShowMenuHabit, filters.IsDetail)

	return r
}

func (h *HabitBot) ShowMenuHabit(update *tgbotapi.Update) {
	data := strings.Split(strings.Replace(update.CallbackQuery.Data, "habit__", "", 1), "__")

	habitID := data[0]
	page := data[1]

	habitIDInt, err := strconv.Atoi(habitID)

	if err != nil {
		log.Error().Err(err)
		return
	}

	pageInt, err := strconv.Atoi(page)

	if err != nil {
		log.Error().Err(err)
		return
	}

	userID := update.CallbackQuery.From.ID

	habit, err := h.HabitStorage.Get(userID, habitIDInt)

	if err != nil {
		msgError := tgbotapi.NewMessage(userID, messages.ErrorGetHabits)
		_, _ = h.Bot.Send(msgError)
		return
	}

	update.CallbackQuery.Message.Text = messages.HabitDetailMsg(habit)
	msg := keyboards.EditInlineKeyboard(keyboards.MenuHabitKeyboard(habit, pageInt), update)
	msg.ParseMode = tgbotapi.ModeHTML
	_, _ = h.Bot.Send(msg)
}

func (h *HabitBot) AskNewHabitTitle(update *tgbotapi.Update) {

}
