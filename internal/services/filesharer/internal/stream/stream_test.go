package stream

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	hls_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	hls_config "github.com/number571/hidden-lake/internal/kernel/pkg/config"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SStreamError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestStream(t *testing.T) {
	t.Parallel()

	filename := "file.txt"
	fileBytes, err := os.ReadFile("./testdata/" + filename)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	stream, err := BuildStream(
		ctx,
		0,
		newTsHLSClient(fileBytes),
		"alias_name",
		func() IFileInfo {
			hash := hashing.NewHasher(fileBytes).ToString()
			return NewFileInfo(filename, hash, uint64(len(fileBytes)))
		}(),
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
	_ hls_client.IClient = &tsHLSClient{}
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

func (p *tsHLSClient) FetchRequest(context.Context, string, request.IRequest) (response.IResponse, error) {
	resp := response.NewResponseBuilder().WithCode(200).WithBody([]byte{p.fFileBytes[p.fCounter]})
	p.fCounter++
	return resp.Build(), nil
}
