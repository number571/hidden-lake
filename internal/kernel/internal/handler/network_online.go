package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sort"
	"sync"

	"github.com/number571/go-peer/pkg/logger"
	pkg_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/pkg/api/adapters/http/client"
)

func HandleNetworkOnlineAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pEPClients []client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		lenEPClients := len(pEPClients)
		errorList := make([]error, lenEPClients)

		switch pR.Method {
		case http.MethodGet:
			inOnlines := make([][]string, lenEPClients)

			wg := &sync.WaitGroup{}
			wg.Add(lenEPClients)
			for i, c := range pEPClients {
				go func(i int, c client.IClient) {
					defer wg.Done()
					inOnlines[i], errorList[i] = c.GetOnlines(pCtx)
				}(i, c)
			}
			wg.Wait()

			if err := errors.Join(errorList...); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("get_connections"))
				_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: get online connections")
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, mergeAndSortOnlines(inOnlines))
		case http.MethodDelete:
			connectBytes, err := io.ReadAll(pR.Body)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
				_ = api.Response(pW, http.StatusConflict, "failed: read connect bytes")
				return
			}

			wg := &sync.WaitGroup{}
			wg.Add(lenEPClients)
			for i, c := range pEPClients {
				go func(i int, c client.IClient) {
					defer wg.Done()
					errorList[i] = c.DelOnline(pCtx, string(connectBytes))
				}(i, c)
			}
			wg.Wait()

			if err := errors.Join(errorList...); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("del_connection"))
				_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: delete online connections")
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: delete online connection")
		}
	}
}

func mergeAndSortOnlines(inOnlines [][]string) []string {
	size := 0
	for _, v := range inOnlines {
		size += len(v)
	}
	mergeInOnlines := make([]string, 0, size)
	for _, s := range inOnlines {
		mergeInOnlines = append(mergeInOnlines, s...)
	}
	sort.SliceStable(mergeInOnlines, func(i, j int) bool {
		return mergeInOnlines[i] < mergeInOnlines[j]
	})
	return mergeInOnlines
}
