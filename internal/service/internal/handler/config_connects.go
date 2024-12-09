package handler

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/adapters/pkg/client"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleConfigConnectsAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pHlaClient client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.GServiceName.Short(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if pR.Method == http.MethodGet {
			connects, err := pHlaClient.GetConnections(pCtx)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("get_connections"))
				_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: get connections")
				return
			}
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, connects)
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
			if err := pHlaClient.AddConnection(pCtx, connect); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("add_connections"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: add connections")
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: add connections")
			return

		case http.MethodDelete:
			if err := pHlaClient.DelConnection(pCtx, connect); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("del_connection"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: del connection")
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: del connection")
		}
	}
}

// func uniqAppendToSlice(pSlice []string, pStr string) []string {
// 	if slices.Contains(pSlice, pStr) {
// 		return pSlice
// 	}
// 	return append(pSlice, pStr)
// }
