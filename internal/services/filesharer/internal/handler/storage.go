package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/logger"
	hls_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/stream"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	hls_filesharer_client "github.com/number571/hidden-lake/internal/services/filesharer/pkg/client"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/webui"
)

type sStorage struct {
	*sTemplate
	FPage      uint64
	FAliasName string
	FFilesList []hls_filesharer_settings.SFileInfo
}

func StoragePage(
	pCtx context.Context,
	pLogger logger.ILogger,
	pCfg config.IConfig,
	pPathTo string,
	pHlsClient hls_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_filesharer_settings.GetShortAppName(), pR)

		if pR.URL.Path != "/friends/storage" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		query := pR.URL.Query()
		aliasName := query.Get("alias_name")
		if aliasName == "" {
			ErrorPage(pLogger, pCfg, "alias_name_error", "incorrect alias name")(pW, pR)
			return
		}

		if fileName := query.Get("file_name"); fileName != "" {
			downloadFile(pCtx, pLogger, pCfg, pPathTo, pW, pR, pHlsClient)
			return
		}

		page, err := strconv.Atoi(query.Get("page"))
		if err != nil {
			page = 0
		}

		hlfClient := hls_filesharer_client.NewClient(
			hls_filesharer_client.NewBuilder(),
			hls_filesharer_client.NewRequester(pHlsClient),
		)

		filesList, err := hlfClient.GetListFiles(pCtx, aliasName, uint64(page)) //nolint:gosec
		if err != nil {
			ErrorPage(pLogger, pCfg, "get_files_list", "failed get list of files")(pW, pR)
			return
		}

		result := sStorage{
			sTemplate:  getTemplate(pCfg),
			FPage:      uint64(page), //nolint:gosec
			FAliasName: aliasName,
			FFilesList: filesList,
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = webui.MustParseTemplate("index.html", "filesharer/storage.html").Execute(pW, result)
	}
}

func downloadFile(
	pCtx context.Context,
	pLogger logger.ILogger,
	pCfg config.IConfig,
	pPathTo string,
	pW http.ResponseWriter,
	pR *http.Request,
	pHlsClient hls_client.IClient,
) {
	query := pR.URL.Query()

	var (
		aliasName = query.Get("alias_name")
		fileName  = query.Get("file_name")
		fileHash  = query.Get("file_hash")
	)

	fileSize, err := strconv.ParseUint(query.Get("file_size"), 10, 64)
	if err != nil {
		ErrorPage(pLogger, pCfg, "file_size_error", "incorrect file size")(pW, pR)
		return
	}

	myPubKey, err := pHlsClient.GetPubKey(pCtx)
	if err != nil {
		ErrorPage(pLogger, pCfg, "get_pub_key", "get public key")(pW, pR)
		return
	}

	friends, err := pHlsClient.GetFriends(pCtx)
	if err != nil {
		ErrorPage(pLogger, pCfg, "get_friends", "get friends")(pW, pR)
		return
	}

	fPubKey, ok := friends[aliasName]
	if !ok {
		ErrorPage(pLogger, pCfg, "get_friend_pub_key", "get friend pub key")(pW, pR)
		return
	}

	tempFile := filepath.Join(pPathTo, tempFilename(myPubKey, fPubKey, fileHash))
	if _, err := os.Stat(tempFile); errors.Is(err, os.ErrNotExist) {
		if _, err := os.Create(tempFile); err != nil { // nolint: gosec
			ErrorPage(pLogger, pCfg, "create_temp_file", "create temp file")(pW, pR)
			return
		}
	}

	pW.Header().Set("Content-Type", "application/octet-stream")
	pW.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(fileName))
	pW.Header().Set("Content-Length", strconv.FormatUint(fileSize, 10))

	chCtx, cancel := context.WithCancel(pCtx)
	defer cancel()

	stream, err := stream.BuildStream(
		chCtx,
		pCfg.GetSettings().GetRetryNum(),
		tempFile,
		pHlsClient,
		aliasName,
		stream.NewFileInfo(fileName, fileHash, fileSize),
	)
	if err != nil {
		ErrorPage(pLogger, pCfg, "build_stream", "build stream")(pW, pR)
		return
	}

	go func() {
		select {
		case <-chCtx.Done():
		case <-pR.Context().Done():
			cancel()
		}
	}()

	http.ServeContent(pW, pR, fileName, time.Now(), stream)
}

func tempFilename(pMyPubKey, pFPubKey asymmetric.IPubKey, pFileHash string) string {
	hash := hashing.NewHasher(bytes.Join([][]byte{
		pMyPubKey.ToBytes(),
		pFPubKey.ToBytes(),
		[]byte(pFileHash),
	}, []byte{})).ToString()
	return fmt.Sprintf(hls_filesharer_settings.CPathTMP, hash[:8])
}
