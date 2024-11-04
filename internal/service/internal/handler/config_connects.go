package handler

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/hidden-lake/internal/service/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/slices"
)

func HandleConfigConnectsAPI(
	pCtx context.Context,
	pWrapper config.IWrapper,
	pLogger logger.ILogger,
	pNode anonymity.INode,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.CServiceName, pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if pR.Method == http.MethodGet {
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, pWrapper.GetConfig().GetConnections())
			return
		}

		connectBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: read connect bytes")
			return
		}

		connect := strings.TrimSpace(string(connectBytes))
		if connect == "" {
			pLogger.PushWarn(logBuilder.WithMessage("read_connect"))
			_ = api.Response(pW, http.StatusTeapot, "failed: connect is nil")
			return
		}

		switch pR.Method {
		case http.MethodPost:
			connects := slices.UniqAppendToSlice(
				pWrapper.GetConfig().GetConnections(),
				connect,
			)
			if err := pWrapper.GetEditor().UpdateConnections(connects); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("update_connections"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: update connections")
				return
			}

			_ = pNode.GetNetworkNode().AddConnection(pCtx, connect) // connection may be refused (closed)

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: update connections")
			return

		case http.MethodDelete:
			connects := slices.DeleteFromSlice(pWrapper.GetConfig().GetConnections(), connect)
			if err := pWrapper.GetEditor().UpdateConnections(connects); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("update_connections"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: delete connection")
				return
			}

			_ = pNode.GetNetworkNode().DelConnection(connect) // connection may be refused (closed)

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: delete connection")
		}
	}
}
