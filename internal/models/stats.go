package models

import "time"

type Rejection struct {
	ID       int
	Text     string
	DateTime time.Time
	HabitID  int
}
