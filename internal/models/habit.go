package models

type Habit struct {
	Title      string
	UserID     int64
	Timestamps []Timestamp
}

type Timestamp struct {
	Day     string
	Time    string
	HabitID int
}

func NewTimestamps(days []string, times []string) []Timestamp {
	var timestamps []Timestamp

	for _, day := range days {
		for _, time := range times {
			timestamps = append(timestamps, Timestamp{
				Day:  day,
				Time: time,
			})
		}
	}

	return timestamps
}
