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
	"Понедельник": time.Monday,
	"Вторник":     time.Tuesday,
	"Среда":       time.Wednesday,
	"Четверг":     time.Thursday,
	"Пятница":     time.Friday,
	"Суббота":     time.Saturday,
	"Воскресенье": time.Sunday,
}

func AddManyHabitsToTiming(habits []models.Habit, habitBot *handlers.HabitBot) {
	var mu = &sync.Mutex{}
	for _, habit := range habits {
		mu.Lock()
		go AddHabitToTiming(habit, habitBot)
		mu.Unlock()
	}
}

func AddHabitToTiming(habit models.Habit, habitBot *handlers.HabitBot) {

	for _, ts := range habit.Timestamps {
		day, ok := DaysWeek[ts.Day]

		if !ok {

			// todo: добавить логирование ошибки при невалидном дне недели в модели Habit и в бд
			continue
		}

		ts.Time = "time__" + ts.Time

		hour, minute, err := utils.ExtractHoursAndMinutes(ts.Time)

		if err != nil {
			log.Error().Err(err).Msg("Не удалось распарсить время в модели Habit")
		}

		chanControl := make(chan string)
		chanComplete := make(chan string)

		habitBot.CompleteChanMap[habit.ID] = &chanComplete
		habitBot.ControlChanMap[habit.ID] = &chanControl

		ScheduleTaskAWeek(&chanControl, day, hour, minute, func() {
			habitBot.SendNotification(habit)
			f := func() {
				msgText := "Вы не успели выполнить привычку " + "<b>" + habit.Title + "</b>!  😤"
				msg := tgbotapi.NewMessage(habit.UserID, msgText)
				msg.ParseMode = tgbotapi.ModeHTML
				habitBot.Bot.Send(msg)
			}
			delay := time.Duration(habit.CompletedTime) * time.Minute
			time.AfterFunc(delay, f)
		})

		dayWarning, hourWarning, minuteWarning := utils.GetWarningHoursAndMinutes(day, hour, minute, habit.WarningTime)

		ScheduleTaskAWeek(&chanControl, dayWarning, hourWarning, minuteWarning, func() {
			habitBot.SendWarningBeforeNotification(habit)
		})
	}
}
