package handlers

import (
	"HabitsBot/internal/filters"
	"HabitsBot/internal/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

const (
	checkFSMState = "check_fsm_state"
	cancelCommand = "cancel"
)

func NewCommandRouter(habitBot *HabitBot) *Router {
	r := NewRouter(habitBot)

	r.Message(filters.F(filters.IsCommandStart), habitBot.HandleStartCommand)
	r.Message(filters.F(filters.IsCommandCancel), habitBot.HandleCancelCommand)
	r.Message(filters.F(filters.IsCommandAddNewHabit), habitBot.HandleAddNewHabitCommand)
	r.Message(filters.F(filters.IsCommandMyHabits), habitBot.Habits)

	return r
}

func (h *HabitBot) HandleStartCommand(update *tgbotapi.Update) {
	h.FSM(update).SetState(StartState)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.StartMsg)
	_, err := h.Bot.Send(msg)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func (h *HabitBot) HandleAddNewHabitCommand(update *tgbotapi.Update) {
	h.FSM(update).SetState(GetHabitNameState)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.InputHabitNameMsg)
	h.Bot.Send(msg)
}

func (h *HabitBot) HandleCancelCommand(update *tgbotapi.Update) {
	h.FSM(update).SetState(StartState)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.CancelMsg)
	h.Bot.Send(msg)
	h.Clear(update, "habit_id", "messages_ids")
}
