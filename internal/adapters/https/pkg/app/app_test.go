package app

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/payload"
	testutils_gopeer "github.com/number571/go-peer/test/utils"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/adapters/https/pkg/app/config"
	"github.com/number571/hidden-lake/internal/adapters/https/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/pkg/api/adapters/http/client"
	"github.com/number571/hidden-lake/pkg/network/adapters"
	testutils "github.com/number571/hidden-lake/test/utils"
)

var (
	tgFlags = flag.NewFlagsBuilder(
		flag.NewFlagBuilder("-v", "--version").
			WithDescription("print version of application"),
		flag.NewFlagBuilder("-h", "--help").
			WithDescription("print information about application"),
		flag.NewFlagBuilder("-p", "--path").
			WithDescription("set path to config, database files").
			WithDefinedValue("."),
		flag.NewFlagBuilder("-n", "--network").
			WithDescription("set network key of connections from build").
			WithDefinedValue(""),
	).Build()
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

const tcPathConfig = "./testdata/"
const tcPathConfigCerts = tcPathConfig + "/certs"

const tcDataConfig = `settings:
  message_size_bytes: 8192
`

func TestInitApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles(tcPathConfig)
	defer testDeleteFiles(tcPathConfig)

	if err := os.WriteFile(tcPathConfig+"hla-http.yml", []byte(tcDataConfig), 0600); err != nil {
		t.Fatal(err)
	}

	app, err := InitApp([]string{"--path", tcPathConfig}, tgFlags)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := InitApp([]string{"--path", tcPathConfig + "/failed"}, tgFlags); err == nil {
		t.Fatal("success init app with invalid config")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
}

func testDeleteFiles(path string) {
	_ = os.RemoveAll(path + settings.CPathYML)
	_ = os.RemoveAll(path + settings.CPathDB)
	_ = os.RemoveAll(path + settings.CPathKey)
	_ = os.RemoveAll(path + settings.CPathCert)
}

func testGetPEMFromTLS(cert tls.Certificate) (string, error) {
	var b bytes.Buffer
	for _, derBytes := range cert.Certificate {
		err := pem.Encode(&b, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: derBytes,
		})
		if err != nil {
			return "", err
		}
	}
	return b.String(), nil
}

func TestApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles("./")
	defer testDeleteFiles("./")

	// Run application
	cfg, err := config.BuildConfig(settings.CPathYML, &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessageSizeBytes: 8192,
			FWorkSizeBits:     10,
			FNetworkKey:       "_",
			FDatabaseEnabled:  true,
		},
		FAddress: &config.SAddress{
			FCertTLS:  "localhost",
			FInternal: testutils.TgAddrs[23],
			FExternal: testutils.TgAddrs[28],
		},
		FAuthMapper: map[string]string{
			"username1": "password1",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	cert, err := getCertificate(tcPathConfigCerts, cfg.GetAddress().GetCertTLS())
	if err != nil {
		t.Fatal(err)
	}

	certPool, err := getCertPool(tcPathConfigCerts)
	if err != nil {
		t.Fatal(err)
	}

	app := NewApp(cfg, &cert, certPool, ".")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	intClient := client.NewClient(
		client.NewRequester(
			testutils.TgAddrs[23],
			&http.Client{Timeout: time.Minute},
			adapters.NewSettings(&adapters.SSettings{
				FMessageSizeBytes: 8192,
				FWorkSizeBits:     10,
				FNetworkKey:       "_",
			}),
		),
	)

	err1 := testutils_gopeer.TryN(
		50,
		10*time.Millisecond,
		func() error {
			err := intClient.GetIndex(context.Background(), settings.CAppAdapterName)
			return err
		},
	)
	if err1 != nil {
		t.Fatal(err1)
	}

	msgBytes := []byte("hello, world!")
	msgBytes = append(msgBytes, random.NewRandom().GetBytes(uint64(8192-len(msgBytes)))...) //nolint:gosec
	netMsg := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{
				FWorkSizeBits: 10,
				FNetworkKey:   "_",
			}),
		}),
		payload.NewPayload32(build.GetSettings().FProtoMask.FNetwork, msgBytes),
	)

	if err := intClient.ProduceMessage(ctx, netMsg); err != nil {
		t.Fatal(err)
	}

	rootCAs := x509.NewCertPool()
	pemCert, err := testGetPEMFromTLS(cert)
	if err != nil {
		t.Fatal(err)
	}
	if ok := rootCAs.AppendCertsFromPEM([]byte(pemCert)); !ok {
		t.Fatal("failed append certs from pem")
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS13,
				RootCAs:    rootCAs,
			},
		},
	}

	msg2Bytes := []byte("hello, world!")
	msg2Bytes = append(msg2Bytes, random.NewRandom().GetBytes(uint64(8192-len(msg2Bytes)))...) //nolint:gosec
	net2Msg := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{
				FWorkSizeBits: 10,
				FNetworkKey:   "_",
			}),
		}),
		payload.NewPayload32(build.GetSettings().FProtoMask.FNetwork, msg2Bytes),
	)

	_, err = api.Request(
		ctx,
		httpClient,
		http.MethodPost,
		"https://"+testutils.TgAddrs[28]+settings.CHandleAdapterProducePath+"?sid=username1",
		http.Header{settings.CAuthTokenHeader: []string{"password1"}},
		net2Msg.ToString(),
	)
	if err != nil {
		t.Fatal(err)
	}

	// try twice running
	go func() {
		if err := app.Run(ctx); err == nil {
			t.Error("success double run")
		}
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	// try twice running
	go func() {
		if err := app.Run(ctx1); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
}
