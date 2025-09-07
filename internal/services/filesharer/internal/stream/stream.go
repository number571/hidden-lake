package stream

import (
	"bytes"
	"context"
	"crypto/sha512"
	"errors"
	"hash"
	"io"
	"os"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	hls_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	internal_utils "github.com/number571/hidden-lake/internal/services/filesharer/internal/utils"
	hls_filesharer_client "github.com/number571/hidden-lake/internal/services/filesharer/pkg/client"
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
	_ IReadSeeker = &sStream{}
)

type sStream struct {
	fContext   context.Context
	fRetryNum  uint64
	fTempFile  string
	fHlfClient hls_filesharer_client.IClient
	fAliasName string
	fFileInfo  IFileInfo
	fBuffer    []byte
	fPosition  uint64
	fTempBytes []byte
	fHasher    hash.Hash
	fChunkSize uint64
}

func BuildStream(
	pCtx context.Context,
	pRetryNum uint64,
	pTempFile string,
	pHlsClient hls_client.IClient,
	pAliasName string,
	pFileInfo IFileInfo,
) (IReadSeeker, error) {
	chunkSize, err := internal_utils.GetMessageLimit(pCtx, pHlsClient)
	if err != nil {
		return nil, errors.Join(ErrGetMessageLimit, err)
	}
	tempBytes, err := os.ReadFile(pTempFile) // nolint: gosec
	if err != nil {
		panic(err)
	}
	return &sStream{
		fContext:   pCtx,
		fRetryNum:  pRetryNum,
		fTempFile:  pTempFile,
		fTempBytes: tempBytes,
		fHlfClient: hls_filesharer_client.NewClient(
			hls_filesharer_client.NewBuilder(),
			hls_filesharer_client.NewRequester(pHlsClient),
		),
		fAliasName: pAliasName,
		fHasher:    sha512.New384(),
		fChunkSize: chunkSize,
		fFileInfo:  pFileInfo,
	}, nil
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
		p.fBuffer = chunk
		if err := p.appendToTempFile(p.fBuffer); err != nil {
			return 0, errors.Join(ErrWriteFileChunk, err)
		}
	}

	n := copy(b, p.fBuffer)
	p.fBuffer = p.fBuffer[n:]
	p.fPosition += uint64(n) //nolint:gosec

	if _, err := p.fHasher.Write(b[:n]); err != nil {
		return 0, errors.Join(ErrWriteFileChunk, err)
	}

	if p.fPosition < p.fFileInfo.GetSize() {
		return n, nil
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
