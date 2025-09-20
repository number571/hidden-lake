package handler

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/adapters/http/pkg/app/config"
	hla_settings "github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleConfigConnectsAPI(
	pCtx context.Context,
	pWrapper config.IWrapper,
	pLogger logger.ILogger,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hla_settings.CAppShortName, pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if pR.Method == http.MethodGet {
			connects := pWrapper.GetConfig().GetConnections()
			result := make([]string, 0, len(connects))
			for _, addr := range connects {
				result = append(result, hla_settings.CAdapterName+"://"+addr)
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
		if u.Scheme != hla_settings.CAdapterName {
			pLogger.PushWarn(logBuilder.WithMessage("scheme_rejected"))
			_ = api.Response(pW, http.StatusAccepted, "rejected: scheme != tcp")
			return
		}

		switch pR.Method {
		case http.MethodPost:
			connects := uniqAppendToSlice(pWrapper.GetConfig().GetConnections(), u.Host)
			if err := pWrapper.GetEditor().UpdateConnections(connects); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("update_connections"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: update connections")
				return
			}

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
