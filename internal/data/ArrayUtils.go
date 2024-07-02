package data

func IsEmptyStringArray(strArray []string) bool {
	for _, str := range strArray {
		if str != "" {
			return false
		}
	}
	return true
}

func RemoveEmptyStringArrays(arr [][]string) [][]string {
	var result [][]string
	for _, innerArray := range arr {
		if !IsEmptyStringArray(innerArray) {
			result = append(result, innerArray)
		}
	}
	return result
}
