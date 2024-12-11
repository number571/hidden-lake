package handler

import (
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/number571/go-peer/pkg/anonymity/adapters"
	"github.com/number571/go-peer/pkg/logger"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/pkg/adapters/tcp"
)

func HandleNetworkOnlineAPI(
	pLogger logger.ILogger,
	pAdapter adapters.IAdapter,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.GServiceName.Short(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		switch pR.Method {
		case http.MethodGet:
			networkNode := pAdapter.(tcp.ITCPAdapter).GetConnKeeper().GetNetworkNode()
			connects := networkNode.GetConnections()

			inOnline := make([]string, 0, len(connects))
			for addr := range connects {
				inOnline = append(inOnline, "tcp://"+addr)
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

			networkNode := pAdapter.(tcp.ITCPAdapter).GetConnKeeper().GetNetworkNode()
			if err := networkNode.DelConnection(u.Host); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("del_connection"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: delete online connection")
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: delete online connection")
		}
	}
}
