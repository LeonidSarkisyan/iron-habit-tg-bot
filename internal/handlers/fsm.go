package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
)

const (
	startState        = "start"
	getHabitNameState = "gettingHabitName"
)

func (h *HabitBot) FSM(update tgbotapi.Update) *fsm.FSM {
	userID := update.Message.From.ID
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.FSMMap[userID]; !ok {
		h.FSMMap[userID] = fsm.NewFSM(
			"start",
			fsmEvents(),
			fsm.Callbacks{},
		)
	}

	return h.FSMMap[userID]
}

func fsmEvents() fsm.Events {
	return fsm.Events{
		{Name: "start", Src: []string{getHabitNameState}, Dst: startState},
		{Name: "addNewHabit", Src: []string{startState}, Dst: getHabitNameState},
	}
}
