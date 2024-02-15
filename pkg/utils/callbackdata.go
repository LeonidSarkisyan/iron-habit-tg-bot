package utils

import (
	"strconv"
	"strings"
)

func GetIDFromCallbackData(prefix, callBackData string) (int, error) {
	id := strings.Replace(callBackData, prefix, "", 1)

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return 0, err
	}

	return idInt, nil
}
