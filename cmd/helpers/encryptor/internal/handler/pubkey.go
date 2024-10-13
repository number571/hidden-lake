package handler

import (
	"net/http"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/api"
	http_logger "github.com/number571/hidden-lake/internal/logger/http"

	hle_settings "github.com/number571/hidden-lake/cmd/helpers/encryptor/pkg/settings"
)

func HandleServicePubKeyAPI(pLogger logger.ILogger, pPubKey asymmetric.IPubKey) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hle_settings.CServiceName, pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, pPubKey.ToString())
	}
}