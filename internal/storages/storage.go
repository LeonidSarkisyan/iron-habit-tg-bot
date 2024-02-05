package storages

import "HabitsBot/internal/models"

type HabitStorage interface {
	Create(habit models.Habit) error
	Get(userID int64) (models.Habit, error)
}
