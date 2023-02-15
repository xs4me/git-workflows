package utils

func IndexOf(search string, data []string) int {
	for k, v := range data {
		if search == v {
			return k
		}
	}
	return -1
}
