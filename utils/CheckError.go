package utils

func CheckError(err error) {
	if err != nil {
		PrintError(err.Error())
	}
}
