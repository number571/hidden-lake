package handler

import (
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/messenger/pkg/app/config"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func IndexPage(pLogger logger.ILogger, pCfg config.IConfig) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.GServiceName.Short(), pR)

		if pR.URL.Path != "/" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		http.Redirect(pW, pR, "/about", http.StatusFound)
	}
}
