package stream

import (
	"bytes"
	"context"
	"crypto/sha512"
	"errors"
	"hash"
	"io"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	internal_utils "github.com/number571/hidden-lake/internal/applications/filesharer/internal/utils"
	hlf_client "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/client"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
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
	fHlfClient hlf_client.IClient
	fAliasName string
	fFileInfo  IFileInfo
	fBuffer    []byte
	fPosition  uint64
	fHasher    hash.Hash
	fChunkSize uint64
}

func BuildStream(
	pCtx context.Context,
	pRetryNum uint64,
	pHlsClient hls_client.IClient,
	pAliasName string,
	pFileInfo IFileInfo,
) (IReadSeeker, error) {
	chunkSize, err := internal_utils.GetMessageLimit(pCtx, pHlsClient)
	if err != nil {
		return nil, errors.Join(ErrGetMessageLimit, err)
	}

	return &sStream{
		fContext:  pCtx,
		fRetryNum: pRetryNum,
		fHlfClient: hlf_client.NewClient(
			hlf_client.NewBuilder(),
			hlf_client.NewRequester(pHlsClient),
		),
		fAliasName: pAliasName,
		fHasher:    sha512.New384(),
		fChunkSize: chunkSize,
		fFileInfo:  pFileInfo,
	}, nil
}

func (p *sStream) Read(b []byte) (int, error) {
	if len(p.fBuffer) == 0 {
		chunk, err := p.loadFileChunk()
		if err != nil {
			return 0, errors.Join(ErrLoadFileChunk, err)
		}
		if _, err := p.fHasher.Write(chunk); err != nil {
			return 0, errors.Join(ErrWriteFileChunk, err)
		}
		p.fBuffer = chunk
	}

	n := copy(b, p.fBuffer)
	p.fBuffer = p.fBuffer[n:]
	p.fPosition += uint64(n)

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
	var pos int64
	switch whence {
	case io.SeekStart:
		pos = offset
	case io.SeekCurrent:
		pos = int64(p.fPosition) + offset
	case io.SeekEnd:
		pos = int64(p.fFileInfo.GetSize()) + offset
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
