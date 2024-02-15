package storages

import (
	"HabitsBot/internal/models"
)

type HabitStorage interface {
	Create(habit models.Habit) (int, error)
	Get(userID int64, habitID int) (models.Habit, error)
	GetAll(userID int64, offset int) ([]models.Habit, bool, error)
	Name(habitID int, userID int64) (string, error)
	Habits() ([]models.TimeShedulerData, error)
}

type ExecutionStorage interface {
	Create(execution models.Execution) error
}

type RejectionStorage interface {
	Create(rejection models.Rejection) error
}

type TimestampStorage interface {
	Create(timestamp models.Timestamp) error
	GetByHabitID(habitID int, offset int) ([]models.Timestamp, bool, error)
}
