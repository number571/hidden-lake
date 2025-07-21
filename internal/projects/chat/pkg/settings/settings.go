package settings

import (
	"github.com/number571/hidden-lake/internal/projects"
	"github.com/number571/hidden-lake/internal/utils/name"
)

var (
	gAppName = name.LoadAppName(CProjectFullName)
)

func GetAppName() name.IAppName {
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
