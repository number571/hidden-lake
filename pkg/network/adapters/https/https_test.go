package https

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/network/adapters"
	testutils "github.com/number571/hidden-lake/test/utils"

	"github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"
	hla_https_settings "github.com/number571/hidden-lake/internal/adapters/https/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
)

const (
	tcCertPEM = `-----BEGIN CERTIFICATE-----
MIIB+TCCAWKgAwIBAgIBATANBgkqhkiG9w0BAQsFADAcMRowGAYDVQQKExFJbi1N
ZW1vcnkgQ2VydCBDbzAeFw0yNjA0MTMyMDI4MDRaFw0yNzA0MTMyMDI4MDRaMBwx
GjAYBgNVBAoTEUluLU1lbW9yeSBDZXJ0IENvMIGfMA0GCSqGSIb3DQEBAQUAA4GN
ADCBiQKBgQCupauZEA1aikzXV7OTO98XDsS02ZoLeTSeTT+rzv9bl3uGAyys5KDW
eNLNjB0DuXcGtRJEZHztQRGC17bwgZ8INdxpBZuJgQKrYlYhPzmuob/pAqcYkjS5
I3t6VOC6GzIbGG8h9xxZrZolgL6S2JaW2VZ4yVaxbKUKPTES4vQhcQIDAQABo0sw
STAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDAYDVR0TAQH/
BAIwADAUBgNVHREEDTALgglsb2NhbGhvc3QwDQYJKoZIhvcNAQELBQADgYEAHL6f
tf1NX8ayJEGjd494zF9Ef0fQ61qf7eCNvVSMtErCg8m9FfFjH8NfEKQOlIfUeQOC
dJih8I1q1Kak2Q42nmiJxmBTmkVh+AjDIawhot9Bd+WtrTjK8wpZLglQTWKmJKS9
qH41/WZQjphw4+HvkKLvkhfZtzajQzKMBDp62cc=
-----END CERTIFICATE-----`

	// nolint:gosec
	tcKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQCupauZEA1aikzXV7OTO98XDsS02ZoLeTSeTT+rzv9bl3uGAyys
5KDWeNLNjB0DuXcGtRJEZHztQRGC17bwgZ8INdxpBZuJgQKrYlYhPzmuob/pAqcY
kjS5I3t6VOC6GzIbGG8h9xxZrZolgL6S2JaW2VZ4yVaxbKUKPTES4vQhcQIDAQAB
AoGAAQp4IsBnil7N2HBR+XjfR1CAnm23++uFnPZSGjo9Z8e+WVOGUT16y+xw05kx
lXnmGk8hkPDJLqAEDhouoYshB7WSvOVabBntk0tGy9ZAgk9hi/TMs9wKoZi1MxxA
17HflHvtLo2K0fNzxMrmZ0CwzGsH+HZAGnfvX4yBMrxPlTECQQDe2a1xeYRJq6Vt
KpsV+zlEj2JInO0/Bfuzq1hovAVSF/VLw1tI8x+5Y5C25UBZKmG7o9emrt5RQN/b
K+qWBVE9AkEAyKBgjJV3a3WufjiOx+L9PEBMCbKzKS1nfOxvGBYPJ7FxO2QF6So2
wreqgLVv7ePQCJns4teXTyQDI4eOVzXsRQJBANd8Lw15ziQaeLStrRa9POwBpazH
KVV2qKNcPPnRTWfLSOMAvSU2CmgOUaG43dcadzSkwmMn1ktFavCYb5au/5UCQQCy
Ou+q1LmzaGds0HffkYKgvQoP74YERcbTDwQepLIv9A4A0foCSrM9RocdMpJOBv1w
NrZgS2CrOPXk4W8NgOT1AkEAkD8sHBQi/NE3xi0Th4Znt9SDde1JjpQwxipvkMff
7Sqs+GxMrUa1bI4EarT6VjlxhEU3qR2cmt9c2UnVbJHpKw==
-----END RSA PRIVATE KEY-----`
)

// func testGenerateInMemCert() (tls.Certificate, error) {
// 	priv, _ := rsa.GenerateKey(rand.Reader, 1024)

// 	template := x509.Certificate{
// 		SerialNumber: big.NewInt(1),
// 		Subject: pkix.Name{
// 			Organization: []string{"In-Memory Cert Co"},
// 		},
// 		NotBefore:             time.Now(),
// 		NotAfter:              time.Now().AddDate(1, 0, 0), // Valid for 1 year
// 		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
// 		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
// 		BasicConstraintsValid: true,
// 		DNSNames:              []string{"localhost"},
// 	}

// 	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)

// 	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
// 	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

// 	return tls.X509KeyPair(certPEM, privPEM)
// }

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestHTTPSAdapter(t *testing.T) { // nolint: gocyclo, maintidx
	t.Parallel()

	// cert, err := testGenerateInMemCert()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	cert, err := tls.X509KeyPair([]byte(tcCertPEM), []byte(tcKeyPEM))
	if err != nil {
		t.Fatal(err)
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM([]byte(tcCertPEM)); !ok {
		t.Fatal("failed to parse any certificates from PEM data")
	}

	adapterSettings := adapters.NewSettings(&adapters.SSettings{FMessageSizeBytes: 8192})

	adapter1 := NewHTTPSAdapter(
		NewSettings(&SSettings{
			FAdapterSettings: adapterSettings,
			FServeSettings: &SServeSettings{
				FAddress: testutils.TgAddrs[19],
				FAuthMapper: map[string]string{
					"username1": "password1",
					"username2": "password2",
					"username3": "password3",
					"username4": "password4",
				},
				FRateLimitParams: [2]float64{.1, 1},
				FDataBrokerParam: 1,
			},
		}),
		cache.NewLRUCache(1024),
		func() []string { return nil },
		&cert,
		nil,
	)

	adapter2 := NewHTTPSAdapter(
		NewSettings(&SSettings{
			FAdapterSettings: adapterSettings,
			FServeSettings:   &SServeSettings{},
		}),
		cache.NewLRUCache(1024),
		func() []string {
			return []string{
				fmt.Sprintf("%s:%s@%s", "username1", "password1", testutils.TgAddrs[19]),
			}
		},
		nil,
		certPool,
	)

	onlines := adapter2.GetOnlines()
	if len(onlines) != 1 || onlines[0] != fmt.Sprintf("%s:%s@%s", "username1", "password1", testutils.TgAddrs[19]) {
		t.Fatal("adapter: get onlines")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { _ = adapter2.Run(ctx) }()
	go func() { _ = adapter1.Run(ctx) }()

	msgBytes := []byte("hello, world!")
	msgBytes = append(msgBytes, random.NewRandom().GetBytes(uint64(8192-len(msgBytes)))...) //nolint:gosec
	netMsg := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		payload.NewPayload32(build.GetSettings().FProtoMask.FNetwork, msgBytes),
	)

	time.Sleep(time.Second)

	if err := adapter2.Produce(ctx, netMsg); err != nil {
		t.Fatal(err)
	}
	// retry produce
	if err := adapter2.Produce(ctx, netMsg); err != nil {
		t.Fatal(err)
	}

	msg2, err := adapter1.Consume(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(msg2.GetHmac(), netMsg.GetHmac()) {
		t.Fatal("invalid hmac (2)")
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS13,
				RootCAs:    certPool,
			},
		},
	}

	_, err = api.Request(
		context.Background(),
		httpClient,
		http.MethodPost,
		"https://"+testutils.TgAddrs[19]+settings.CHandleAdapterConsumePath+"?sid=username",
		nil,
		nil,
	)
	if err == nil {
		t.Fatal("success request with invalid method")
	}

	_, err = api.Request(
		ctx,
		httpClient,
		http.MethodGet,
		"https://"+testutils.TgAddrs[19]+settings.CHandleAdapterProducePath+"?sid=username1",
		http.Header{hla_https_settings.CAuthTokenHeader: []string{"password1"}},
		netMsg.ToBytes(),
	)
	if err == nil {
		t.Fatal("success request with invalid method")
	}

	_, err = api.Request(
		ctx,
		httpClient,
		http.MethodPost,
		"https://"+testutils.TgAddrs[19]+settings.CHandleAdapterProducePath+"?sid=undefined",
		http.Header{hla_https_settings.CAuthTokenHeader: []string{"password1"}},
		netMsg.ToBytes(),
	)
	if err == nil {
		t.Fatal("success request with invalid sid")
	}

	msgBytesX := []byte("hello, world!")
	msgBytesX = append(msgBytesX, random.NewRandom().GetBytes(uint64(8192-len(msgBytesX)))...) //nolint:gosec
	netMsgX := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		payload.NewPayload32(111, msgBytesX),
	)

	_, err = api.Request(
		ctx,
		httpClient,
		http.MethodPost,
		"https://"+testutils.TgAddrs[19]+settings.CHandleAdapterProducePath+"?sid=username4",
		http.Header{hla_https_settings.CAuthTokenHeader: []string{"password4"}},
		netMsgX.ToBytes(),
	)
	if err == nil {
		t.Fatal("success request with invalid proto network")
	}

	_, err = api.Request(
		ctx,
		httpClient,
		http.MethodPost,
		"https://"+testutils.TgAddrs[19]+settings.CHandleAdapterProducePath+"?sid=username1",
		http.Header{hla_https_settings.CAuthTokenHeader: []string{"password1"}},
		encoding.HexEncode([]byte{1}),
	)
	if err == nil {
		t.Fatal("success request with invalid body (1)")
	}

	size := adapterSettings.GetMessageSizeBytes() + layer1.CMessageHeadSize
	_, err = api.Request(
		ctx,
		httpClient,
		http.MethodPost,
		"https://"+testutils.TgAddrs[19]+settings.CHandleAdapterProducePath+"?sid=username1",
		http.Header{hla_https_settings.CAuthTokenHeader: []string{"password1"}},
		encoding.HexEncode(random.NewRandom().GetBytes(size)),
	)
	if err == nil {
		t.Fatal("success request with invalid body (2)")
	}

	msgBytes2 := []byte("hello, world (222)!")
	msgBytes2 = append(msgBytes2, random.NewRandom().GetBytes(uint64(8192-len(msgBytes2)))...) //nolint:gosec
	netMsg2 := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		payload.NewPayload32(build.GetSettings().FProtoMask.FNetwork, msgBytes2),
	)

	chErr1 := make(chan error, 1)
	chErr2 := make(chan error, 1)

	go func() {
		time.Sleep(time.Second)
		_, err := api.Request(
			ctx,
			httpClient,
			http.MethodGet,
			"https://"+testutils.TgAddrs[19]+settings.CHandleAdapterConsumePath+"?sid=username2",
			http.Header{hla_https_settings.CAuthTokenHeader: []string{"password2"}},
			nil,
		)
		if err == nil {
			chErr1 <- errors.New("got double consume success") // nolint: err113
			return
		}
		_, err = api.Request(
			ctx,
			httpClient,
			http.MethodPost,
			"https://"+testutils.TgAddrs[19]+settings.CHandleAdapterProducePath+"?sid=username2",
			http.Header{hla_https_settings.CAuthTokenHeader: []string{"password2"}},
			netMsg2.ToBytes(),
		)
		chErr1 <- err
	}()

	go func() {
		msgBytes, err = api.Request(
			ctx,
			httpClient,
			http.MethodGet,
			"https://"+testutils.TgAddrs[19]+settings.CHandleAdapterConsumePath+"?sid=username2",
			http.Header{hla_https_settings.CAuthTokenHeader: []string{"password2"}},
			nil,
		)
		gotMsg, err := layer1.LoadMessage(adapterSettings, msgBytes)
		if err != nil {
			chErr2 <- err
			return
		}
		if gotMsg.ToString() != netMsg2.ToString() {
			chErr2 <- errors.New("got invalid message") // nolint: err113
			return
		}
		chErr2 <- nil
	}()

	if err := <-chErr1; err != nil {
		t.Fatal(err)
	}
	if err := <-chErr2; err != nil {
		t.Fatal(err)
	}

	msgBytesY := []byte("hello, world!")
	msgBytesY = append(msgBytesY, random.NewRandom().GetBytes(uint64(8192-len(msgBytesY)))...) //nolint:gosec
	netMsgY := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		payload.NewPayload32(build.GetSettings().FProtoMask.FNetwork, msgBytesY),
	)
	_, err = api.Request(
		ctx,
		httpClient,
		http.MethodPost,
		"https://"+testutils.TgAddrs[19]+settings.CHandleAdapterProducePath+"?sid=username2",
		http.Header{hla_https_settings.CAuthTokenHeader: []string{"password2"}},
		netMsgY.ToBytes(),
	)
	if err == nil {
		t.Fatal("success request with overflow rate limit")
	}

	msgBytesZ := []byte("hello, world!")
	msgBytesZ = append(msgBytesZ, random.NewRandom().GetBytes(uint64(8192-len(msgBytesZ)))...) //nolint:gosec
	netMsgZ := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		payload.NewPayload32(build.GetSettings().FProtoMask.FNetwork, msgBytesZ),
	)
	_, err = api.Request(
		ctx,
		httpClient,
		http.MethodPost,
		"https://"+testutils.TgAddrs[19]+settings.CHandleAdapterProducePath+"?sid=username3",
		http.Header{hla_https_settings.CAuthTokenHeader: []string{"password3"}},
		netMsgZ.ToBytes(),
	)
	if err == nil {
		t.Fatal("success request with overflow broker data")
	}

	_, err = api.Request(
		ctx,
		httpClient,
		http.MethodGet,
		"https://"+testutils.TgAddrs[19]+settings.CHandleAdapterConsumePath+"?sid=undefined",
		http.Header{hla_https_settings.CAuthTokenHeader: []string{"password2"}},
		nil,
	)
	if err == nil {
		t.Fatal("success consume with invalid sid")
	}

	if _, err := parseURL("\000"); err == nil || errors.Is(err, ErrNoPassword) {
		t.Fatal("parse url: success with invalid")
	}
	if _, err := parseURL("username@host"); !errors.Is(err, ErrNoPassword) {
		t.Fatal("parse url: success without password")
	}
}

func TestSettings(t *testing.T) {
	t.Parallel()

	sett := NewSettings(nil)
	if sett.GetAdapterSettings() == nil {
		t.Fatal("invalid adapter settings")
	}
}
