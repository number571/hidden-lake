package stream

import (
	"bytes"
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	hls_filesharer_client "github.com/number571/hidden-lake/internal/services/filesharer/pkg/client"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/utils"
)

func init() {
	v := []byte("init_value")
	h := sha512.Sum384(v)

	// maintaining the overall level of security and uniformity of the algorithms used
	if !bytes.Equal(h[:], hashing.NewHasher(v).ToBytes()) {
		panic("uses diff hash functions")
	}
}

var (
	_ io.ReadSeeker = &sStream{}
)

type ICallbackFunc func([]byte, uint64, uint64)

type sStream struct {
	fContext   context.Context
	fRetryNum  uint64
	fTempFile  string
	fHlfClient hls_filesharer_client.IClient
	fAliasName string
	fFileInfo  hls_filesharer_client.IFileInfo
	fBuffer    []byte
	fPosition  uint64
	fTempBytes []byte
	fHasher    hash.Hash
	fChunkSize uint64
	fCallback  ICallbackFunc
}

func BuildStream(
	pCtx context.Context,
	pRetryNum uint64,
	pInputPath string,
	pAliasName string,
	pHlkClient hlk_client.IClient,
	pFilename string,
	pCallback ICallbackFunc,
) (io.ReadSeeker, error) {
	hlfClient := hls_filesharer_client.NewClient(
		hls_filesharer_client.NewBuilder(),
		hls_filesharer_client.NewRequester(pHlkClient),
	)
	chunkSize, err := utils.GetMessageLimitOnLoadPage(pCtx, pHlkClient)
	if err != nil {
		return nil, errors.Join(ErrGetMessageLimit, err)
	}
	fileInfo, err := hlfClient.GetFileInfo(pCtx, pAliasName, pFilename)
	if err != nil {
		return nil, errors.Join(ErrGetFileInfo, err)
	}
	tempName := fmt.Sprintf(hls_filesharer_settings.CPathTMP, fileInfo.GetHash()[:8])
	tempFile := filepath.Join(pInputPath, tempName)
	tempBytes, err := readTempFile(tempFile, chunkSize)
	if err != nil {
		return nil, errors.Join(ErrReadTempFile, err)
	}
	return &sStream{
		fContext:   pCtx,
		fRetryNum:  pRetryNum,
		fTempFile:  tempFile,
		fTempBytes: tempBytes,
		fHlfClient: hlfClient,
		fAliasName: pAliasName,
		fHasher:    sha512.New384(),
		fChunkSize: chunkSize,
		fFileInfo:  fileInfo,
		fCallback:  pCallback,
	}, nil
}

func readTempFile(pTempFile string, pChunkSize uint64) ([]byte, error) {
	if _, err := os.Stat(pTempFile); errors.Is(err, os.ErrNotExist) {
		if _, err := os.Create(pTempFile); err != nil { // nolint: gosec
			return nil, err
		}
	}
	tempBytes, err := os.ReadFile(pTempFile) // nolint: gosec
	if err != nil {
		return nil, err
	}
	if uint64(len(tempBytes))%pChunkSize != 0 {
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

	if p.fCallback != nil {
		p.fCallback(b[:n], p.fPosition, p.fFileInfo.GetSize())
	}

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

func (p *sStream) Seek(offset int64, whence int) (int64, error) {
	select {
	case <-p.fContext.Done():
		return 0, io.ErrClosedPipe
	default:
	}

	var pos int64
	switch whence {
	case io.SeekStart:
		pos = offset
	case io.SeekCurrent:
		pos = int64(p.fPosition) + offset //nolint:gosec
	case io.SeekEnd:
		pos = int64(p.fFileInfo.GetSize()) + offset //nolint:gosec
	default:
		return 0, ErrInvalidWhence
	}
	if pos < 0 {
		return 0, ErrNegativePosition
	}
	p.fBuffer = p.fBuffer[:0]
	p.fPosition = uint64(pos)
	return pos, nil
}

func (p *sStream) loadFileChunk() ([]byte, error) {
	var lastErr error
	for i := uint64(0); i <= p.fRetryNum; i++ {
		chunk, err := p.fHlfClient.LoadFileChunk(
			p.fContext,
			p.fAliasName,
			p.fFileInfo.GetName(),
			p.fPosition/p.fChunkSize,
		)
		if err != nil {
			lastErr = err
			continue
		}
		return chunk, nil
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
