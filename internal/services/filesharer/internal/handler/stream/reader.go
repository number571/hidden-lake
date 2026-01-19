package stream

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/number571/go-peer/pkg/encoding"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/handler/incoming/limiters"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/client/fileinfo"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

var (
	_ io.Reader = &sStream{}
)

type sStream struct {
	fContext   context.Context
	fRetryNum  uint64
	fTempFile  string
	fHlkClient hlk_client.IClient
	fAliasName string
	fFileInfo  fileinfo.IFileInfo
	fBuffer    []byte
	fPosition  uint64
	fTempBytes []byte
	fHasher    hash.Hash
	fChunkSize uint64
}

func BuildStreamReader(
	pCtx context.Context,
	pRetryNum uint64,
	pInputPath string,
	pAliasName string,
	pHlkClient hlk_client.IClient,
	pFileInfo fileinfo.IFileInfo,
) (io.Reader, string, error) {
	chunkSize, err := limiters.GetLimitOnLoadResponseSize(pCtx, pHlkClient)
	if err != nil {
		return nil, "", errors.Join(ErrGetMessageLimit, err)
	}

	tempName := fmt.Sprintf(hls_filesharer_settings.CPathTMP, pFileInfo.GetHash()[:8])
	tempFile := filepath.Join(pInputPath, tempName)

	// TODO: fix
	tempBytes, err := readTempFile(tempFile, chunkSize)
	if err != nil {
		return nil, "", errors.Join(ErrReadTempFile, err)
	}

	return &sStream{
		fContext:   pCtx,
		fRetryNum:  pRetryNum,
		fTempFile:  tempFile,
		fTempBytes: tempBytes,
		fHlkClient: pHlkClient,
		fAliasName: pAliasName,
		fHasher:    sha512.New384(),
		fChunkSize: chunkSize,
		fFileInfo:  pFileInfo,
	}, tempFile, nil
}

func readTempFile(pTempFile string, pChunkSize uint64) ([]byte, error) {
	if _, err := os.Stat(pTempFile); errors.Is(err, os.ErrNotExist) {
		if _, err := os.Create(pTempFile); err != nil { // nolint: gosec
			return nil, err
		}
	}
	tempBytes, err := os.ReadFile(pTempFile) // nolint: gosec
	if err != nil {
		_ = os.Remove(pTempFile)
		return nil, err
	}
	if uint64(len(tempBytes))%pChunkSize != 0 {
		_ = os.Remove(pTempFile)
		return nil, errors.New("len(tempBytes) %% chunkSize != 0") // nolint: err113
	}
	return tempBytes, nil
}

func (p *sStream) Read(b []byte) (int, error) {
	select {
	case <-p.fContext.Done():
		return 0, io.ErrClosedPipe
	default:
	}

	if p.fTempBytes != nil {
		p.fBuffer = p.fTempBytes
		p.fTempBytes = nil
	}

	if len(p.fBuffer) == 0 {
		chunk, err := p.loadFileChunk()
		if err != nil {
			return 0, errors.Join(ErrLoadFileChunk, err)
		}
		if err := p.appendToTempFile(chunk); err != nil {
			return 0, errors.Join(ErrAppendToTempFile, err)
		}
		p.fBuffer = chunk
	}

	n := copy(b, p.fBuffer)
	p.fBuffer = p.fBuffer[n:]
	p.fPosition += uint64(n) //nolint:gosec

	if _, err := p.fHasher.Write(b[:n]); err != nil {
		return 0, errors.Join(ErrHashWriteChunk, err)
	}

	if p.fPosition < p.fFileInfo.GetSize() {
		return n, nil
	}

	if err := os.Remove(p.fTempFile); err != nil {
		return 0, errors.Join(ErrDeleteTempFile, err)
	}

	hashSum := encoding.HexEncode(p.fHasher.Sum(nil))
	if hashSum != p.fFileInfo.GetHash() {
		return 0, ErrInvalidHash
	}

	return n, io.EOF
}

func (p *sStream) loadFileChunk() ([]byte, error) {
	var lastErr error
	for i := uint64(0); i <= p.fRetryNum; i++ {
		req := newLoadChunkRequest(p.fFileInfo.GetName(), p.fPosition/p.fChunkSize)
		resp, err := p.fHlkClient.FetchRequest(p.fContext, p.fAliasName, req)
		if err != nil {
			lastErr = err
			continue
		}
		if resp.GetCode() != http.StatusOK {
			lastErr = err
			continue
		}
		return resp.GetBody(), nil
	}
	return nil, errors.Join(ErrRetryFailed, lastErr)
}

func (p *sStream) appendToTempFile(chunk []byte) error {
	f, err := os.OpenFile(p.fTempFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	if _, err := f.Write(chunk); err != nil {
		return err
	}
	return nil
}

func newLoadChunkRequest(pFileName string, pChunk uint64) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_filesharer_settings.CAppShortName).
		WithPath(fmt.Sprintf(
			"%s?name=%s&chunk=%d",
			hls_filesharer_settings.CLoadPath,
			url.QueryEscape(pFileName),
			pChunk,
		)).
		Build()
}
