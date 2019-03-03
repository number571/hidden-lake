package utils

func CheckWarning(warn error) {
	if warn != nil {
		PrintWarning(warn.Error())
	}
}
