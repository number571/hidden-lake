package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/utils"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/query"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	fileinfo "github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
	"github.com/number571/hidden-lake/pkg/api/services/filesharer/request"
)

func HandleRemoteFileInfoAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pHlkClient hlk_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)

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

		fileName := queryParams.Get("name")
		if utils.FileNameIsInvalid(fileName) {
			pLogger.PushWarn(logBuilder.WithMessage("got_invalid_name"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: got invalid name")
			return
		}

		isPersonal, err := query.GetBoolValueFromQuery(queryParams, "personal")
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("parse_personal"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: parse personal")
			return
		}

		req := request.NewInfoRequest(fileName, isPersonal)
		resp, err := pHlkClient.FetchRequest(pCtx, aliasName, req)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("fetch_request"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: fetch request")
			return
		}

		if code := resp.GetCode(); code != http.StatusOK {
			pLogger.PushWarn(logBuilder.WithMessage("status_error"))
			_ = api.Response(pW, http.StatusTeapot, fmt.Sprintf("failed: status %d", code))
			return
		}

		info, err := fileinfo.LoadFileInfo(resp.GetBody())
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("decode_response"))
			_ = api.Response(pW, http.StatusTeapot, "failed: decode response")
			return
		}

		if info.GetName() != fileName {
			pLogger.PushWarn(logBuilder.WithMessage("invalid_response"))
			_ = api.Response(pW, http.StatusTeapot, "failed: invalid response")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, info.ToString(), api.WithCustomContentType(api.CApplicationJSON))
	}
}
