package handler

import (
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	pkg_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleIndexAPI(pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.GetAppShortNameFMT(), pR)
		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))

		_ = api.Response(pW, http.StatusOK, pkg_settings.CAppFullName)
	}
}
