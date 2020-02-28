package db

func InUsers(username string) bool {
	id := GetUserId(username)
	if id == -1 {
		return false
	}
	return true
}
