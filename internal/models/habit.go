package models

type TimeShedulerData struct {
	Habit
	Timestamp
}

type Habit struct {
	ID            int
	Title         string
	UserID        int64
	WarningTime   int
	CompletedTime int
	Timestamps    []Timestamp
}

type Timestamp struct {
	ID      int
	Day     string
	Time    string
	HabitID int
}
