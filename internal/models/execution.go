package models

import "time"

type Execution struct {
	ID               int
	DatetimeComplete time.Time
	HabitID          int
}
