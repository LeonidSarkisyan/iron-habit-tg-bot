package messages

import (
	"sort"
	"strings"
)

const (
	StartMsg          = "Привет! Добро пожаловать! Я бот для управления привычками."
	InputHabitNameMsg = "Введите название новой привычки:"

	InputHabitDaysMsg = "Выберите дни недели, когда вам напоминать о привычке:"
	InputHabitTimeMsg = "Выберите время дня, в которое необходимо напоминать вам о новом событии:\n\nP.S. Время указано по МСК"
)

var daysOfWeek = map[string]string{
	"1": "Понедельник",
	"2": "Вторник",
	"3": "Среда",
	"4": "Четверг",
	"5": "Пятница",
	"6": "Суббота",
	"7": "Воскресенье",
}

func ShowHabitNameAndDaysMsg(name string, days []string) string {

	sort.Strings(days)

	stringsNumberToDays(&days)

	return "Название привычки: " + name + "\n" + "Дни недели: " + strings.Join(days, ", ")
}

func ShowSaveHabitMsg(name string, days []string, times []string) string {

	sort.Strings(times)

	return "<b>Привычка \"" + name + "\" успешно создана!</b>\n\n" + "<b>Дни недели:</b> " + strings.Join(days, ", ") + "\n\n" + "<b>Время:</b> " + strings.Join(times, ", ")
}

func stringsNumberToDays(stringsNumbers *[]string) {
	for i, number := range *stringsNumbers {
		(*stringsNumbers)[i] = daysOfWeek[number]
	}
}
