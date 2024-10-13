package handler

import (
	"github.com/number571/hidden-lake/cmd/applications/filesharer/internal/config"
	"github.com/number571/hidden-lake/internal/language"
)

type sTemplate struct {
	FLanguage language.ILanguage
}

func getTemplate(pCfg config.IConfig) *sTemplate {
	return &sTemplate{
		FLanguage: pCfg.GetSettings().GetLanguage(),
	}
}
