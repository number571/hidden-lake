package handler

import (
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	"github.com/number571/hidden-lake/internal/utils/language"
)

type sTemplate struct {
	FLanguage      language.ILanguage
	FHeaderAddLink string
	FHeaderAddName [3]string
}

func getTemplate(pCfg config.IConfig) *sTemplate {
	return &sTemplate{
		FLanguage: pCfg.GetSettings().GetLanguage(),
	}
}
