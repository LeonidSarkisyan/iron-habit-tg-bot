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
	RejectionStorage storages.RejectionStorage
	TimeShedulerChan chan models.Habit
	ControlChanMap   map[int]*chan string
}

type Handler func(update *tgbotapi.Update, done *bool)
