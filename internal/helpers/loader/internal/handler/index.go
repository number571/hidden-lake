package handler

import (
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	hll_settings "github.com/number571/hidden-lake/internal/helpers/loader/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleIndexAPI(pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hll_settings.GServiceName.Short(), pR)
		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))

		_ = api.Response(pW, http.StatusOK, hll_settings.CServiceFullName)
	}
}
