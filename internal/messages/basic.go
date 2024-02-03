package messages

import (
	"HabitsBot/internal/keyboards"
	"strings"
)

const (
	StartMsg          = "Привет! Добро пожаловать! Я бот для управления привычками."
	InputHabitNameMsg = "Введите название новой привычки:"

	InputHabitDaysMsg = "Выберите дни недели, когда вам напоминать о привычке:"
	InputHabitTimeMsg = "Выберите время дня, в которое необходимо напоминать вам о новом событии:\n\nP.S. Время указано по МСК"
)

func ShowHabitNameAndDaysMsg(name string, days []string) string {

	daysOfWeek := []keyboards.CallBackData{
		{Name: "Понедельник", Data: "1"},
		{Name: "Вторник", Data: "2"},
		{Name: "Среда", Data: "3"},
		{Name: "Четверг", Data: "4"},
		{Name: "Пятница", Data: "5"},
		{Name: "Суббота", Data: "6"},
		{Name: "Воскресенье", Data: "7"},
	}

	var daysName []string

	for _, day := range days {
		for _, dayOfWeek := range daysOfWeek {
			if day == dayOfWeek.Data {
				daysName = append(daysName, dayOfWeek.Name)
				break
			}
		}
	}

	return "Название привычки: " + name + "\n" + "Дни недели: " + strings.Join(daysName, ", ")
}
