package utils

func RemoveByElem(slice []string, elem string) []string {
	var position = len(slice) - 1
	for position >= 0 {
		if slice[position] == elem {
			slice = RemoveByIndex(slice, position)
		}
		position--
	}
	return slice
}
