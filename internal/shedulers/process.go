package shedulers

import (
	"HabitsBot/internal/handlers"
	"HabitsBot/internal/models"
	"HabitsBot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

var DaysWeek = map[string]time.Weekday{
	"1": time.Monday,
	"2": time.Tuesday,
	"3": time.Wednesday,
	"4": time.Thursday,
	"5": time.Friday,
	"6": time.Saturday,
	"7": time.Sunday,
}

func AddManyHabitsToTiming(tsds []models.TimeShedulerData, habitBot *handlers.HabitBot) {
	var mu = &sync.Mutex{}
	for _, tsd := range tsds {
		go AddHabitToTiming(tsd.Habit, tsd.Timestamp, habitBot, mu)
	}
}

func AddHabitToTiming(habit models.Habit, ts models.Timestamp, habitBot *handlers.HabitBot, mu *sync.Mutex) {
	day, ok := DaysWeek[ts.Day]

	if !ok {
		log.Error().Msg("не был выбран день недели, have = " + ts.Day + " need = [1, 7]")
		return
	}

	ts.Time = "time__" + ts.Time

	hour, minute, err := utils.ExtractHoursAndMinutes(ts.Time)

	if err != nil {
		log.Error().Err(err).Msg("Не удалось распарсить время в модели Habit")
	}

	chanComplete := make(chan string)
	chanControl := make(chan string)

	mu.Lock()
	habitBot.CompleteChanMap[habit.ID] = &chanComplete
	habitBot.ControlChanMap[habit.ID] = &chanControl
	mu.Unlock()

	ScheduleTaskAWeek(&chanControl, day, hour, minute, func() {
		habitBot.SendNotification(habit)
		timeout := time.Minute * time.Duration(habit.CompletedTime)

		logic := func() {
			msgText := "Вы не успели выполнить привычку " + "<b>" + habit.Title + "</b>!  😤"
			msg := tgbotapi.NewMessage(habit.UserID, msgText)
			msg.ParseMode = tgbotapi.ModeHTML
			habitBot.Bot.Send(msg)

			r := models.Rejection{
				Text:     "Не выполнение привычки вовремя",
				DateTime: time.Now(),
				HabitID:  habit.ID,
			}

			err = habitBot.RejectionStorage.Create(r)

			if err != nil {
				log.Error().Err(err).Msg("ошибка при создании отказа по дедлайну")
			}
		}

		select {
		case <-time.After(timeout):
			logic()
		case <-chanComplete:
			log.Info().Msg("Пользователь выполнил привычку")
			return
		}
	})

	dayWarning, hourWarning, minuteWarning := utils.GetWarningHoursAndMinutes(day, hour, minute, habit.WarningTime)

	ScheduleTaskAWeek(&chanControl, dayWarning, hourWarning, minuteWarning, func() {
		habitBot.SendWarningBeforeNotification(habit)
	})
}
