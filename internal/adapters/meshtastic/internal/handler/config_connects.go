package handler

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/adapters/meshtastic/pkg/app/config"
	hla_settings "github.com/number571/hidden-lake/internal/adapters/meshtastic/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleConfigConnectsAPI(
	pCtx context.Context,
	pWrapper config.IConfig,
	pLogger logger.ILogger,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hla_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if pR.Method == http.MethodGet {
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, []string{})
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
		if u.Scheme != hla_settings.CAppAdapterName {
			pLogger.PushWarn(logBuilder.WithMessage("scheme_rejected"))
			_ = api.Response(pW, http.StatusAccepted, "rejected: scheme != meshtastic")
			return
		}

		switch pR.Method {
		case http.MethodPost:
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusNoContent, "not supported")
			return
		case http.MethodDelete:
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusNoContent, "not supported")
		}
	}
}
