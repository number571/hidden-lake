package stream

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hls_config "github.com/number571/hidden-lake/pkg/api/kernel/config"
	fileinfo "github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestStreamReader(t *testing.T) {
	t.Parallel()

	inputPath := "./testdata/"
	filename := "file.txt"
	tempname := "hls-filesharer-bf880b2a.tmp"

	_ = os.Remove(inputPath + tempname)
	defer func() { _ = os.Remove(inputPath + tempname) }()

	fileBytes, err := os.ReadFile(inputPath + filename)
	if err != nil {
		t.Fatal(err)
	}

	offset := 100
	if err := os.WriteFile(inputPath+tempname, fileBytes[:offset], 0600); err != nil {
		t.Fatal(err)
	}

	stream, err := BuildStreamReader(
		context.Background(),
		0,
		inputPath+tempname,
		"alias_name",
		newTsHLSClient(0, fileBytes, offset),
		newFileInfoFromBytes(filename, fileBytes),
		false,
	)
	if err != nil {
		t.Fatal(err)
	}

	r := make([]byte, 0, 128)
	b := make([]byte, 1)
	for {
		n, err := stream.Read(b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				r = append(r, b[0])
				break
			}
			t.Fatal(err)
		}
		if n != 1 {
			t.Fatal("n != 1")
		}
		r = append(r, b[0])
	}

	if string(r) != string(fileBytes) {
		t.Fatal("string(r) != string(fileBytes)")
	}

	_ = os.Remove(inputPath + tempname)
	stream2, err := BuildStreamReader(
		context.Background(),
		0,
		inputPath+tempname,
		"alias_name",
		newTsHLSClient(1, fileBytes, 0),
		newFileInfoFromBytes(filename, fileBytes),
		false,
	)
	if err != nil {
		t.Fatal(err)
	}

	for {
		n, err := stream2.Read(b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				t.Fatal("none invalid hash error (1)")
			}
			if errors.Is(err, ErrGotAnotherHash) {
				break // ok
			}
			t.Fatal(err)
		}
		if n != 1 {
			t.Fatal("n != 1")
		}
	}

	_ = os.Remove(inputPath + tempname)
	stream3, err := BuildStreamReader(
		context.Background(),
		0,
		inputPath+tempname,
		"alias_name",
		newTsHLSClient(2, fileBytes, 0),
		newFileInfoFromBytes(filename, fileBytes),
		false,
	)
	if err != nil {
		t.Fatal(err)
	}

	for {
		n, err := stream3.Read(b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				t.Fatal("none invalid hash error (2)")
			}
			if errors.Is(err, ErrInvalidHash) {
				break // ok
			}
			t.Fatal(err)
		}
		if n != 1 {
			t.Fatal("n != 1")
		}
	}

	_ = os.Remove(inputPath + tempname)
	stream4, err := BuildStreamReader(
		context.Background(),
		0,
		inputPath+tempname,
		"alias_name",
		newTsHLSClient(3, fileBytes, 0),
		newFileInfoFromBytes(filename, fileBytes),
		false,
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := stream4.Read(b); !errors.Is(err, ErrLoadFileChunk) {
		t.Fatal("success read with invalid load chunk")
	}

	_ = os.Remove(inputPath + tempname)
	stream5, err := BuildStreamReader(
		context.Background(),
		0,
		inputPath+tempname,
		"alias_name",
		newTsHLSClient(4, fileBytes, 0),
		newFileInfoFromBytes(filename, fileBytes),
		false,
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := stream5.Read(b); !errors.Is(err, ErrInvalidResponseCode) {
		t.Fatal("success read with invalid response code")
	}
}

var (
	_ hlk_client.IClient = &tsHLSClient{}
)

type tsHLSClient struct {
	fCounter   int
	fErrType   int
	fFileBytes []byte
	fPrivKey   asymmetric.IPrivKey
}

func newTsHLSClient(pErrType int, pFileBytes []byte, pOffset int) *tsHLSClient {
	return &tsHLSClient{
		fFileBytes: pFileBytes,
		fErrType:   pErrType,
		fPrivKey:   asymmetric.NewPrivKey(),
		fCounter:   pOffset,
	}
}

func (p *tsHLSClient) GetIndex(context.Context) error { return nil }
func (p *tsHLSClient) GetSettings(context.Context) (hls_config.IConfigSettings, error) {
	return &hls_config.SConfigSettings{
		FPayloadSizeBytes: 197 + 1, // gRespSize + 1
	}, nil
}

func (p *tsHLSClient) GetPubKey(context.Context) (asymmetric.IPubKey, error) {
	return p.fPrivKey.GetPubKey(), nil
}

func (p *tsHLSClient) GetOnlines(context.Context) ([]string, error) { return nil, nil }
func (p *tsHLSClient) DelOnline(context.Context, string) error      { return nil }

func (p *tsHLSClient) GetFriends(context.Context) (map[string]asymmetric.IPubKey, error) {
	return nil, nil
}

func (p *tsHLSClient) AddFriend(context.Context, string, asymmetric.IPubKey) error { return nil }
func (p *tsHLSClient) DelFriend(context.Context, string) error                     { return nil }

func (p *tsHLSClient) GetConnections(context.Context) ([]string, error) { return nil, nil }
func (p *tsHLSClient) AddConnection(context.Context, string) error      { return nil }
func (p *tsHLSClient) DelConnection(context.Context, string) error      { return nil }

func (p *tsHLSClient) SendRequest(context.Context, string, request.IRequest) error {
	return nil
}

func (p *tsHLSClient) FetchRequest(c context.Context, s string, r request.IRequest) (response.IResponse, error) {
	var resp response.IResponseBuilder
	switch {
	case strings.Contains(r.GetPath(), "/info"):
		fileInfo := newFileInfoFromBytes("file.txt", p.fFileBytes)
		resp = response.NewResponseBuilder().WithCode(200).WithBody(encoding.SerializeJSON(fileInfo))
	case strings.Contains(r.GetPath(), "/load"):
		switch p.fErrType {
		case 3:
			return nil, errors.New("error") // nolint: err113
		case 4:
			return response.NewResponseBuilder().
				WithCode(500).
				WithHead(map[string]string{}).
				WithBody([]byte{}).Build(), nil
		}
		const localChunk = 5
		var respBytes []byte
		if p.fCounter+localChunk >= len(p.fFileBytes) {
			respBytes = p.fFileBytes[p.fCounter : p.fCounter+localChunk]
			p.fCounter += localChunk
		} else {
			respBytes = p.fFileBytes[p.fCounter:len(p.fFileBytes)]
			p.fCounter += len(p.fFileBytes) - p.fCounter
		}
		hash := "bf880b2af9d0babacc67a988d2b7b9b6630e131d3ad6a0b78aefac0eaca162e4a3453b27a16de790f7879df4cda4b8c9"
		switch p.fErrType {
		case 0:
			// ok
		case 1:
			hash = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
		case 2:
			respBytes[0] ^= 1
		}
		resp = response.NewResponseBuilder().
			WithCode(200).
			WithHead(map[string]string{
				"Content-Type":                          api.CApplicationOctetStream,
				hls_filesharer_settings.CHeaderFileHash: hash,
			}).
			WithBody(respBytes)
	default:
		return nil, errors.New("unknown path") // nolint:err113
	}
	return resp.Build(), nil
}

type sFileInfo struct {
	FName string
	FHash string
	FSize uint64
}

func newFileInfoFromBytes(pName string, b []byte) fileinfo.IFileInfo { // nolint: unparam
	return &sFileInfo{
		FName: pName,
		FHash: hashing.NewHasher(b).ToString(),
		FSize: uint64(len(b)),
	}
}

func (p *sFileInfo) GetName() string {
	return p.FName
}

func (p *sFileInfo) GetHash() string {
	return p.FHash
}

func (p *sFileInfo) GetSize() uint64 {
	return p.FSize
}

func (p *sFileInfo) ToBytes() []byte {
	return encoding.SerializeJSON(p)
}

func (p *sFileInfo) ToString() string {
	return string(p.ToBytes())
}
