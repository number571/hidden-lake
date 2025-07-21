package settings

import (
	"github.com/number571/hidden-lake/internal/projects"
	"github.com/number571/hidden-lake/internal/utils/appname"
)

var (
	gAppName = appname.LoadAppName(CProjectFullName)
)

func GetAppName() appname.IAppName {
	return gAppName
}

const (
	CProjectName = "chat"
)

const (
	CProjectFullName    = projects.CProjectFullPrefix + "=" + CProjectName
	CProjectDescription = "console anonymous group chat"
)

const (
	CPathDB = "hlp-chat.db"
)
