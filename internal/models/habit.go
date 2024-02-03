package models

type Habit struct {
	Name string
}

func NewHabit(name string) *Habit {
	return &Habit{
		Name: name,
	}
}
