package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hlk_config "github.com/number571/hidden-lake/pkg/api/kernel/config"
	hls_client "github.com/number571/hidden-lake/pkg/api/services/filesharer/client"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
)

const (
	tsFriendPubKeyStr = "PubKey{03A52F4F442B6A05A36BB9AB4DE727F9765641483EC4D9C014B31991B1A2DD365A966BCC78E27198F58254FB7B7926B9DDBAB6B104C71612A5C52210219637858261622C76690341E30ABC9A50B48202774929AF455B25E89572DBE6B5FB0A3594E9B74E16AC52F8A97B13BCD0425585DB6E19C2B8C474A668C96AEA23953AD83714511594FBC058864931975A9FE78F8CAA7A0F4AA10A0A9CEB19166BE75317209B740219CFE71C84360A44086F79AC82F1CAB437650942775024C5B28DCBCE2A87290780B3FAEA5D5672CB3B54B4C61981C37474BF807224320380AB1A6872CB2A735E3110BE2BD7C03430042DC1116F886CCD2BA9EEE3527D010ACE1BBDE2E7309ACB96B766B42C66CB348205B386600AACC73AA9B6DB390568D7720C0C5166556A8B10A2F5B83EC788260DA36B67579896983F24001CAB38C1489511CCB33699322C17EC1F7CC052DE628D9C101BB8999922091310F71D39CBAF7EC46F94FBB0914B3D00E80D2C4B3B6ED25E0D01A30AD730907140636119B5448D65B757628707E73B0A0B63766825BB2889721DD8B255E88A38502E5D7A4D9E6992B8B32C2946068DD799F7A922BBE8A1F32A7B94A6C0A9E59CDD485BD16C96A65076B68CBAB95CA40BDA2BA1443A8B6A2B339AB8B67896BCD652B5B68405F91EF3FA629B3681893AB5E15201D6795715726E2DB90E6E303C45255FF727BCB9584966636988331406E1237CF41F216A26DFB6C40D55B5A00654DF3A6F20C945507817D9616D94F190D07B02E5EA2B29E21127EABBE771639320B855B515D30CC8AE3B2C3B90AE4929C3F8824500EB5840A80CA3A06F7B83213B6745697194D0086E9D96B4555621B0826E6975131211A1C81B620FFA9DD8186CA9841879056D84369CB434BA8829C7E771B3BC8494DCC3804408C70DA8CE2CA63939655481D64DA3456901A184E5698F54828E5549684AC165DFC42C9A972D75A58BAC9160A522CC204ABFB5117DF77ACFDABBC1B491962295B514F5A330A80A0654BFF46C33DC496BAB23225328A0E07148E473BC5C9A9DFAEC9DE05C998E291273D412804489427B3D81BA11D3E26665EA230733ABA9C60A4048B57E7C5388F8B9B1B162F5A8620A475A37AC80F304AE4B49933B49675FF85359E34D5FC8A480989F9AE014023283BE17CF14D96D52F85DCCE0A56A44AD13E955A2EA0901CAC7122634C43B2E267ABAED7029310BCAB8B27BE38A0AA0673EE2625F89056331A86AB7D5220DCA32EA613A4613C1909439AF211D29E955EEB52663ECAAA0DC517BD64E99523BE2357E4CF8BD9CEC51413451E3C6C15C9C90BC1420DB3CC07677BA06388BAAACBAD22185A1EC3BFB23C8D901A1056B8A89413DD6FB3F8147C6EB00303D6919E95C130D8AC086928A69A032907145EAB47E36B38C6520AE772CBF86DA318A44BDD3BB8020E90907E82743D3421EAB72E125CC9CF14095005B1DB84C20B5C54150B0042147ACCB90E558973BD294A1D25571183217C8A2335316F5220922856C7307BC5AC174F695BCD2167AE5A641849C3A4F4208AEE9AB11BC732C41518021C5A04A6600A21B464157BE61812B4A72BDE91083CA2674720492BC533953556AE91FA47823A36B0969884427C037B1C4FE1F4CBE11FDB8EF0C779BE2E33904CEF1B61627D22A7857FC784E1DE66CE8A81B2552F16D547285042C57E7524DDEA0C429F7598AD87899BAB8859EC7F66EF7BAC8769EF7FFC4139BCE8A0D2AED3F4F91D91EFB85821EF4111E16022AA96BF5F4F8BAE71D3693BF1739CA2C53A62BF658D1EB97F6B10B3FCA40E4995A851E2F1AB307393856ACD0F132636F08431012A2AC3FF7717559B9BB05324949E1A836480BF0802AA2E5E313A282130132D3CF9C38D73C9AF964311AE25E29C55C3184EB43D91743B039E63F60B7ACAB8A008C4D266240BE3C87A2D6E1D929E465205F7BDD4AF2C2551BB08FCD964286BCF5D9E51BDC68259974FE897A6F21E0AE2DDB7F25C0824BF1AB2D50E76550BEC9C25C29E9522EC74FE76BA6F63064BDAF95000F9DA774EB95070369C75356C24CDD2DA7DFE23012D7D42553A26677D44F2D040F492C713878B022E90F04C9E70CB1AC14287582B29B1DD9B1A8620E4AB04331D7B80A3F2DABEDF52471021C9E41B55CCF3ECE79EBFAEE2280B03F28D081CA791AA1E7A3C414CC144A16B1FE013B7B961D9D3059B3D0B01B5C33F1C73EA0183C6FBC3B2F780792168277AFD721AE19C13C3A2EACB8978C8EC3879FB3B281FC94A05C1665149592BB7E11CCF49270F1F1E0B59FE5B29B0C8E4D103E11DCE01942F553534CE9F73B98959F2C6FEDC8A0C946D55DA6DE6EFBBEA29A984EABA45F625CB1C940A450803FBB051B86E5EE10F8BF916A6CC9A24948C4AD34C20A75792C36804FAABCE29D2800BA81CCD2BB5AD6C0EE9F0C79969FFD63C29B25541348427B02E97AB0321E402F893ED092F8BD137D968E7FF1A21A9F29AE838664EB6A4B5D615724D36635B1F99A0870026FD015957E4692F5F2675A1A9A04E001A1172645761CEB4631F10531E02A1965E5B04DF81D5246629FC72204DD34145966B4F8427FE0ABDE73C643C106D82AAF479A9124C6F12174A5E6FF171EAEC4F7C50338D6806ED66592569EE29E75E553A6ECCC905E8B7965E11C262E8B2981DF3CACD942DDA7D583970BF5D8B7C1D253066FD7AE47391EE2A32D539B6F5470898A0958723A1AF30CD144EAE9C5A6FF0306E4A9B9786FE7467BE485AA3203F381074203F7FDFD6B6A7C4EFB4F7889956419883C3A171E0ADD66305DA8B2628A9F124095E75A1C58473D56E5354C371DA847761A0DDFE9128B6B8161B7756925AE8CF2391501F3CAECC555EA10AEBE0A87F8B2735F3B40B89F8B20357366E27E38D88F8B80E39E07052BCC10DF3870D9F27A6199034ECA1EA26D4E0FEA96E1D8BA75AAF4CB5867BA3DE41523D93A9E4D580297F7329BBAC6CBA1C460550A39807E967958521189AB7671C16ECB8771478362863E31B504081A0F9A6B5E220C503FC60C2602621F7351CEC162073D3181A06438E31C76C1871918978957B09753D8F687BDC5CB071769DBD5D8CD748E17E4D2AE8FB662FC244182C02BFA7E84CD2CBCDCF70860A87270121C205BF7A55330C181E14990E500E756791B09D1BBC118E2159E000CD6D0CEDDE3E5C979AC53A1F28AAB54F8007BEAD2412E562E31A34E06AFB58C303E8C2955840F2B6C04A800FAF2943829E137801B1F45A3689DF9498DF4E6F5FDE49795C69E1C444D03A934F33747F9DEE4863DB6209F23AF34A4502791885A4AA23982EBC7A5CB4774D5B4BD647F342CFF43266FA59BC9D0FC9A9D453A4D97FF6EDB1B5381DC809ABD72F56BCCC2121EE28A09D1674F49881B39151807C2FB011416C007860B123148D3521A7F4CFBEFA3EDA2B679C0CA6BB5AF29BD34258B30B7E9391093539C8F75E3710E7CC420722EBF22B5D34FB401867C15767AB447D03A28F433817303A8E4F1E9CD8643EA5024AF611358AE767C80E00D59CBCA12FA2D99C9AEC18503837BCC033ADEFFCD04B0FFE626FAA92A7FB64A2633A45370D12A1021CE8E865D7FFEB885401E19FE929263F11E94B4918374BD53B104460512F0E7960FC1C13726FFFAAE3E9E2274D3A2E5A32C6641E516CC383602FDE75F60F37303666CA71E5A68B9BD6B05F133B1D466D1CA72A91D0B187B02411A1C6C93B086EC89A396A49FB5470E3D03EF9CA4AB42D67F796B9CF26E7B27E25B2F837D3905C944F60722A35B1A68DA2E85CFB668EBFDC150DFDC1E144632B471F35D5511DE9F98560E56B20CD0F088150D86202ECC0FD4805D1CBA64B004369DB1A93BDCD75FEDF53E907929C18485AFBA054A7015833E216D47335F5BE7D0EFCEE545A791E45C95D6F1FD5C33223465CA7B935653F2A40C5134FD5EE51C7C7BB1824F75B35C066CC4A4857A96736734140967EEBF95DA786E96DEE1F9D965BB4BF6D76DAB886B35B97313A3C09BA34B7403BFBD0BE4F96A0B74FCBC96E882328E43E5CE430F5C8CD422A3DEB071EFFE914AD1E92B2232F5F201E9240B0D891F034018D41FC95B488ED66696DF4D952B956A786D42B7B627202B7AB11514CE1F4DD9CC1FF9E5300CB8611BAE6559ED509DE432BD38840010779A3726DE7C52F421696FBBD56C4462C228BC63BF42490BB0EA150DBBAB87D3BA750B3AC2FC93CA471536E674B9BBE818F8A72ECA97034B2B8BC88D771F352DAFC7DDBFC1FEE15F261DD292EEAF6917DBF92F80CEBFE1CEDD6E6EEE3D2E5D142FC1296949B139CA18E5A0FB99355D4C42353F36A548E6D43F4E59883E6ED72F85193B9E0349AF9AAE865DF437C289C30D4408E4BD5742365ADFB268C82212BFFC5BA0689A12B1C262189DC2AF2CB8DE560265841292D04122F80B641FDE}"
)

var (
	tsFriendPubKey = asymmetric.LoadPubKey(tsFriendPubKeyStr)
)

func TestErrorsAPI(t *testing.T) {
	t.Parallel()

	client := hls_client.NewClient(
		hls_client.NewRequester("", &http.Client{}),
	)

	ctx := context.Background()
	if _, err := client.GetIndex(ctx); err == nil {
		t.Fatal("success incorrect getIndex")
	}
	if _, err := client.GetRemoteList(ctx, "", 0, false); err == nil {
		t.Fatal("success incorrect getRemoteList")
	}
	if _, err := client.GetRemoteFile(nil, ctx, "", "", false); err == nil {
		t.Fatal("success incorrect getRemoteFile")
	}
	if _, err := client.GetRemoteFileInfo(ctx, "", "", false); err == nil {
		t.Fatal("success incorrect getRemoteFileInfo")
	}
	if _, err := client.GetLocalList(ctx, "", 0); err == nil {
		t.Fatal("success incorrect getLocalList")
	}
	if err := client.GetLocalFile(nil, ctx, "", ""); err == nil {
		t.Fatal("success incorrect getLocalFile")
	}
	if err := client.PutLocalFile(ctx, "", "", nil); err == nil {
		t.Fatal("success incorrect putLocalFile")
	}
	if err := client.DelLocalFile(ctx, "", ""); err == nil {
		t.Fatal("success incorrect delLocalFile")
	}
	if _, err := client.GetLocalFileInfo(ctx, "", ""); err == nil {
		t.Fatal("success incorrect getLocalFileInfo")
	}
}

func TestHandleIndexAPI(t *testing.T) {
	t.Parallel()

	httpLogger := std_logger.NewStdLogger(
		func() std_logger.ILogging {
			logging, err := std_logger.LoadLogging([]string{})
			if err != nil {
				panic(err)
			}
			return logging
		}(),
		func(_ logger.ILogArg) string {
			return ""
		},
	)

	handler := HandleIndexAPI(httpLogger)
	if err := indexAPIRequestOK(handler); err != nil {
		t.Fatal(err)
	}
}

func indexAPIRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

var (
	_ hlk_client.IClient = &tsHLKClient{}
	_ hlk_client.IClient = &tsHLKClientWrapper{}
)

type tsHLKClientWrapper struct {
	fClient hlk_client.IClient
}

func (p *tsHLKClientWrapper) GetIndex(a context.Context) (string, error) {
	return p.fClient.GetIndex(a)
}
func (p *tsHLKClientWrapper) GetSettings(a context.Context) (hlk_config.IConfigSettings, error) {
	return p.fClient.GetSettings(a)
}
func (p *tsHLKClientWrapper) GetPubKey(a context.Context) (asymmetric.IPubKey, error) {
	return p.fClient.GetPubKey(a)
}
func (p *tsHLKClientWrapper) GetOnlines(a context.Context) ([]string, error) {
	return p.fClient.GetOnlines(a)
}
func (p *tsHLKClientWrapper) DelOnline(a context.Context, b string) error {
	return p.fClient.DelOnline(a, b)
}
func (p *tsHLKClientWrapper) GetFriends(a context.Context) (map[string]asymmetric.IPubKey, error) {
	return p.fClient.GetFriends(a)
}
func (p *tsHLKClientWrapper) AddFriend(a context.Context, b string, c asymmetric.IPubKey) error {
	return p.fClient.AddFriend(a, b, c)
}
func (p *tsHLKClientWrapper) DelFriend(a context.Context, b string) error {
	return p.fClient.DelFriend(a, b)
}
func (p *tsHLKClientWrapper) GetConnections(a context.Context) ([]string, error) {
	return p.fClient.GetConnections(a)
}
func (p *tsHLKClientWrapper) AddConnection(a context.Context, b string) error {
	return p.fClient.AddConnection(a, b)
}
func (p *tsHLKClientWrapper) DelConnection(a context.Context, b string) error {
	return p.fClient.DelConnection(a, b)
}
func (p *tsHLKClientWrapper) SendRequest(a context.Context, b string, c request.IRequest) error {
	return p.fClient.SendRequest(a, b, c)
}
func (p *tsHLKClientWrapper) FetchRequest(a context.Context, b string, c request.IRequest) (response.IResponse, error) {
	return p.fClient.FetchRequest(a, b, c)
}

type tsHLKClient struct {
	fFetchType  int
	fSettingsOK bool
	fPrivKey    asymmetric.IPrivKey
}

func newTsHLKClient(pFetchType int, pSettingsOK bool) *tsHLKClient {
	return &tsHLKClient{
		fFetchType:  pFetchType,
		fSettingsOK: pSettingsOK,
		fPrivKey:    asymmetric.NewPrivKey(),
	}
}

func (p *tsHLKClient) GetIndex(context.Context) (string, error) { return "", nil }
func (p *tsHLKClient) GetSettings(context.Context) (hlk_config.IConfigSettings, error) {
	if !p.fSettingsOK {
		return nil, errors.New("error") // nolint: err113
	}
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
	return map[string]asymmetric.IPubKey{
		"abc": tsFriendPubKey,
	}, nil
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

func (p *tsHLKClient) FetchRequest(_ context.Context, _ string, pReq request.IRequest) (response.IResponse, error) {
	switch p.fFetchType {
	case 3:
		if strings.HasPrefix(pReq.GetPath(), "/load") {
			resp := response.NewResponseBuilder().WithCode(200).WithHead(map[string]string{settings.CHeaderFileHash: "7d0c64e050a2c31cd2d5266b2923ca51b95e97e2dedfc39e4ce220b477683975ba032c6c3141bad8442af4943f91ac43"}).WithBody([]byte("hello, world!\n"))
			return resp.Build(), nil
		}
		resp := response.NewResponseBuilder().WithCode(200).WithBody([]byte(`{"name":"example.txt","size":14,"hash":"7d0c64e050a2c31cd2d5266b2923ca51b95e97e2dedfc39e4ce220b477683975ba032c6c3141bad8442af4943f91ac43"}`))
		return resp.Build(), nil
	case 2:
		resp := response.NewResponseBuilder().WithCode(200).WithBody([]byte(`[{"name":"example.txt","size":14,"hash":"7d0c64e050a2c31cd2d5266b2923ca51b95e97e2dedfc39e4ce220b477683975ba032c6c3141bad8442af4943f91ac43"}]`))
		return resp.Build(), nil
	case 1:
		return nil, errors.New("error") // nolint: err113
	case 0:
		resp := response.NewResponseBuilder().WithCode(200).WithBody([]byte(`{"name":"example.txt","size":14,"hash":"7d0c64e050a2c31cd2d5266b2923ca51b95e97e2dedfc39e4ce220b477683975ba032c6c3141bad8442af4943f91ac43"}`))
		return resp.Build(), nil
	case -1:
		resp := response.NewResponseBuilder().WithCode(500).WithBody([]byte(`500`))
		return resp.Build(), nil
	case -2:
		resp := response.NewResponseBuilder().WithCode(200).WithBody([]byte{1})
		return resp.Build(), nil
	case -3:
		resp := response.NewResponseBuilder().WithCode(200).WithBody([]byte(`{"name":"example1.txt","size":14,"hash":"7d0c64e050a2c31cd2d5266b2923ca51b95e97e2dedfc39e4ce220b477683975ba032c6c3141bad8442af4943f91ac43"}`))
		return resp.Build(), nil
	}
	panic("unknown fetch type")
}

type tsConfig struct {
}

func (p *tsConfig) GetSettings() config.IConfigSettings {
	return &config.SConfigSettings{
		FPageOffset: 10,
		FRetryNum:   3,
	}
}
func (p *tsConfig) GetAddress() config.IAddress {
	return nil
}
func (p *tsConfig) GetLogging() std_logger.ILogging {
	return nil
}
func (p *tsConfig) GetConnection() string {
	return ""
}
