package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
)

const (
	startState            = "start"
	getHabitNameState     = "gettingHabitName"
	getHabitDaysState     = "gettingHabitDays"
	getWarningTimeState   = "gettingWarningTime"
	getCompletedTimeState = "gettingCompletedTime"
)

func (h *HabitBot) Clear(update *tgbotapi.Update, keys ...string) {
	for _, key := range keys {
		h.FSM(update).DeleteMetadata(key)
	}
}

func (h *HabitBot) FSM(update *tgbotapi.Update) *fsm.FSM {
	var userID int64

	if update.Message != nil {
		userID = update.Message.From.ID
	} else if update.CallbackQuery != nil {
		userID = update.CallbackQuery.From.ID
	} else {
		panic("Невозможно идентифицитировать юзера")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if existingFSM, ok := h.FSMMap[userID]; ok {
		return existingFSM
	}

	newFSM := fsm.NewFSM(
		"initial",
		fsmEvents(),
		fsm.Callbacks{},
	)

	h.FSMMap[userID] = newFSM
	return newFSM
}

func fsmEvents() fsm.Events {
	return []fsm.EventDesc{}
}
