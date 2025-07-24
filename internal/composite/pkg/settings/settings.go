package settings

import "github.com/number571/hidden-lake/pkg/utils/appname"

var (
	gFmtAppName = appname.LoadAppName(CAppFullName)
)

func GetFmtAppName() appname.IFmtAppName {
	return gFmtAppName
}

const (
	CAppFullName    = "hidden-lake-composite"
	CAppDescription = "runs many HL applications as one application"
)

const (
	CPathYML = "hlc.yml"
)
