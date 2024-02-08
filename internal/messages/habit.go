package messages

import (
	"HabitsBot/internal/models"
	"HabitsBot/pkg/utils"
	"fmt"
	"sort"
)

func HabitListMsg(habits []models.Habit) string {
	var msg string

	if len(habits) > 0 {
		msg = "<b>Cписок ваших привычек  ⤵️</b>\n\n"
	} else {
		msg = "Список ваших напоминаний пуст. Добавить новую привычку - " + "/add_new_habit"
	}

	for i, habit := range habits {
		i += 1

		switch i {
		case 1:
			msg += "1️⃣"
		case 2:
			msg += "2️⃣"
		case 3:
			msg += "3️⃣"
		case 4:
			msg += "4️⃣"
		case 5:
			msg += "5️⃣"
		case 6:
			msg += "6️⃣"
		case 7:
			msg += "7️⃣"
		case 8:
			msg += "8️⃣"
		case 9:
			msg += "9️⃣"
		case 10:
			msg += "🔟"
		}

		msg += fmt.Sprintf("  <b>%s</b>\n\n", habit.Title)

		daysListMsg, timesListMsg := daysTimesListMsg(habit.Timestamps)

		msg += daysListMsg + "\n\n"
		msg += timesListMsg + "\n\n"
	}

	return msg
}

func daysTimesListMsg(timestamps []models.Timestamp) (string, string) {
	msgDays := "\t\t\t\t\t\t\t\t📅  <b>Дни недели:</b> "
	msgTimes := "\t\t\t\t\t\t\t\t⏰  <b>Время:</b> "

	var days []string
	var times []string

	for _, timestamp := range timestamps {
		day := timestamp.Day
		time := timestamp.Time

		days = utils.AddUniqueValueToSlice(day, days)
		times = utils.AddUniqueValueToSlice(time, times)
	}

	sort.Strings(times)

	msgDays += listItemsWithCommas(days)
	msgTimes += listItemsWithCommas(times)

	return msgDays, msgTimes
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
