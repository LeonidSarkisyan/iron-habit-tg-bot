package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
	"sync"
)

type HabitBot struct {
	Bot          *tgbotapi.BotAPI
	FSMMap       map[int64]*fsm.FSM
	mu           sync.Mutex
	HabitStorage map[int64]string
}

type Handler func(update *tgbotapi.Update, done *bool)
