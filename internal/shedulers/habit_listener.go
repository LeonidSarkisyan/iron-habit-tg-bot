package shedulers

import "HabitsBot/internal/handlers"

func HabitListener(habitBot *handlers.HabitBot) {
	for habit := range habitBot.TimeShedulerChan {
		go AddHabitToTiming(habit, habitBot)
	}
}
