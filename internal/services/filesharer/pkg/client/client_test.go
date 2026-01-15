package client

// import (
// 	"context"
// 	"testing"

// 	"github.com/number571/go-peer/pkg/crypto/asymmetric"
// 	hls_config "github.com/number571/hidden-lake/internal/kernel/pkg/config"
// 	"github.com/number571/hidden-lake/pkg/request"
// 	"github.com/number571/hidden-lake/pkg/response"
// )

// func TestError(t *testing.T) {
// 	t.Parallel()

// 	str := "value"
// 	err := &SClientError{str}
// 	if err.Error() != errPrefix+str {
// 		t.Fatal("incorrect err.Error()")
// 	}
// }

// func TestClient(t *testing.T) {
// 	t.Parallel()

// 	client := NewClient(
// 		NewBuilder(),
// 		NewRequester(newTsHLSClient()),
// 	)

// 	ctx := context.Background()
// 	if _, err := client.GetListFiles(ctx, "alias_name", 0); err != nil {
// 		t.Fatal(err)
// 	}
// 	if _, err := client.LoadFileChunk(ctx, "alias_name", "filename", 0); err != nil {
// 		t.Fatal(err)
// 	}
// }

// type tsHLSClient struct {
// 	fPrivKey asymmetric.IPrivKey
// }

// func newTsHLSClient() *tsHLSClient {
// 	return &tsHLSClient{
// 		fPrivKey: asymmetric.NewPrivKey(),
// 	}
// }

// func (p *tsHLSClient) GetIndex(context.Context) (string, error) { return "", nil }
// func (p *tsHLSClient) GetSettings(context.Context) (hls_config.IConfigSettings, error) {
// 	return nil, nil
// }

// func (p *tsHLSClient) GetPubKey(context.Context) (asymmetric.IPubKey, error) {
// 	return p.fPrivKey.GetPubKey(), nil
// }

// func (p *tsHLSClient) GetOnlines(context.Context) ([]string, error) { return nil, nil }
// func (p *tsHLSClient) DelOnline(context.Context, string) error      { return nil }

// func (p *tsHLSClient) GetFriends(context.Context) (map[string]asymmetric.IPubKey, error) {
// 	return nil, nil
// }

// func (p *tsHLSClient) AddFriend(context.Context, string, asymmetric.IPubKey) error { return nil }
// func (p *tsHLSClient) DelFriend(context.Context, string) error                     { return nil }

// func (p *tsHLSClient) GetConnections(context.Context) ([]string, error) { return nil, nil }
// func (p *tsHLSClient) AddConnection(context.Context, string) error      { return nil }
// func (p *tsHLSClient) DelConnection(context.Context, string) error      { return nil }

// func (p *tsHLSClient) SendRequest(context.Context, string, request.IRequest) error {
// 	return nil
// }

// func (p *tsHLSClient) FetchRequest(context.Context, string, request.IRequest) (response.IResponse, error) {
// 	resp := response.NewResponseBuilder().WithCode(200).WithBody([]byte(`[{"name":"file.txt","hash":"114a856f792c4c292599dba6fa41adba45ef4f851b1d17707e2729651968ff64be375af9cff6f9547b878d5c73c16a11","size":500}]`))
// 	return resp.Build(), nil
// }
