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

type TimeShedulerData struct {
	models.Habit
	models.Timestamp
}

type TimeShedulerChan = chan TimeShedulerData

type HabitBot struct {
	Bot              *tgbotapi.BotAPI
	FSMMap           map[int64]*fsm.FSM
	mu               sync.Mutex
	HabitStorage     storages.HabitStorage
	TimestampStorage storages.TimestampStorage
	RejectionStorage storages.RejectionStorage
	ExecutionStorage storages.ExecutionStorage
	TimeShedulerChan TimeShedulerChan
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
		TimestampStorage: postgres.NewTimestamp(db),
		ExecutionStorage: postgres.NewExecutionStorage(db),
		TimeShedulerChan: make(TimeShedulerChan),
		ControlChanMap:   make(map[int]*chan string),
		CompleteChanMap:  make(map[int]*chan string),
	}
}

type Handler func(update *tgbotapi.Update, done *bool)

type HandlerFunc func(update *tgbotapi.Update)
