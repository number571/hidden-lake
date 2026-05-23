package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/database"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
)

func HandleChatLoadAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pConfig config.IConfig,
	pHlkClient hlk_client.IClient,
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

		index, err := strconv.ParseUint(queryParams.Get("index"), 10, 64)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("parse_index"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: parse index")
			return
		}

		size := pDatabase.Size(aliasName)
		if index >= size {
			pLogger.PushWarn(logBuilder.WithMessage("index_gte_size"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: index >= size")
			return
		}

		dbMsgs, err := pDatabase.Load(aliasName, index, 1)
		if err != nil || len(dbMsgs) != 1 {
			pLogger.PushWarn(logBuilder.WithMessage("load_messages"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: load messages")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, dbMsgs[0])
	}
}
