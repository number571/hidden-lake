package handler

import (
	"github.com/number571/hidden-lake/internal/applications/filesharer/internal/config"
	"github.com/number571/hidden-lake/internal/modules/language"
)

type sTemplate struct {
	FLanguage language.ILanguage
}

func getTemplate(pCfg config.IConfig) *sTemplate {
	return &sTemplate{
		FLanguage: pCfg.GetSettings().GetLanguage(),
	}
}
