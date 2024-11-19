package handler

import (
	"html/template"
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/filesharer/pkg/app/config"
	"github.com/number571/hidden-lake/internal/webui"

	hlf_settings "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

type sAbout struct {
	*sTemplate
	FAppFullName string
	FDescription [3]string
}

func AboutPage(pLogger logger.ILogger, pCfg config.IConfig) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlf_settings.CServiceName, pR)

		if pR.URL.Path != "/about" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		t, err := template.ParseFS(
			webui.GetTemplatePath(),
			"index.html",
			"about.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = t.Execute(pW, &sAbout{
			sTemplate:    getTemplate(pCfg),
			FAppFullName: "Hidden Lake Filesharer",
			FDescription: [3]string{
				"The HLF is a file sharing service based on the Anonymous Network Core (HLS) with theoretically provable anonymity. A feature of this file sharing service is the anonymity of the fact of transactions (file downloads), taking into account the existence of a global observer.",
				"HLF - это служба обмена файлами, основанная на ядре анонимной сети (HLS) с теоретически доказуемой анонимностью. Особенностью данного файлообменного сервиса является анонимность факта транзакций (скачиваний файла), принимая во внимание существование глобального наблюдателя.",
				"HLF estas dosierpartumado servo bazita sur La Anonima Reto Kerno (hls) kun teorie pruvebla  anonimeco. Karakterizaĵo de ĉi tiu dosierpartuma servo estas la anonimeco de la fakto de transakcioj (dosieraj elŝutoj), konsiderante la ekziston de tutmonda observanto.",
			},
		})
	}
}
