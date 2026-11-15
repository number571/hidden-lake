package sealed

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/qb/logger"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/kernel/pkg/app/config"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	internal_anon_logger "github.com/number571/hidden-lake/internal/utils/logger/anon"
	"github.com/number571/hidden-lake/pkg/network/handler"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
)

const (
	CHeaderAuth = "Hlk-Auth"

	cMaxBodySize = 64 << 10 // 64KiB: internal API request body, drained service response
)

// NewServiceHandler forwards network requests POST /push of the hls-messenger
// service only (to http://service_addr/push, redirects are not followed).
// Peer headers are never copied: the request carries only Hlk-Sender-Name
// (alias of the sender from the config) and Hlk-Auth (push token).
// The service response is never returned to the network (mode off).
func NewServiceHandler(pCfg config.IConfig, pPushToken string, pLogger logger.ILogger) handler.IHandlerF {
	// no redirects and no proxy: Hlk-Auth must never leave service_addr
	httpClient := &http.Client{
		Timeout:   build.GetSettings().GetHttpHandleTimeout(),
		Transport: &http.Transport{Proxy: nil},
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return func(
		pCtx context.Context,
		pSender layer2.IParticipantKey,
		pRequest request.IRequest,
	) (response.IResponse, error) {
		logBuilder := anon_logger.NewLogBuilder(hlk_settings.GetAppShortNameFMT())

		if pRequest.GetHost() != hls_messenger_settings.CAppShortName {
			pLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnUndefinedService))
			return nil, ErrUndefinedService
		}
		service, ok := pCfg.GetService(pRequest.GetHost())
		if !ok {
			pLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnUndefinedService))
			return nil, ErrUndefinedService
		}

		// Forward only an exact POST to the fixed /push path; the request URL is
		// built from constants and the configured service address, never from
		// anything the peer sends. Any other method or path is rejected.
		if pRequest.GetMethod() != http.MethodPost || pRequest.GetPath() != hls_messenger_settings.CPushPath {
			pLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnInvalidRequestMethod))
			return nil, ErrInvalidRequest
		}

		aliasName, ok := getAliasName(pCfg, pSender)
		if !ok {
			pLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnRequestToService))
			return nil, ErrUnknownSender
		}

		pushReq, err := http.NewRequestWithContext(
			pCtx,
			http.MethodPost,
			"http://"+service+hls_messenger_settings.CPushPath,
			bytes.NewReader(pRequest.GetBody()),
		)
		if err != nil {
			pLogger.PushErro(logBuilder.WithType(internal_anon_logger.CLogErroProxyRequestType))
			return nil, errors.Join(ErrBuildRequest, err)
		}
		pushReq.Header.Set(hlk_settings.CHeaderSenderName, aliasName)
		pushReq.Header.Set(CHeaderAuth, pPushToken)

		resp, err := httpClient.Do(pushReq)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnRequestToService))
			return nil, errors.Join(ErrRequestToService, err)
		}
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, cMaxBodySize))
		_ = resp.Body.Close()

		pLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogBaseResponseModeFromService))
		return nil, nil
	}
}

func getAliasName(pCfg config.IConfig, pSender layer2.IParticipantKey) (string, bool) {
	sender := pSender.ToString()
	for aliasName, key := range pCfg.GetFriends() {
		if key.ToString() == sender {
			return aliasName, true
		}
	}
	return "", false
}
