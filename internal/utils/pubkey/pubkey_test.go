package pubkey

import (
	"context"
	"errors"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hlk_config "github.com/number571/hidden-lake/pkg/api/kernel/config"
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

func TestFriendPubKeyByAliasName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	if _, err := GetFriendPubKeyByAliasName(ctx, newTsHLKClient(false), "abc"); err == nil {
		t.Fatal("success get friend's public key with error")
	}
	hlkClient := newTsHLKClient(true)
	if _, err := GetFriendPubKeyByAliasName(ctx, hlkClient, "111"); err == nil {
		t.Fatal("success get public key of undefined friend")
	}
	pubKey, err := GetFriendPubKeyByAliasName(ctx, hlkClient, "abc")
	if err != nil {
		t.Fatal(err)
	}
	if pubKey.ToString() != hlkClient.fPrivKey.GetPubKey().ToString() {
		t.Fatal("got invalid public key")
	}
}

var (
	_ hlk_client.IClient = &tsHLKClient{}
)

type tsHLKClient struct {
	fFriendOK bool
	fPrivKey  asymmetric.IPrivKey
}

func newTsHLKClient(pFriendOK bool) *tsHLKClient {
	return &tsHLKClient{
		fFriendOK: pFriendOK,
		fPrivKey:  asymmetric.NewPrivKey(),
	}
}

func (p *tsHLKClient) GetIndex(context.Context) error { return nil }
func (p *tsHLKClient) GetSettings(context.Context) (hlk_config.IConfigSettings, error) {
	return &hlk_config.SConfigSettings{
		FPayloadSizeBytes: 1024,
	}, nil
}

func (p *tsHLKClient) GetPubKey(context.Context) (asymmetric.IPubKey, error) {
	return p.fPrivKey.GetPubKey(), nil
}

func (p *tsHLKClient) GetOnlines(context.Context) ([]string, error) {
	return []string{"tcp://aaa"}, nil
}
func (p *tsHLKClient) DelOnline(context.Context, string) error { return nil }

func (p *tsHLKClient) GetFriends(context.Context) (map[string]asymmetric.IPubKey, error) {
	if !p.fFriendOK {
		return nil, errors.New("error") // nolint: err113
	}
	return map[string]asymmetric.IPubKey{"abc": p.fPrivKey.GetPubKey()}, nil
}

func (p *tsHLKClient) AddFriend(context.Context, string, asymmetric.IPubKey) error { return nil }
func (p *tsHLKClient) DelFriend(context.Context, string) error                     { return nil }

func (p *tsHLKClient) GetConnections(context.Context) ([]string, error) {
	return []string{"tcp://aaa"}, nil
}
func (p *tsHLKClient) AddConnection(context.Context, string) error { return nil }
func (p *tsHLKClient) DelConnection(context.Context, string) error { return nil }

func (p *tsHLKClient) SendRequest(context.Context, string, request.IRequest) error {
	return nil
}

func (p *tsHLKClient) FetchRequest(context.Context, string, request.IRequest) (response.IResponse, error) {
	resp := response.NewResponseBuilder().WithCode(200).WithBody([]byte(`[{"name":"file.txt","size":500,"hash":"114a856f792c4c292599dba6fa41adba45ef4f851b1d17707e2729651968ff64be375af9cff6f9547b878d5c73c16a11"}]`))
	return resp.Build(), nil
}
