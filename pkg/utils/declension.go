package utils

func DeclensionDays(n int) string {
	return declension(n, "день", "дня", "дней")
}

func DeclensionHours(n int) string {
	return declension(n, "час", "часа", "часов")
}

func DeclensionMinutes(n int) string {
	return declension(n, "минута", "минуты", "минут")
}

func DeclensionSeconds(n int) string {
	return declension(n, "секунда", "секунды", "секунд")
}

func declension(n int, form1, form2, form3 string) string {
	if n%100 >= 11 && n%100 <= 14 {
		return form3
	}
	switch n % 10 {
	case 1:
		return form1
	case 2, 3, 4:
		return form2
	default:
		return form3
	}
}
