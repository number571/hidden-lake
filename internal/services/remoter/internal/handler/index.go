package handler

import (
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	hls_settings "github.com/number571/hidden-lake/internal/services/remoter/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleIndexAPI(pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)
		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))

		_ = api.Response(pW, http.StatusOK, hls_settings.CAppFullName)
	}
}
