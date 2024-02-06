package shedulers

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"time"
)

func ScheduleTaskAWeek(chanControl *chan string, weekday time.Weekday, hour, minute int, task func()) {
	go func() {
		now := time.Now()
		currentWeekday := now.Weekday()

		daysUntilNextWeekday := (int(weekday) - int(currentWeekday) + 7) % 7
		nextWeekday := now.Add(time.Duration(daysUntilNextWeekday) * 24 * time.Hour)
		nextTime := time.Date(nextWeekday.Year(), nextWeekday.Month(), nextWeekday.Day(), hour, minute, 0, 0, nextWeekday.Location())

		timeUntilFirstExecution := nextTime.Sub(now)

		if timeUntilFirstExecution < 0 {
			timeUntilFirstExecution += 7 * 24 * time.Hour
		}

		fmt.Println(timeUntilFirstExecution)

		timer := time.NewTimer(timeUntilFirstExecution)

		taskFunctionWrapper := func() {
			task()

			timer.Reset(7 * 24 * time.Hour)
		}

		for {
			select {
			case <-timer.C:
				taskFunctionWrapper()
			case operation := <-*chanControl:
				log.Info().Str("operation", operation).Msg("Пользователь делает операцию: " + operation)
				switch operation {
				case "cancel":
					nextWeekday = getNextWeekday(now, weekday)
					nextTime = time.Date(nextWeekday.Year(), nextWeekday.Month(), nextWeekday.Day(), hour, minute, 0, 0, nextWeekday.Location())
					timeUntilNextExecution := nextTime.Sub(now)
					timer.Reset(timeUntilNextExecution)
					log.Info().Msg("Пользователь отменил привычку, теперь привычка сработает через ")
					fmt.Println(timeUntilNextExecution)
				case "delete":
					return
				}
			}
		}
	}()
	event := log.Info().Str("weekday", weekday.String()).Int("hour", hour).Int("minute", minute)
	event.Msg("Задание запланировано")
}

func getNextWeekday(currentTime time.Time, targetWeekday time.Weekday) time.Time {
	daysToAdd := (int(targetWeekday) - int(currentTime.Weekday()) + 7) % 7
	if daysToAdd == 0 {
		daysToAdd = 7 // Если текущий день недели уже является указанным днем недели, добавляем 7 дней
	}
	return currentTime.AddDate(0, 0, daysToAdd)
}
