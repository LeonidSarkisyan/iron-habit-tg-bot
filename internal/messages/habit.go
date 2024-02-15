package messages

import (
	"HabitsBot/internal/models"
	"fmt"
)

var (
	CancelCreateHabitMsg = "Вы отменили создание привычки."
	ErrorCreateHabitMsg  = "Произошла ошибка при создании привычки. Попробуйте ещё раз."
)

func BeforeCreateHabitMsg(name string) string {
	return fmt.Sprintf("Название вашей привычки: <b>%s</b>. Создаём?", name)
}

func HabitCreatedMsg(name string) string {
	return "<b>Привычка \"" + name + "\" успешно создана!</b>  🥳"
}

func HabitListMsg(habits []models.Habit) string {
	var msg string

	if len(habits) > 0 {
		msg = "<b>Cписок ваших привычек  ⤵️</b>\n\n" +
			"Если хотите изменить привычку, нажмите на ее название в списке"
	} else {
		msg = "Список ваших напоминаний пуст. Добавить новую привычку - " + "/add_new_habit"
	}

	return msg
}

func HabitDetailMsg(habit models.Habit) string {
	return "Ваша привычка: <b>" + habit.Title + "</b>\n\n" +
		"Что вы хотите сделать с ней?"
}

func listItemsWithCommas(items []string) string {
	msg := ""

	for i, item := range items {
		msg += item
		if i != len(items)-1 {
			msg += ", "
		}
	}

	return msg
}
