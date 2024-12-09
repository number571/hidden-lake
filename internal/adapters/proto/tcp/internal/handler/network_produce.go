package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/anonymity/adapters"
	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/hidden-lake/internal/adapters/proto/tcp/pkg/app/config"
	hla_settings "github.com/number571/hidden-lake/internal/adapters/proto/tcp/pkg/settings"
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

		sett := pConfig.GetSettings()
		msgLen := uint64(sett.GetMessageSizeBytes()+net_message.CMessageHeadSize) << 1 // nolint: unconvert
		msgStr := make([]byte, msgLen)
		n, err := io.ReadFull(pR.Body, msgStr)
		if err != nil || uint64(n) != msgLen {
			pW.WriteHeader(http.StatusBadRequest)
			return
		}

		msg, err := net_message.LoadMessage(pConfig.GetSettings(), string(msgStr))
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("load_message"))
			pW.WriteHeader(http.StatusNotAcceptable)
			return
		}

		if err := pProducer.Produce(pCtx, msg); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("produce_message"))
			pW.WriteHeader(http.StatusBadGateway)
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, hla_settings.CServiceFullName)
	}
}
