package app

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/types"
	build "github.com/number571/hidden-lake/build/environment"
	"github.com/number571/hidden-lake/internal/adapters/https/pkg/app/config"
	hla_https_settings "github.com/number571/hidden-lake/internal/adapters/https/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func InitApp(pArgs []string, pFlags flag.IFlags) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(pFlags.Get("-p").GetStringValue(pArgs), "/")
	if err := os.MkdirAll(inputPath, 0700); err != nil {
		return nil, errors.Join(ErrMkdirPath, err)
	}

	okLoaded, err := build.SetBuildByPath(inputPath)
	if err != nil {
		return nil, errors.Join(ErrSetBuild, err)
	}

	cfgPath := filepath.Join(inputPath, hla_https_settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil, pFlags.Get("-n").GetStringValue(pArgs))
	if err != nil {
		return nil, errors.Join(ErrInitConfig, err)
	}

	stdfLogger := std_logger.NewStdLogger(cfg.GetLogging(), std_logger.GetLogFunc())
	build.LogLoadedBuildFiles(hla_https_settings.GetAppShortNameFMT(), stdfLogger, okLoaded)

	cert, err := getCertificate(inputPath, cfg.GetAddress().GetCertTLS())
	if err != nil {
		return nil, errors.Join(ErrGetCertificate, err)
	}

	certPool, err := getCertPool(inputPath)
	if err != nil {
		return nil, errors.Join(ErrGetCertPool, err)
	}

	return NewApp(cfg, &cert, certPool, inputPath), nil
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
			return nil, ErrAddCertToPool
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
