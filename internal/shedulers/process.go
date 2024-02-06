package shedulers

import (
	"HabitsBot/internal/handlers"
	"HabitsBot/internal/models"
	"HabitsBot/pkg/utils"
	"github.com/rs/zerolog/log"
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

		habitBot.ControlChanMap[habit.ID] = &chanControl

		minute = 46

		ScheduleTaskAWeek(&chanControl, day, hour, minute, func() {
			habitBot.SendNotification(habit)
		})

		dayWarning, hourWarning, minuteWarning := utils.GetWarningHoursAndMinutes(day, hour, minute, habit.WarningTime)

		ScheduleTaskAWeek(&chanControl, dayWarning, hourWarning, minuteWarning, func() {
			habitBot.SendWarningBeforeNotification(habit)
		})
	}
}
