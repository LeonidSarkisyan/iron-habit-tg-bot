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

		warningTime := 15

		if minute-warningTime < 0 {
			hour -= 1
			minute = minute - warningTime + 60
		}

		ScheduleTask(day, hour, 42-warningTime, func() {
			habitBot.SendWarningBeforeNotification(habit)
		})

		ScheduleTask(day, hour, 32, func() {
			habitBot.SendNotification(habit)
		})
	}

}
