package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2/hybrid"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/network/adapters"
	"github.com/number571/hidden-lake/pkg/network/adapters/https"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"

	hla_https_settings "github.com/number571/hidden-lake/internal/adapters/https/pkg/settings"
)

const (
	echoTemplate = "echo: %s;"
)

func main() {
	networks := build.GetNetworks()
	delete(networks, build.CDefaultNetwork)

	lenNetworks := len(networks)
	if lenNetworks == 0 {
		panic("networks is null")
	}

	retries, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(lenNetworks)

	for networkKey := range networks {
		go func(nk string) {
			defer wg.Done()

			networkByKey, _ := build.GetNetwork(networkKey)
			connections := networkByKey.FConnections.GetByScheme(hla_https_settings.CAppAdapterName)
			if len(connections) == 0 {
				return // pass another adapter
			}

			respTime, err := doTestRequest(nk, retries)
			if err != nil {
				log.Printf("%s: network '%s' has error: %s", hla_https_settings.CAppAdapterName, nk, err.Error())
				return
			}
			log.Printf("%s: network '%s' is working successfully; response time %s", hla_https_settings.CAppAdapterName, nk, respTime)
		}(networkKey)
	}

	wg.Wait()
}

func doTestRequest(networkKey string, retries int) (time.Duration, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		node1, key1 = newNode(networkKey, "node1")
		node2, key2 = newNode(networkKey, "node2")
	)

	go func() { _ = node1.Run(ctx) }()
	go func() { _ = node2.Run(ctx) }()

	_, pKey := exchangeKeys(node1, node2, key1, key2)
	startTime := time.Now()

	msg := "hello, world!"
	for i := 0; i < retries; i++ {
		rsp, err := node1.FetchRequest(
			ctx,
			pKey,
			request.NewRequestBuilder().WithBody([]byte(msg)).Build(),
		)
		if err != nil {
			return 0, err
		}
		if string(rsp.GetBody()) != fmt.Sprintf(echoTemplate, msg) {
			return 0, errors.New("got invalid response") // nolint: err113
		}
	}

	return time.Since(startTime), nil
}

func newNode(networkKey string, name string) (network.IHiddenLakeNode, layer2.IParticipantKey) {
	privKey := asymmetric.NewPrivKey()
	adapterSettings := adapters.NewSettingsByNetworkKey(networkKey)
	node, err := network.NewHiddenLakeNode(
		network.NewSettings(&network.SSettings{
			FAdapterSettings: adapterSettings,
		}),
		func() layer2.IScheme {
			scheme, _ := hybrid.NewScheme(privKey, adapterSettings.GetMessageSizeBytes())
			return scheme
		}(),
		layer2.NewKeysContainer(),
		func() database.IKVDatabase {
			kv, err := database.NewKVDatabase(name + "_" + networkKey + ".db")
			if err != nil {
				panic(err)
			}
			return kv
		}(),
		https.NewHTTPSAdapter(
			https.NewSettings(&https.SSettings{
				FAdapterSettings: adapterSettings,
			}),
			cache.NewLRUCache(build.GetSettings().FStorageManager.FCacheHashesCap),
			func() []string {
				networkByKey, _ := build.GetNetwork(networkKey)
				connections := networkByKey.FConnections.GetByScheme(hla_https_settings.CAppAdapterName)
				if len(connections) == 0 {
					panic("len conns == 0")
				}
				p, err := getPasswordByName("users.json", networkKey, name)
				if err != nil {
					panic(err)
				}
				c := strings.NewReplacer(
					"__username__", name,
					"__password__", p,
				).Replace(connections[0])
				return []string{c}
			},
			func() *tls.Certificate {
				cert, err := getCertificate(".", "localhost")
				if err != nil {
					panic(err)
				}
				return &cert
			}(),
			func() *x509.CertPool {
				certPool, err := getCertPool(".")
				if err != nil {
					panic(err)
				}
				return certPool
			}(),
		),
		func(_ context.Context, _ layer2.IParticipantKey, r request.IRequest) (response.IResponse, error) {
			rsp := []byte(fmt.Sprintf(echoTemplate, string(r.GetBody())))
			return response.NewResponseBuilder().WithBody(rsp).Build(), nil
		},
	)
	if err != nil {
		panic(err)
	}
	return node, privKey.GetPubKey()
}

func exchangeKeys(hlNode1, hlNode2 network.IHiddenLakeNode, key1, key2 layer2.IParticipantKey) (layer2.IParticipantKey, layer2.IParticipantKey) {
	node1 := hlNode1.GetOriginNode()
	node2 := hlNode2.GetOriginNode()

	node1.GetKeysContainer().Add(key2)
	node2.GetKeysContainer().Add(key1)

	return key1, key2
}

/*
	{
	    "8Jkl93Mdk93md1bz": {
	        "node1": "<INSERT>",
	        "node2": "<INSERT>"
	    }
	}
*/
func getPasswordByName(pPath, pNetworkKey, pName string) (string, error) {
	creds := map[string]map[string]string{}

	data, err := os.ReadFile(pPath) // nolint: gosec
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(data, &creds); err != nil {
		return "", err
	}

	network, ok := creds[pNetworkKey]
	if !ok {
		return "", errors.New("network not found") // nolint: err113
	}

	p, ok := network[pName]
	if !ok {
		return "", errors.New("user not found") // nolint: err113
	}

	return p, nil
}

func getCertPool(pPath string) (*x509.CertPool, error) {
	certsDir := filepath.Join(pPath, hla_https_settings.CPathCerts)

	if err := os.MkdirAll(certsDir, 0700); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(certsDir)
	if err != nil {
		return nil, err
	}

	certPool, err := x509.SystemCertPool()
	if err == nil {
		certPool = x509.NewCertPool()
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		certBytes, err := os.ReadFile(filepath.Join(certsDir, entry.Name())) // nolint:gosec
		if err != nil {
			return nil, err
		}
		if ok := certPool.AppendCertsFromPEM(certBytes); !ok {
			return nil, errors.New("append certs from pem") // nolint: err113
		}
	}

	return certPool, nil
}

func getCertificate(pPath, externalAddr string) (tls.Certificate, error) {
	var (
		certFilePath = filepath.Join(pPath, hla_https_settings.CPathCert)
		keyFilePath  = filepath.Join(pPath, hla_https_settings.CPathKey)
	)

	var (
		_, err1 = os.Stat(certFilePath)
		_, err2 = os.Stat(keyFilePath)
	)
	if err1 == nil && err2 == nil {
		certBytes, err := os.ReadFile(certFilePath) // nolint:gosec
		if err != nil {
			return tls.Certificate{}, err
		}
		keyBytes, err := os.ReadFile(keyFilePath) // nolint:gosec
		if err != nil {
			return tls.Certificate{}, err
		}
		return tls.X509KeyPair(certBytes, keyBytes)
	}

	if !os.IsNotExist(err1) || !os.IsNotExist(err2) {
		return tls.Certificate{}, errors.Join(err1, err2)
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	r := random.NewRandom().GetString(16)
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Organization: []string{r}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(5, 0, 0), // Valid for 5 years
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := net.ParseIP(externalAddr); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, externalAddr)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	certBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err := os.WriteFile(certFilePath, certBytes, 0600); err != nil {
		return tls.Certificate{}, err
	}

	keyBytes := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	if err := os.WriteFile(keyFilePath, keyBytes, 0600); err != nil {
		return tls.Certificate{}, err
	}

	return tls.X509KeyPair(certBytes, keyBytes)
}
