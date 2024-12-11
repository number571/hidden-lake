package handler

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/number571/go-peer/pkg/anonymity/adapters"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/pkg/adapters/tcp"
)

func HandleConfigConnectsAPI(
	pCtx context.Context,
	pWrapper config.IWrapper,
	pLogger logger.ILogger,
	pAdapter adapters.IAdapter,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.GServiceName.Short(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if pR.Method == http.MethodGet {
			connects := pWrapper.GetConfig().GetConnections()
			result := make([]string, 0, len(connects))
			for _, addr := range connects {
				result = append(result, "tcp://"+addr)
			}
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, result)
			return
		}

		connectBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: read connect bytes")
			return
		}

		u, err := url.Parse(strings.TrimSpace(string(connectBytes)))
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("read_connect"))
			_ = api.Response(pW, http.StatusTeapot, "failed: connect is nil")
			return
		}
		if u.Scheme != "tcp" {
			pLogger.PushWarn(logBuilder.WithMessage("scheme_rejected"))
			_ = api.Response(pW, http.StatusAccepted, "rejected: scheme != tcp")
			return
		}

		switch pR.Method {
		case http.MethodPost:
			connects := uniqAppendToSlice(
				pWrapper.GetConfig().GetConnections(),
				u.Host,
			)
			if err := pWrapper.GetEditor().UpdateConnections(connects); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("update_connections"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: update connections")
				return
			}

			networkNode := pAdapter.(tcp.ITCPAdapter).GetConnKeeper().GetNetworkNode()
			_ = networkNode.AddConnection(pCtx, u.Host) // connection may be refused (closed)

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: update connections")
			return

		case http.MethodDelete:
			connects := slices.DeleteFunc(
				pWrapper.GetConfig().GetConnections(),
				func(p string) bool { return p == u.Host },
			)
			if err := pWrapper.GetEditor().UpdateConnections(connects); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("update_connections"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: delete connection")
				return
			}

			networkNode := pAdapter.(tcp.ITCPAdapter).GetConnKeeper().GetNetworkNode()
			_ = networkNode.DelConnection(u.Host) // connection may be refused (closed)

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: delete connection")
		}
	}
}

func uniqAppendToSlice(pSlice []string, pStr string) []string {
	if slices.Contains(pSlice, pStr) {
		return pSlice
	}
	return append(pSlice, pStr)
}
