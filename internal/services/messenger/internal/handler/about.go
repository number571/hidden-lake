package handler

import (
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/app/config"
	"github.com/number571/hidden-lake/internal/webui"

	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

type sAbout struct {
	*sTemplate
	FAppFullName string
	FDescription [3]string
}

func AboutPage(pLogger logger.ILogger, pCfg config.IConfig) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_messenger_settings.GetServiceName().Short(), pR)

		if pR.URL.Path != "/about" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = webui.MustParseTemplate("index.html", "about.html").Execute(pW, &sAbout{
			sTemplate:    getTemplate(pCfg),
			FAppFullName: "Hidden Lake Service (Messenger)",
			FDescription: [3]string{
				"The HLS=messenger is a messenger based on the core of an anonymous network with theoretically provable anonymity of HLK. A feature of this messenger is the provision of anonymity of the fact of transactions (sending, receiving).",
				"HLS=messenger - это мессенджер, основанный на ядре анонимной сети (HLK) с теоретически доказуемой анонимностью. Особенностью данного мессенджера является анонимность факта совершения транзакций (отправление, получение).",
				"HLS=messenger estas mesaĝisto bazita sur La Anonima Retkerno (HLK) kun teorie pruvebla anonimeco. La propreco de ĉi tiu mesaĝisto estas la anonimeco de la fakto de komisiono transakcioj (sendo, ricevo).",
			},
		})
	}
}
