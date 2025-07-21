package handler

import (
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/app/config"
	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	"github.com/number571/hidden-lake/internal/webui"

	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

type sError struct {
	*sTemplate
	FTitle   string
	FMessage string
}

func ErrorPage(pLogger logger.ILogger, pCfg config.IConfig, pTitle, pMessage string) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_messenger_settings.GetAppName().Short(), pR)

		pW.WriteHeader(http.StatusNotFound)

		pLogger.PushWarn(logBuilder.WithMessage(pTitle))
		_ = webui.MustParseTemplate("index.html", "error.html").Execute(pW, &sError{
			sTemplate: getTemplate(pCfg),
			FTitle:    pTitle,
			FMessage:  pMessage,
		})
	}
}

func NotFoundPage(pLogger logger.ILogger, pCfg config.IConfig) http.HandlerFunc {
	return ErrorPage(pLogger, pCfg, "404_page", "page not found")
}
