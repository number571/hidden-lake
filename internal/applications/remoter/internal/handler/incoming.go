package handler

import (
	"context"
	"io"
	"net/http"
	"os/exec"
	"strings"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/remoter/pkg/app/config"
	hlr_settings "github.com/number571/hidden-lake/internal/applications/remoter/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/chars"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleIncomingExecHTTP(pCtx context.Context, pConfig config.IConfig, pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hls_settings.CHeaderResponseMode, hls_settings.CHeaderResponseModeON)

		logBuilder := http_logger.NewLogBuilder(hlr_settings.GServiceName.Short(), pR)

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		sett := pConfig.GetSettings()
		if pR.Header.Get(hlr_settings.CHeaderPassword) != sett.GetPassword() {
			pLogger.PushWarn(logBuilder.WithMessage("forbidden"))
			_ = api.Response(pW, http.StatusForbidden, "failed: request forbidden")
			return
		}

		cmdBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: response message")
			return
		}

		cmdStr := string(cmdBytes)
		if chars.HasNotGraphicCharacters(cmdStr) {
			pLogger.PushWarn(logBuilder.WithMessage("has_not_graphic_chars"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: has not graphic characters")
			return
		}

		ctx, cancel := context.WithTimeout(pCtx, sett.GetExecTimeout())
		defer cancel()

		cmdSplited := strings.Split(cmdStr, hlr_settings.CExecSeparator)
		out, err := exec.CommandContext(ctx, cmdSplited[0], cmdSplited[1:]...).Output() // nolint: gosec
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("execute_command"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: "+err.Error())
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(cmdStr))
		_ = api.Response(pW, http.StatusOK, out)
	}
}
