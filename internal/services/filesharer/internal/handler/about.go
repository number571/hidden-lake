package handler

import (
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	"github.com/number571/hidden-lake/internal/webui"

	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

type sAbout struct {
	*sTemplate
	FAppFullName string
	FDescription [3]string
}

func AboutPage(pLogger logger.ILogger, pCfg config.IConfig) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_filesharer_settings.GetServiceName().Short(), pR)

		if pR.URL.Path != "/about" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = webui.MustParseTemplate("index.html", "about.html").Execute(pW, &sAbout{
			sTemplate:    getTemplate(pCfg),
			FAppFullName: "Hidden Lake Service (Filesharer)",
			FDescription: [3]string{
				"The HLS=filesharer is a file sharing service based on the Anonymous Network Core (HLK) with theoretically provable anonymity. A feature of this file sharing service is the anonymity of the fact of transactions (file downloads), taking into account the existence of a global observer.",
				"HLS=filesharer - это служба обмена файлами, основанная на ядре анонимной сети (HLK) с теоретически доказуемой анонимностью. Особенностью данного файлообменного сервиса является анонимность факта транзакций (скачиваний файла), принимая во внимание существование глобального наблюдателя.",
				"HLS=filesharer estas dosierpartumado servo bazita sur La Anonima Reto Kerno (hlk) kun teorie pruvebla  anonimeco. Karakterizaĵo de ĉi tiu dosierpartuma servo estas la anonimeco de la fakto de transakcioj (dosieraj elŝutoj), konsiderante la ekziston de tutmonda observanto.",
			},
		})
	}
}
