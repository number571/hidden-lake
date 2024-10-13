package handler

import (
	"github.com/number571/hidden-lake/internal/applications/messenger/internal/config"
	"github.com/number571/hidden-lake/internal/utils/language"
)

type sTemplate struct {
	FLanguage language.ILanguage
}

func getTemplate(pCfg config.IConfig) *sTemplate {
	return &sTemplate{
		FLanguage: pCfg.GetSettings().GetLanguage(),
	}
}
