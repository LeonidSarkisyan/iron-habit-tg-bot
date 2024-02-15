package messages

import (
	"HabitsBot/pkg/utils"
	"fmt"
	"strconv"
)

const (
	ErrorCreateReminder = "Ошибка при создании напоминания. Попробуйте ещё раз."
)

func CreateReminderMsg(habitName string, day string, time string) string {
	weekDay, ok := daysOfWeek[day]

	if !ok {
		return ""
	}

	headerMsg := fmt.Sprintf("Напоминание для привычки <b>%s</b> создано  🎉\n\n", habitName)
	aboutMsg := fmt.Sprintf("<b>Время:</b> %s, %s", weekDay, time)

	return headerMsg + aboutMsg
}

func TimeWhenDoMsg(days, hours, minutes int) string {
	days_ := strconv.Itoa(days)
	hours_ := strconv.Itoa(hours)
	minutes_ := strconv.Itoa(minutes)

	baseMsg := "Первое напоминание о привычке через "

	if days > 0 {
		baseMsg += days_ + " " + utils.DeclensionDays(days) + " "
	}

	if hours > 0 {
		baseMsg += hours_ + " " + utils.DeclensionHours(hours) + " "
	}

	if minutes > 0 {
		baseMsg += minutes_ + " " + utils.DeclensionMinutes(minutes)
	}

	return baseMsg
}
