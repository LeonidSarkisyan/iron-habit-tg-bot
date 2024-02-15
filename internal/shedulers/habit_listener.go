package shedulers

import (
	"HabitsBot/internal/handlers"
	"sync"
)

func HabitListener(habitBot *handlers.HabitBot) {
	mu := &sync.Mutex{}
	for data := range habitBot.TimeShedulerChan {
		go AddHabitToTiming(data.Habit, data.Timestamp, habitBot, mu)
	}
}
