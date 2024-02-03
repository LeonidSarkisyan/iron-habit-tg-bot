package handlers

import (
	"HabitsBot/internal/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

const (
	startCommand       = "start"
	addNewHabitCommand = "add_new_habit"
	checkFSMState      = "check_fsm_state"
)

func (h *HabitBot) HandleCommand(update *tgbotapi.Update, done *bool) {
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case startCommand:
			h.handleStartCommand(update, done)
		case addNewHabitCommand:
			h.handleAddNewHabitCommand(update, done)
		case checkFSMState:
			h.handleFSMCheckCommand(update, done)
		}
	}
}

func (h *HabitBot) handleStartCommand(update *tgbotapi.Update, done *bool) {
	h.FSM(update).SetState(startState)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.StartMsg)
	_, err := h.Bot.Send(msg)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	*done = true
}

func (h *HabitBot) handleAddNewHabitCommand(update *tgbotapi.Update, done *bool) {
	h.FSM(update).SetState(getHabitNameState)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.InputHabitNameMsg)
	h.Bot.Send(msg)
	*done = true
}

func (h *HabitBot) handleFSMCheckCommand(update *tgbotapi.Update, done *bool) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "вы в состоянии: "+h.FSM(update).Current())
	h.Bot.Send(msg)
	*done = true
}
