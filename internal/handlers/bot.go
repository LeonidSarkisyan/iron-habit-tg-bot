package handlers

import (
	"HabitsBot/internal/models"
	"HabitsBot/internal/storages"
	"HabitsBot/internal/storages/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
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
	CompleteChanMap  map[int]*chan string
}

func New(bot *tgbotapi.BotAPI, db *pgx.Conn) *HabitBot {
	return &HabitBot{
		Bot:              bot,
		mu:               sync.Mutex{},
		FSMMap:           make(map[int64]*fsm.FSM),
		HabitStorage:     postgres.New(db),
		RejectionStorage: postgres.NewRejection(db),
		TimeShedulerChan: make(chan models.Habit),
		ControlChanMap:   make(map[int]*chan string),
		CompleteChanMap:  make(map[int]*chan string),
	}
}

type Handler func(update *tgbotapi.Update, done *bool)

type HandlerFunc func(update *tgbotapi.Update)
