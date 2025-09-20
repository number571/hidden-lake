package handler

import (
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/number571/go-peer/pkg/logger"
	pkg_settings "github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	hla_http "github.com/number571/hidden-lake/pkg/adapters/http"
)

func HandleNetworkOnlineAPI(
	pLogger logger.ILogger,
	pAdapter hla_http.IHTTPAdapter,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		switch pR.Method {
		case http.MethodGet:
			connects := pAdapter.GetOnlines()
			inOnline := make([]string, 0, len(connects))
			for _, addr := range connects {
				inOnline = append(inOnline, pkg_settings.CAdapterName+"://"+addr)
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
			if u.Scheme != pkg_settings.CAdapterName {
				pLogger.PushWarn(logBuilder.WithMessage("scheme_rejected"))
				_ = api.Response(pW, http.StatusAccepted, "rejected: scheme != tcp")
				return
			}

			// delete not supported in http adapter
			_ = u

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusNoContent, "not supported")
		}
	}
}
