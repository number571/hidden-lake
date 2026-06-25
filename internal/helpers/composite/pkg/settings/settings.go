package settings

import (
	"strings"

	"github.com/number571/hidden-lake/internal/helpers"
)

var (
	gAppShortNameFMT = strings.ToUpper(CAppShortName)
)

func GetAppShortNameFMT() string {
	return gAppShortNameFMT
}

const (
	CAppShortName = helpers.CHelperShortName + "-" + CAppHelperName
	CAppFullName  = helpers.CHelperDomainName + "=" + CAppHelperName
)

const (
	CAppHelperName  = "composite"
	CAppDescription = "runs many HL applications as one application"
)

const (
	CPathYML = CAppShortName + ".yml"
)
