package handler

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/notifier/internal/database"
	"github.com/number571/hidden-lake/internal/services/notifier/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/chars"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/query"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	"github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"
	"github.com/number571/hidden-lake/pkg/api/services/notifier/request"
)

func HandleChanLoadAPI(
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

		var (
			index  = queryParams.Get("index")
			msgUID = queryParams.Get("uid")
		)

		if index != "" && msgUID != "" {
			pLogger.PushWarn(logBuilder.WithMessage("has_many_ids"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: has many ids")
			return
		}

		aliasName := queryParams.Get("friend")
		if aliasName == "" {
			pLogger.PushWarn(logBuilder.WithMessage("alias_name_is_null"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: alias name is null")
			return
		}

		isRemote, err := query.GetBoolValueFromQuery(queryParams, "remote")
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("remote_error"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: remote is unknown")
			return
		}

		var dbMsg dto.IMessage

		switch {
		case index != "":
			dbMsg, err = pDatabase.LoadByEnum(aliasName, parseIndex(index))
		case msgUID != "":
			if !dto.MessageUIDIsValid(msgUID) {
				pLogger.PushWarn(logBuilder.WithMessage("msg_uid"))
				_ = api.Response(pW, http.StatusBadRequest, "failed: message uid")
				return
			}
			dbMsg, err = pDatabase.LoadByUID(aliasName, msgUID)
		default:
			pLogger.PushWarn(logBuilder.WithMessage("id_undefined"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: id undefined")
			return
		}

		if err == nil {
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, dbMsg.ToString(), api.WithCustomContentType(api.CApplicationJSON))
			return
		}

		if !errors.Is(err, database.ErrNotFoundMessage) {
			pLogger.PushErro(logBuilder.WithMessage("load_message"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: load message")
			return
		}

		if isRemote { // nolint:nestif
			if !dto.MessageUIDIsValid(msgUID) {
				pLogger.PushWarn(logBuilder.WithMessage("msg_uid"))
				_ = api.Response(pW, http.StatusBadRequest, "failed: message uid")
				return
			}
			req := request.NewMsgboxLoadRequest(msgUID)
			rsp, err := pHlkClient.FetchRequest(pCtx, aliasName, req)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("no_response"))
				_ = api.Response(pW, http.StatusServiceUnavailable, "failed: no response")
				return
			}
			if rsp.GetCode() != http.StatusOK {
				pLogger.PushWarn(logBuilder.WithMessage("bad_status_code"))
				_ = api.Response(pW, http.StatusTeapot, "failed: status code")
				return
			}
			msgBody := string(rsp.GetBody())
			if chars.HasNotGraphicCharacters(msgBody) {
				pLogger.PushWarn(logBuilder.WithMessage("bad_response"))
				_ = api.Response(pW, http.StatusTeapot, "failed: body response")
				return
			}
			msg := dto.NewMessage(
				dto.NewMessageHead(dto.CIncomingState, msgUID),
				dto.NewMessageBody(msgBody, time.Now()),
			)
			if _, err := pDatabase.Push(aliasName, msg); err != nil {
				pLogger.PushErro(logBuilder.WithMessage("push_message"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
				return
			}
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, msg.ToString(), api.WithCustomContentType(api.CApplicationJSON))
			return
		}

		pLogger.PushWarn(logBuilder.WithMessage("not_found_message"))
		_ = api.Response(pW, http.StatusNotFound, "failed: message not found")
	}
}

func parseIndex(strIndex string) uint64 {
	index, err := strconv.ParseUint(strIndex, 10, 64)
	if err != nil {
		return math.MaxUint64
	}
	return index
}
