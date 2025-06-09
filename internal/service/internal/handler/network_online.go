package handler

import (
	"context"
	"io"
	"net/http"
	"sort"

	"github.com/number571/go-peer/pkg/logger"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/pkg/adapters/http/client"
)

func HandleNetworkOnlineAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pEPClients []client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.GetServiceName().Short(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		switch pR.Method {
		case http.MethodGet:
			inOnline := make([]string, 0, 128)
			for _, client := range pEPClients {
				gotConns, err := client.GetOnlines(pCtx)
				if err != nil {
					pLogger.PushWarn(logBuilder.WithMessage("get_connections"))
					_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: get online connections")
					return
				}
				inOnline = append(inOnline, gotConns...)
			}

			sort.SliceStable(inOnline, func(i, j int) bool {
				return inOnline[i] < inOnline[j]
			})

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, inOnline)
		case http.MethodDelete:
			connectBytes, err := io.ReadAll(pR.Body)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
				_ = api.Response(pW, http.StatusConflict, "failed: read connect bytes")
				return
			}

			for _, client := range pEPClients {
				if err := client.DelOnline(pCtx, string(connectBytes)); err != nil {
					pLogger.PushWarn(logBuilder.WithMessage("del_connection"))
					_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: delete online connections")
					return
				}
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: delete online connection")
		}
	}
}
