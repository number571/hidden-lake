package settings

import "strings"

var (
	gAppShortNameFMT = strings.ToUpper(CAppShortName)
)

func GetAppShortNameFMT() string {
	return gAppShortNameFMT
}

const (
	CAppShortName   = "hlc"
	CAppFullName    = "hidden-lake-composite"
	CAppDescription = "runs many HL applications as one application"
)

const (
	CPathYML = CAppShortName + ".yml"
)
