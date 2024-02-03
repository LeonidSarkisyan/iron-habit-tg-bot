package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

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
