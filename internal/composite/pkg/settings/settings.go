package settings

import "github.com/number571/hidden-lake/internal/utils/appname"

var (
	gShortAppName = appname.ToShortAppName(CAppFullName)
)

func GetShortAppName() string {
	return gShortAppName
}

const (
	CAppFullName    = "hidden-lake-composite"
	CAppDescription = "runs many HL applications as one application"
)

const (
	CPathYML = "hlc.yml"
)
