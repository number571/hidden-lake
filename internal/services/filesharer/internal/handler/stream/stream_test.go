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
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	hls_config "github.com/number571/hidden-lake/internal/kernel/pkg/config"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/utils"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SStreamError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestStreamReader(t *testing.T) {
	t.Parallel()

	inputPath := "./testdata/"
	filename := "file.txt"
	fileBytes, err := os.ReadFile(inputPath + filename)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	hlkClient := newTsHLSClient(fileBytes)

	stream, _, err := BuildStreamReader(
		ctx,
		0,
		inputPath,
		"alias_name",
		hlkClient,
		utils.NewFileInfo(filename, hashing.NewHasher(fileBytes).ToString(), uint64(len(fileBytes))),
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
}

var (
	_ hlk_client.IClient = &tsHLSClient{}
)

type tsHLSClient struct {
	fCounter   int
	fFileBytes []byte
	fPrivKey   asymmetric.IPrivKey
}

func newTsHLSClient(pFileBytes []byte) *tsHLSClient {
	return &tsHLSClient{
		fFileBytes: pFileBytes,
		fPrivKey:   asymmetric.NewPrivKey(),
	}
}

func (p *tsHLSClient) GetIndex(context.Context) (string, error) { return "", nil }
func (p *tsHLSClient) GetSettings(context.Context) (hls_config.IConfigSettings, error) {
	return &hls_config.SConfigSettings{
		FPayloadSizeBytes: 104, // gRespSize + 1
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
		fileInfo := utils.NewFileInfoFromBytes("file.txt", p.fFileBytes)
		resp = response.NewResponseBuilder().WithCode(200).WithBody(encoding.SerializeJSON(fileInfo))
	case strings.Contains(r.GetPath(), "/load"):
		resp = response.NewResponseBuilder().WithCode(200).WithBody([]byte{p.fFileBytes[p.fCounter]})
		p.fCounter++
	default:
		return nil, errors.New("unknown path") // nolint:err113
	}
	return resp.Build(), nil
}
