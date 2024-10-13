package handler

import (
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	pkg_settings "github.com/number571/hidden-lake/cmd/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/api"

	http_logger "github.com/number571/hidden-lake/internal/logger/http"
)

func HandleNetworkPubKeyAPI(pLogger logger.ILogger, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.CServiceName, pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, pNode.GetMessageQueue().GetClient().GetPubKey().ToString())
	}
}
