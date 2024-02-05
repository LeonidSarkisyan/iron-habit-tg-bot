package shedulers

import (
	"github.com/rs/zerolog/log"
	"time"
)

func ScheduleTask(weekday time.Weekday, hour, minute int, task func()) {
	go func() {
		now := time.Now()
		currentWeekday := now.Weekday()

		daysUntilNextWeekday := (int(weekday) - int(currentWeekday) + 7) % 7
		nextWeekday := now.Add(time.Duration(daysUntilNextWeekday) * 24 * time.Hour)
		nextTime := time.Date(nextWeekday.Year(), nextWeekday.Month(), nextWeekday.Day(), hour, minute, 0, 0, nextWeekday.Location())

		timeUntilFirstExecution := nextTime.Sub(now)
		timer := time.NewTimer(timeUntilFirstExecution)

		taskFunctionWrapper := func() {
			task()

			timer.Reset(7 * 24 * time.Hour)
		}

		for {
			<-timer.C
			taskFunctionWrapper()
		}
	}()
	event := log.Info().Str("weekday", weekday.String()).Int("hour", hour).Int("minute", minute)
	event.Msg("Задание запланировано")
}
