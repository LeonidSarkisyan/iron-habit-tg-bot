package main

import (
	"HabitsBot/internal/handlers"
	"HabitsBot/internal/models"
	"HabitsBot/internal/shedulers"
	"HabitsBot/internal/storages/postgres"
	"HabitsBot/pkg/systems"
	"context"
	"github.com/looplab/fsm"
	"github.com/rs/zerolog/log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

func main() {
	systems.SetupLogger()
	systems.MustLoadEnv()

	db, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))

	if err != nil {
		log.Fatal().Err(err).Msg("Не удалось подключиться к базе данных")
	}

	token := os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal().Err(err).Msg("Не удалось создать бота")
	}

	handlers.MustCommands(bot)

	habitBot := &handlers.HabitBot{
		Bot:              bot,
		FSMMap:           make(map[int64]*fsm.FSM),
		HabitStorage:     postgres.New(db),
		RejectionStorage: postgres.NewRejection(db),
		TimeShedulerChan: make(chan models.Habit),
		ControlChanMap:   make(map[int]*chan string),
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	callbackRouters := []handlers.Handler{
		habitBot.HandleFSMHabit,
	}

	messageRouters := []handlers.Handler{
		habitBot.HandleCommand,
		habitBot.HandleFSMHabit,
	}

	go func() {
		for habit := range habitBot.TimeShedulerChan {
			go shedulers.AddHabitToTiming(habit, habitBot)
		}
	}()

	for update := range updates {
		done := false

		log.Info().Msgf("%+v", update)

		if update.CallbackQuery != nil {
			passHandlers(&update, &done, callbackRouters...)
		}

		if update.Message != nil {
			passHandlers(&update, &done, messageRouters...)
		}
	}
}

func passHandlers(update *tgbotapi.Update, done *bool, handlers ...handlers.Handler) {
	for _, handler := range handlers {
		handler(update, done)
		if *done {
			break
		}
	}
}
