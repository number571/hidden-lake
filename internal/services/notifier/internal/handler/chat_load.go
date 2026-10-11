package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/notifier/internal/database"
	"github.com/number571/hidden-lake/internal/services/notifier/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleChatLoadAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pConfig config.IConfig,
	pDatabase database.IKVDatabase,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		queryParams := pR.URL.Query()

		aliasName := queryParams.Get("friend")
		if aliasName == "" {
			pLogger.PushWarn(logBuilder.WithMessage("alias_name_is_null"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: alias name is null")
			return
		}

		index, err := strconv.ParseUint(queryParams.Get("index"), 10, 64)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("msg_uid"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: message uid")
			return
		}

		dbMsg, err := pDatabase.LoadByEnum(aliasName, index)
		if err != nil {
			if errors.Is(err, database.ErrNotFoundMessage) {
				pLogger.PushWarn(logBuilder.WithMessage("not_found_message"))
				_ = api.Response(pW, http.StatusNotFound, "failed: message not found")
				return
			}
			pLogger.PushErro(logBuilder.WithMessage("load_message"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: load message")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, dbMsg.ToString(), api.WithCustomContentType(api.CApplicationJSON))
	}
}
