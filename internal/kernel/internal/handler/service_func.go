package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/kernel/pkg/app/config"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/pkg/handler"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/logger"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	internal_anon_logger "github.com/number571/hidden-lake/internal/utils/logger/anon"
)

func HandleServiceFunc(pCfg config.IConfig, pLogger logger.ILogger) handler.IHandlerF {
	return func(
		pCtx context.Context,
		pSender asymmetric.IPubKey,
		pRequest request.IRequest,
	) (response.IResponse, error) {
		logBuilder := anon_logger.NewLogBuilder(hlk_settings.GetAppShortNameFMT())

		// get service's address by hostname
		service, ok := pCfg.GetService(pRequest.GetHost())
		if !ok {
			pLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnUndefinedService))
			return nil, ErrUndefinedService
		}

		// generate new request to serivce
		pushReq, err := http.NewRequestWithContext(
			pCtx,
			pRequest.GetMethod(),
			fmt.Sprintf("http://%s%s", service, pRequest.GetPath()),
			bytes.NewReader(pRequest.GetBody()),
		)
		if err != nil {
			pLogger.PushErro(logBuilder.WithType(internal_anon_logger.CLogErroProxyRequestType))
			return nil, errors.Join(ErrBuildRequest, err)
		}

		// append headers from request & set service headers
		for key, val := range pRequest.GetHead() {
			pushReq.Header.Set(key, val)
		}
		setAliasNameToHeaders(pCfg, pushReq, pSender)

		// send request and receive response from service
		httpClient := &http.Client{Timeout: build.GetSettings().GetHttpHandleTimeout()}
		resp, err := httpClient.Do(pushReq)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnRequestToService))
			return nil, errors.Join(ErrBadRequest, err)
		}
		defer func() { _ = resp.Body.Close() }()

		// get response mode: on/off
		respMode := resp.Header.Get(hlk_settings.CHeaderResponseMode)
		switch respMode {
		case "", hlk_settings.CHeaderResponseModeON:
			// send response to the client
			pLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogInfoResponseFromService))
			return response.NewResponseBuilder().
					WithCode(resp.StatusCode).
					WithHead(getResponseHead(resp)).
					WithBody(getResponseBody(resp)).
					Build(),
				nil
		case hlk_settings.CHeaderResponseModeOFF:
			// response is not required by the client side
			pLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogBaseResponseModeFromService))
			return nil, nil
		default:
			// unknown response mode
			pLogger.PushErro(logBuilder.WithType(internal_anon_logger.CLogBaseResponseModeFromService))
			return nil, ErrInvalidResponseMode
		}
	}
}

func setAliasNameToHeaders(pCfg config.IConfig, pReq *http.Request, pSender asymmetric.IPubKey) {
	for aliasName, pubKey := range pCfg.GetFriends() {
		if pubKey.ToString() == pSender.ToString() {
			pReq.Header.Set(hlk_settings.CHeaderSenderName, aliasName)
			break
		}
	}
}

func getResponseHead(pResp *http.Response) map[string]string {
	headers := make(map[string]string, len(pResp.Header))
	for k := range pResp.Header {
		if _, ok := gIgnoreHeaders[k]; ok {
			continue
		}
		headers[k] = pResp.Header.Get(k)
	}
	return headers
}

func getResponseBody(pResp *http.Response) []byte {
	data, err := io.ReadAll(pResp.Body)
	if err != nil {
		return nil
	}
	return data
}
