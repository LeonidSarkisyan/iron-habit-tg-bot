package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const (
	MoscowTimezone = "Europe/Moscow"
)

var WeekDays = map[int]time.Weekday{
	1: time.Monday,
	2: time.Tuesday,
	3: time.Wednesday,
	4: time.Thursday,
	5: time.Friday,
	6: time.Saturday,
	7: time.Sunday,
}

func WeekDayFromInt(day int) (time.Weekday, error) {
	weekDay, ok := WeekDays[day]

	if !ok {
		return 0, fmt.Errorf("неверный номер дня недели: have = %d, want = [0, 7]", day)
	}

	return weekDay, nil
}

func ExtractHoursAndMinutes(input string) (int, int, error) {
	re := regexp.MustCompile(`time__(\d{2}):(\d{2})`)
	matches := re.FindStringSubmatch(input)

	if len(matches) != 3 {
		return 0, 0, fmt.Errorf("неверный формат строки времени: %s", input)
	}

	hours, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, fmt.Errorf("невозможно преобразовать часы в число: %v", err)
	}

	minutes, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, fmt.Errorf("невозможно преобразовать минуты в число: %v", err)
	}

	return hours, minutes, nil
}

func GetWarningHoursAndMinutes(day time.Weekday, hour, minute, warningTime int) (time.Weekday, int, int) {

	if minute-warningTime < 0 {
		minute = minute - warningTime + 60
		if hour == 0 {
			hour = 23
			if day == 0 {
				day = 6
			} else {
				day -= 1
			}
		} else {
			hour -= 1
		}
	} else {
		minute -= warningTime
	}

	return day, hour, minute
}

func TimeUntil(dayOfWeek time.Weekday, hour, minute int) (days, hours, minutes int) {
	now := time.Now()
	currentWeekDay := now.Weekday()
	currentHour := now.Hour()
	currentMinute := now.Minute()

	if dayOfWeek == 0 {
		dayOfWeek = 7
	}

	if currentWeekDay == 0 {
		currentWeekDay = 7
	}

	days = int(dayOfWeek - currentWeekDay)
	hours = hour - currentHour
	minutes = minute - currentMinute

	minutesSmallZero := false
	hoursSmallZero := false

	if minutes < 0 {
		minutes += 60
		minutesSmallZero = true
	}

	if minutesSmallZero {
		hours -= 1
	}

	if hours < 0 {
		hours += 24
		hoursSmallZero = true
	}

	if hoursSmallZero {
		days -= 1
	}

	if days < 0 {
		days += 7
	}

	return days, hours, minutes
}
