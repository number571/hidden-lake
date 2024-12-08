package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/anonymity/adapters"
	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/app/config"
	hla_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleNetworkProduceAPI(
	pCtx context.Context,
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pProducer adapters.IProducer,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hla_settings.GServiceName.Short(), pR)

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			pW.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		msgStr, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			pW.WriteHeader(http.StatusBadRequest)
			return
		}

		msg, err := net_message.LoadMessage(pConfig.GetSettings(), string(msgStr))
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("load message"))
			pW.WriteHeader(http.StatusNotAcceptable)
			return
		}

		_ = pProducer.Produce(pCtx, msg)

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, hla_settings.CServiceFullName)
	}
}
