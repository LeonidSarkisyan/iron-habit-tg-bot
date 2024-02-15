package main

import (
	"HabitsBot/internal/handlers"
	"HabitsBot/internal/shedulers"
	"HabitsBot/pkg/systems"
	"context"
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

	habitBot := handlers.New(bot, db)
	habits, err := habitBot.HabitStorage.Habits()

	if err != nil {
		log.Fatal().Err(err).Msg("невозможно запустить таймеры привычек")
	}

	go shedulers.AddManyHabitsToTiming(habits, habitBot)
	go shedulers.HabitListener(habitBot)

	log.Info().Msg("Таймеры привычек были запущены")
	log.Printf("Authorized on account %s", bot.Self.UserName)

	d := handlers.NewDispatcher(habitBot)

	commandRouter := handlers.NewCommandRouter(habitBot)

	habitsFSMRouter := handlers.NewHabitsFSMRouter(habitBot)
	habitsCallBackRouter := handlers.NewHabitsCallBackRouter(habitBot)
	statsRouter := handlers.NewStatsRouter(habitBot)

	habitList := handlers.NewHabitListRouter(habitBot)
	habitMenu := handlers.NewHabitDetailRouter(habitBot)

	reminderList := handlers.NewReminderRouter(habitBot)

	d.IncludeRouter(commandRouter)
	d.IncludeRouter(habitsFSMRouter)
	d.IncludeRouter(habitsCallBackRouter)
	d.IncludeRouter(statsRouter)
	d.IncludeRouter(habitList)
	d.IncludeRouter(habitMenu)
	d.IncludeRouter(reminderList)

	d.Polling()
}
