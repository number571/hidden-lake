package stream

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/handler/incoming/limiters"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	fileinfo "github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
	"github.com/number571/hidden-lake/pkg/api/services/filesharer/request"
)

var (
	_ io.Reader = &sStream{}
)

type sStream struct {
	fContext      context.Context
	fRetryNum     uint64
	fTempEOF      bool
	fTempPosition uint64
	fTempFile     string
	fHlkClient    hlk_client.IClient
	fAliasName    string
	fFileInfo     fileinfo.IFileInfo
	fBuffer       []byte
	fPosition     uint64
	fHasher       hash.Hash
	fChunkSize    uint64
	fPersonal     bool
}

func BuildStreamReader(
	pCtx context.Context,
	pRetryNum uint64,
	pInputPath string,
	pAliasName string,
	pHlkClient hlk_client.IClient,
	pFileInfo fileinfo.IFileInfo,
	pPersonal bool,
) (io.Reader, error) {
	chunkSize, err := limiters.GetLimitOnLoadResponseSize(pCtx, pHlkClient)
	if err != nil {
		return nil, errors.Join(ErrGetMessageLimit, err)
	}

	tempName := fmt.Sprintf(hls_filesharer_settings.CPathTMP, pFileInfo.GetHash()[:8])
	tempFile := filepath.Join(pInputPath, tempName)

	if err := createTempFile(tempFile); err != nil {
		return nil, errors.Join(ErrReadTempFile, err)
	}

	return &sStream{
		fContext:   pCtx,
		fRetryNum:  pRetryNum,
		fTempFile:  tempFile,
		fHlkClient: pHlkClient,
		fAliasName: pAliasName,
		fHasher:    sha512.New384(),
		fChunkSize: chunkSize,
		fFileInfo:  pFileInfo,
		fPersonal:  pPersonal,
	}, nil
}

func (p *sStream) Read(b []byte) (int, error) {
	select {
	case <-p.fContext.Done():
		return 0, io.ErrClosedPipe
	default:
	}

	if len(p.fBuffer) == 0 && !p.fTempEOF {
		switch chunk, err := p.loadFileChunkFromTemp(); {
		case err == nil:
			p.fBuffer = chunk
		case errors.Is(err, io.EOF):
			p.fTempEOF = true
		default:
			_ = os.Remove(p.fTempFile)
			return 0, errors.Join(ErrLoadFileChunk, err)
		}
	}

	if len(p.fBuffer) == 0 {
		chunk, err := p.loadFileChunk()
		if err != nil {
			return 0, errors.Join(ErrLoadFileChunk, err)
		}
		if err := p.appendChunkToTempFile(chunk); err != nil {
			_ = os.Remove(p.fTempFile)
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

	fileSize := p.fFileInfo.GetSize()
	switch {
	case p.fPosition < fileSize:
		return n, nil
	case p.fPosition > fileSize:
		return 0, ErrInvalidSize
	default:
	}

	hashSum := encoding.HexEncode(p.fHasher.Sum(nil))
	if hashSum != p.fFileInfo.GetHash() {
		_ = os.Remove(p.fTempFile)
		return 0, ErrInvalidHash
	}

	return n, io.EOF
}

func createTempFile(pTempFile string) error {
	stat, err := os.Stat(pTempFile)
	if errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(pTempFile) // nolint: gosec
		return err
	}
	if err != nil {
		return err
	}
	expired := time.Since(stat.ModTime()) > (30 * 24 * time.Hour)
	if !expired { // 1 month
		return nil
	}
	if err := os.Remove(pTempFile); err != nil {
		return err
	}
	_, err = os.Create(pTempFile) // nolint: gosec
	return err
}

func (p *sStream) loadFileChunk() ([]byte, error) {
	var lastErr error
	for i := uint64(0); i <= p.fRetryNum; i++ {
		req := request.NewLoadRequest(
			p.fFileInfo.GetName(),
			p.fPosition/p.fChunkSize,
			p.fPersonal,
		)
		resp, err := p.fHlkClient.FetchRequest(p.fContext, p.fAliasName, req)
		if err != nil {
			lastErr = err
			continue
		}
		if resp.GetCode() != http.StatusOK {
			lastErr = err
			continue
		}
		hash, ok := resp.GetHead()[hls_filesharer_settings.CHeaderFileHash]
		if !ok || p.fFileInfo.GetHash() != hash {
			return nil, ErrGotAnotherHash
		}
		return resp.GetBody(), nil
	}
	return nil, errors.Join(ErrRetryFailed, lastErr)
}

func (p *sStream) loadFileChunkFromTemp() ([]byte, error) {
	f, err := os.Open(p.fTempFile)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	if _, err := f.Seek(int64(p.fTempPosition), io.SeekStart); err != nil { // nolint: gosec
		return nil, err
	}
	chunk := make([]byte, p.fChunkSize)
	n, err := f.Read(chunk)
	if err != nil {
		return nil, err
	}
	p.fTempPosition += uint64(n) // nolint: gosec
	return chunk[:n], nil
}

func (p *sStream) appendChunkToTempFile(pChunk []byte) error {
	f, err := os.OpenFile(p.fTempFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	if _, err := f.Write(pChunk); err != nil {
		return err
	}
	return nil
}
