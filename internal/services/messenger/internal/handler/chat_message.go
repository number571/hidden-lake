package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/database"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/app/config"
	hls_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/chars"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/pubkey"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	message "github.com/number571/hidden-lake/pkg/api/services/messenger/client/dto"
	"github.com/number571/hidden-lake/pkg/api/services/messenger/request"
)

func HandleChatMessageAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pConfig config.IConfig,
	pHlkClient hlk_client.IClient,
	pDatabase database.IKVDatabase,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if pR.Method == http.MethodGet {
			size, err := getLimitOnPushRequestSize(pCtx, pHlkClient)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("get_limit"))
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

		strBody := string(body)
		if len(strBody) == 0 || chars.HasNotGraphicCharacters(strBody) {
			pLogger.PushWarn(logBuilder.WithMessage("has_not_graphic_chars"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: has not graphic characters")
			return
		}

		aliasName := pR.URL.Query().Get("friend")
		fPubKey, err := pubkey.GetFriendPubKeyByAliasName(pCtx, pHlkClient, aliasName)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("load_pubkey"))
			_ = api.Response(pW, http.StatusForbidden, "failed: load public key")
			return
		}

		myPubKey, err := pHlkClient.GetPubKey(pCtx)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_public_key"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: get public key from service")
			return
		}

		req := request.NewPushRequest(strBody)
		if err := pHlkClient.SendRequest(pCtx, aliasName, req); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("send_request"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: send request")
			return
		}

		msg := message.NewMessage(false, string(body), time.Now())
		if err := pDatabase.Push(database.NewRelation(myPubKey, fPubKey), msg); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("push_message"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
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

	reqSize := uint64(len(request.NewPushRequest("").ToBytes()))
	pldLimit := sett.GetPayloadSizeBytes()
	if reqSize >= pldLimit {
		return 0, errors.New("request size >= payload limit") // nolint: err113
	}

	return pldLimit - reqSize, nil
}
