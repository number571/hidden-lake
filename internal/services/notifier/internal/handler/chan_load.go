package handler

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/notifier/internal/database"
	"github.com/number571/hidden-lake/internal/services/notifier/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"
)

func HandleChanLoadAPI(
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

		var (
			index = queryParams.Get("index")
			uid   = queryParams.Get("uid")
		)

		aliasName := queryParams.Get("friend")
		if aliasName == "" {
			pLogger.PushWarn(logBuilder.WithMessage("alias_name_is_null"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: alias name is null")
			return
		}

		var (
			dbMsg dto.IMessage
			err   error
		)

		switch {
		case index != "":
			dbMsg, err = pDatabase.LoadByEnum(aliasName, parseIndex(index))
		case uid != "":
			dbMsg, err = pDatabase.LoadByUID(aliasName, uid)
		default:
			pLogger.PushWarn(logBuilder.WithMessage("id_undefined"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: id undefined")
			return
		}

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

func parseIndex(strIndex string) uint64 {
	index, err := strconv.ParseUint(strIndex, 10, 64)
	if err != nil {
		return math.MaxUint64
	}
	return index
}
