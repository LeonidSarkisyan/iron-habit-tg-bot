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
