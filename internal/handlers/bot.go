package handlers

import (
	"HabitsBot/internal/models"
	"HabitsBot/internal/storages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
	"sync"
)

type HabitBot struct {
	Bot              *tgbotapi.BotAPI
	FSMMap           map[int64]*fsm.FSM
	mu               sync.Mutex
	HabitStorage     storages.HabitStorage
	TimeShedulerChan chan models.Habit
}

type Handler func(update *tgbotapi.Update, done *bool)
