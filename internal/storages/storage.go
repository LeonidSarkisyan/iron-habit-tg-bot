package storages

import "HabitsBot/internal/models"

type HabitStorage interface {
	Create(habit models.Habit) (int, error)
	Get(userID int64) (models.Habit, error)
	Name(habitID int, userID int64) (string, error)
}

type RejectionStorage interface {
	Create(rejection models.Rejection) error
}
