package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/notifier/internal/database"
	"github.com/number571/hidden-lake/internal/services/notifier/internal/handler/process"
	"github.com/number571/hidden-lake/internal/services/notifier/pkg/app/config"
	hls_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/chars"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	"github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"
	"github.com/number571/hidden-lake/pkg/api/services/notifier/request"
)

func HandleChatPushAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pConfig config.IConfig,
	pHlkClient hlk_client.IClient,
	pDatabase database.IKVDatabase,
	pProcessManager process.INotifyProcessManager,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost && pR.Method != http.MethodPut {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if pR.Method == http.MethodGet {
			size, err := getLimitOnPushRequestSize(pCtx, pHlkClient)
			if err != nil {
				pLogger.PushErro(logBuilder.WithMessage("get_limit"))
				_ = api.Response(pW, http.StatusBadGateway, "failed: get limit")
				return
			}
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, strconv.FormatUint(size, 10))
			return
		}

		body, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: decode request")
			return
		}

		msgBody := string(body)
		if len(msgBody) == 0 || chars.HasNotGraphicCharacters(msgBody) {
			pLogger.PushWarn(logBuilder.WithMessage("has_not_graphic_chars"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: has not graphic characters")
			return
		}

		queryParams := pR.URL.Query()

		aliasName := queryParams.Get("friend")
		if aliasName == "" {
			pLogger.PushWarn(logBuilder.WithMessage("alias_name_is_null"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: alias name is null")
			return
		}

		msgUID := queryParams.Get("uid")
		if pR.Method == http.MethodPut {
			msgUID = dto.NewMessageUID()
		}

		if !dto.MessageUIDIsValid(msgUID) {
			pLogger.PushWarn(logBuilder.WithMessage("msg_uid"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: message uid")
			return
		}

		msg, err := pDatabase.LoadByUID(aliasName, msgUID)
		if err != nil { // nolint:nestif
			if !errors.Is(err, database.ErrNotFoundMessage) {
				pLogger.PushErro(logBuilder.WithMessage("load_message"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: load message from database")
				return
			}
			// create new message
			msg, err = dto.NewMessageWithParams(
				dto.CSendingState,
				msgBody,
				msgUID,
				time.Now(),
			)
			if err != nil {
				pLogger.PushErro(logBuilder.WithMessage("create_message"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: create message")
				return
			}
			if pR.Method == http.MethodPost {
				if _, err := pDatabase.Push(aliasName, msg); err != nil {
					pLogger.PushErro(logBuilder.WithMessage("push_message"))
					_ = api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
					return
				}
			}
		}

		// state(msg) = incoming OR received
		if msg.GetState() != dto.CSendingState {
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusAccepted, msg.GetTimestamp())
			return
		}

		// fetch process on UID already exist
		if ok := pProcessManager.TryLock(msgUID); !ok {
			pLogger.PushInfo(logBuilder.WithMessage("process_exist"))
			_ = api.Response(pW, http.StatusAlreadyReported, msg.GetTimestamp())
			return
		}
		defer func() { _ = pProcessManager.Unlock(msgUID) }()

		if pR.Method == http.MethodPut {
			req := request.NewPushRequest(msg)
			err := pHlkClient.SendRequest(pCtx, aliasName, req)
			if err != nil {
				pLogger.PushErro(logBuilder.WithMessage("send_request"))
				_ = api.Response(pW, http.StatusBadGateway, "failed: send request")
				return
			}
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, msg.GetTimestamp())
			return
		}

		// try get response from node (delivery check)
		req := request.NewPushRequestWithUID(msg)
		rsp, err := pHlkClient.FetchRequest(pCtx, aliasName, req)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("no_response"))
			_ = api.Response(pW, http.StatusCreated, msg.GetTimestamp())
			return
		}

		switch rsp.GetCode() {
		case http.StatusOK, http.StatusAccepted:
			// ok
		default:
			pLogger.PushWarn(logBuilder.WithMessage("bad_status_code"))
			_ = api.Response(pW, http.StatusTeapot, "failed: status code")
			return
		}

		if err := pDatabase.Recv(aliasName, msgUID); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("recv_response"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: recv database")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, msg.GetTimestamp())
	}
}

func getLimitOnPushRequestSize(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	sett, err := pHlkClient.GetSettings(pCtx)
	if err != nil {
		return 0, err
	}

	msg := dto.NewMessage(dto.CIncomingState, "")
	reqSize := uint64(len(request.NewPushRequest(msg).ToBytes()))
	pldLimit := sett.GetPayloadSizeBytes()
	if reqSize >= pldLimit {
		return 0, errors.New("request size >= payload limit") // nolint: err113
	}

	return pldLimit - reqSize, nil
}
