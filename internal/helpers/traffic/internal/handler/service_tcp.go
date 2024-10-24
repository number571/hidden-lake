package handler

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
	"github.com/number571/go-peer/pkg/network/conn"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/utils"
	"github.com/number571/hidden-lake/internal/utils/api"

	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/config"
	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/storage"
	hlt_settings "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/settings"
)

func HandleServiceTCP(
	pCfg config.IConfig,
	pStorage storage.IMessageStorage,
	pLogger logger.ILogger,
) network.IHandlerF {
	httpClient := &http.Client{Timeout: time.Minute}

	return func(pCtx context.Context, pNode network.INode, pConn conn.IConn, pNetMsg net_message.IMessage) error {
		logBuilder := anon_logger.NewLogBuilder(hlt_settings.CServiceName)

		// enrich logger
		logBuilder.
			WithConn(pConn).
			WithHash(pNetMsg.GetHash()).
			WithProof(pNetMsg.GetProof()).
			WithSize(len(pNetMsg.ToBytes()))

		if _, err := message.LoadMessage(pCfg.GetSettings().GetMessageSizeBytes(), pNetMsg.GetPayload().GetBody()); err != nil {
			pLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnMessageNull))
			return utils.MergeErrors(ErrLoadMessage, err)
		}

		// check message from in storage queue
		if err := pStorage.Push(pNetMsg); err != nil {
			if errors.Is(err, storage.ErrMessageIsExist) || errors.Is(err, storage.ErrHashAlreadyExist) {
				pLogger.PushInfo(logBuilder.WithType(anon_logger.CLogInfoExist))
				return nil
			}
			pLogger.PushErro(logBuilder.WithType(anon_logger.CLogErroDatabaseSet))
			return utils.MergeErrors(ErrPushMessageDB, err)
		}

		// some of connections may be closed
		// need pass return error if exist
		logBroadcastMessage(
			pLogger,
			logBuilder,
			pNode.BroadcastMessage(pCtx, pNetMsg),
		)

		consumers := pCfg.GetConsumers()

		wg := sync.WaitGroup{}
		wg.Add(len(consumers))

		for _, cHost := range consumers {
			go func(cHost string) {
				defer wg.Done()
				_, err := api.Request(
					pCtx,
					httpClient,
					http.MethodPost,
					"http://"+cHost,
					pNetMsg.ToString(),
				)
				if err != nil {
					pLogger.PushWarn(logBuilder.WithType(anon_logger.CLogBaseGetResponse))
					return
				}
				pLogger.PushInfo(logBuilder.WithType(anon_logger.CLogBaseGetResponse))
			}(cHost)
		}

		wg.Wait()
		return nil
	}
}

func logBroadcastMessage(
	pLogger logger.ILogger,
	pLogBuilder anon_logger.ILogBuilder,
	pErr error,
) {
	if pErr != nil {
		pLogger.PushWarn(pLogBuilder.WithType(anon_logger.CLogBaseBroadcast))
		return
	}
	pLogger.PushInfo(pLogBuilder.WithType(anon_logger.CLogBaseBroadcast))
}
