package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/publisher/internal/database"
	"github.com/number571/hidden-lake/internal/services/publisher/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/services/publisher/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/query"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	"github.com/number571/hidden-lake/pkg/api/services/publisher/client/dto"
	"github.com/number571/hidden-lake/pkg/api/services/publisher/request"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func HandleMsgboxLoadAPI(
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
			msgNID    = queryParams.Get("nid")
			msgUID    = queryParams.Get("uid")
			aliasName = queryParams.Get("friend")
		)

		if msgNID != "" && msgUID != "" {
			pLogger.PushWarn(logBuilder.WithMessage("has_many_ids"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: has many ids")
			return
		}
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

		var (
			dbMsg dto.IMessage
			index uint64
		)

		switch {
		case msgNID != "":
			index, err = strconv.ParseUint(msgNID, 10, 64)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("msg_index"))
				_ = api.Response(pW, http.StatusBadRequest, "failed: message index")
				return
			}
			dbMsg, err = pDatabase.LoadByEnum(aliasName, index)
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

		if !isRemote {
			pLogger.PushWarn(logBuilder.WithMessage("not_found_message"))
			_ = api.Response(pW, http.StatusNotFound, "failed: message not found")
			return
		}

		var req hlk_request.IRequest
		if msgUID != "" {
			req = request.NewMsgboxLoadRequestByUID(msgUID)
		} else {
			req = request.NewMsgboxLoadRequestByEnum(index)
		}

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

		msgBody := rsp.GetBody()
		if len(msgBody) == 0 {
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
	}
}
