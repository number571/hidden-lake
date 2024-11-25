package handler

import (
	"github.com/number571/hidden-lake/internal/applications/filesharer/pkg/app/config"
	"github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/language"
)

type sTemplate struct {
	FAppName  string
	FLanguage language.ILanguage
}

func getTemplate(pCfg config.IConfig) *sTemplate {
	return &sTemplate{
		FAppName:  settings.GServiceName.Short(),
		FLanguage: pCfg.GetSettings().GetLanguage(),
	}
}
