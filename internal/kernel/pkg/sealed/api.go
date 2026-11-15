package sealed

import (
	"context"
	"crypto/subtle"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"

	anonymity "github.com/number571/go-peer/pkg/anonymity/qb"
	"github.com/number571/go-peer/pkg/logger"
	gopeer_network "github.com/number571/go-peer/pkg/network"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/kernel/internal/handler"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/pkg/network"
)

// NewAPIHandler returns the internal API of hlk-sealed. Every request must carry
// Hlk-Auth == pAPIToken and Host == pAddr (address of the bound listener),
// otherwise 401 is returned before the body is read. Body is limited to 64KiB.
func NewAPIHandler(
	pCtx context.Context,
	pCfg IConfig,
	pNode network.IHiddenLakeNode,
	pConnNode gopeer_network.INode,
	pAPIToken string,
	pAddr string,
	pLogger logger.ILogger,
) http.Handler {
	origNode := pNode.GetOriginNode()

	mux := http.NewServeMux()
	mux.HandleFunc(hlk_settings.CHandleConfigSettingsPath, handler.HandleConfigSettingsAPI(pCfg, pLogger, origNode))
	mux.HandleFunc(hlk_settings.CHandleConfigFriendsPath, handleConfigFriendsAPI(pCfg, pLogger, origNode))
	mux.HandleFunc(hlk_settings.CHandleConfigConnectsPath, handleConfigConnectsAPI(pCtx, pCfg, pLogger, pConnNode))
	mux.HandleFunc(hlk_settings.CHandleNetworkOnlinePath, handleNetworkOnlineAPI(pLogger, pConnNode))
	mux.HandleFunc(hlk_settings.CHandleNetworkRequestPath, handleNetworkRequestAPI(pCtx, pCfg, pLogger, pNode))

	return authMiddleware(pAPIToken, pAddr, pLogger, mux)
}

// NewAPIServer wraps the handler into a server with read/write timeouts.
func NewAPIServer(pHandler http.Handler) *http.Server {
	buildSettings := build.GetSettings()
	return &http.Server{
		Handler:           pHandler,
		ReadHeaderTimeout: buildSettings.GetHttpReadTimeout(),
		ReadTimeout:       buildSettings.GetHttpReadTimeout(),
		WriteTimeout:      buildSettings.GetHttpHandleTimeout(),
	}
}

func authMiddleware(pAPIToken, pAddr string, pLogger logger.ILogger, pNext http.Handler) http.Handler {
	return http.HandlerFunc(func(pW http.ResponseWriter, pR *http.Request) {
		tokenOK := pAPIToken != "" && subtle.ConstantTimeCompare([]byte(pR.Header.Get(CHeaderAuth)), []byte(pAPIToken)) == 1
		if !tokenOK || pR.Host != pAddr {
			logBuilder := http_logger.NewLogBuilder(hlk_settings.GetAppShortNameFMT(), pR)
			pLogger.PushWarn(logBuilder.WithMessage("unauthorized"))
			_ = api.Response(pW, http.StatusUnauthorized, "failed: unauthorized")
			return
		}
		pR.Body = http.MaxBytesReader(pW, pR.Body, cMaxBodySize)
		pNext.ServeHTTP(pW, pR)
	})
}

// handleConfigFriendsAPI serializes the kernel handler (read friends -> edit ->
// UpdateFriends is not atomic there) and turns a panic of go-peer on an invalid
// key of the correct length into 400.
func handleConfigFriendsAPI(
	pCfg IConfig,
	pLogger logger.ILogger,
	pNode anonymity.INode,
) http.HandlerFunc {
	mutex := sync.Mutex{}
	friendsF := handler.HandleConfigFriendsAPI(pCfg, pLogger, pNode)
	return func(pW http.ResponseWriter, pR *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		defer func() {
			if recover() != nil {
				logBuilder := http_logger.NewLogBuilder(hlk_settings.GetAppShortNameFMT(), pR)
				pLogger.PushWarn(logBuilder.WithMessage("decode_key"))
				_ = api.Response(pW, http.StatusBadRequest, "failed: load public key")
			}
		}()
		friendsF(pW, pR)
	}
}

// handleNetworkRequestAPI: POST sends the request through the QB queue
// without waiting for a response (PUT of the kernel handler).
func handleNetworkRequestAPI(
	pCtx context.Context,
	pCfg IConfig,
	pLogger logger.ILogger,
	pNode network.IHiddenLakeNode,
) http.HandlerFunc {
	sendF := handler.HandleNetworkRequestAPI(pCtx, pCfg, pLogger, pNode)
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.Method != http.MethodPost {
			logBuilder := http_logger.NewLogBuilder(hlk_settings.GetAppShortNameFMT(), pR)
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}
		// the kernel handler treats POST as FetchRequest (waits for a response up to
		// fetch timeout) and PUT as SendRequest (QB queue only); the client expects
		// POST = send without waiting, so the kernel handler gets a PUT copy.
		cp := pR.Clone(pR.Context())
		cp.Method = http.MethodPut
		sendF(pW, cp)
	}
}

func handleConfigConnectsAPI(
	pCtx context.Context,
	pCfg IConfig,
	pLogger logger.ILogger,
	pConnNode gopeer_network.INode,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlk_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if pR.Method == http.MethodGet {
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, pCfg.GetEndpoints())
			return
		}

		connectBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: read connect bytes")
			return
		}

		addr := strings.TrimSpace(string(connectBytes))
		if !isValidEndpoint(addr) {
			pLogger.PushWarn(logBuilder.WithMessage("read_connect"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: invalid connect")
			return
		}

		if pR.Method == http.MethodPost {
			pCfg.AddEndpoint(addr)
			_ = pConnNode.AddConnection(pCtx, addr) // connection may be refused (closed)

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: update connections")
			return
		}

		pCfg.DelEndpoint(addr)
		_ = pConnNode.DelConnection(addr) // connection may be already closed

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, "success: delete connection")
	}
}

func handleNetworkOnlineAPI(pLogger logger.ILogger, pConnNode gopeer_network.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlk_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		connects := pConnNode.GetConnections()
		online := make([]string, 0, len(connects))
		for addr := range connects {
			online = append(online, addr)
		}
		sort.Strings(online)

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, online)
	}
}
