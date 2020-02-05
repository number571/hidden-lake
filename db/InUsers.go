package db

func InUsers(hashpasw string) bool {
	id := GetUserId(hashpasw)
	if id == -1 {
		return false
	}
	return true
}
