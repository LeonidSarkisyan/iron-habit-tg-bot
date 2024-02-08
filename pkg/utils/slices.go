package utils

func RemoveDuplicates(slice []int) []int {
	encountered := make(map[int]bool)
	var result []int

	for _, v := range slice {
		if encountered[v] == false {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}

func ConvertToIntSlice(slice []interface{}) []int {
	result := make([]int, len(slice))

	for i, v := range slice {
		result[i] = v.(int)
	}

	return result
}

func AddUniqueValueToSlice[T string | int](value T, slice []T) []T {
	exists := false

	for _, v := range slice {
		if v == value {
			exists = true
			return slice
		}
	}

	if !exists {
		slice = append(slice, value)
	}

	return slice
}
