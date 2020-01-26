package settings

import (
	"../models"
	"../utils"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"github.com/number571/gopeer"
	"math/big"
	"os"
	"time"
)

func InitializeCFG(cfgname string) {
	if !utils.FileIsExist(cfgname) {
		cfg := newConfig()
		cfgJSON, err := json.MarshalIndent(cfg, "", "\t")
		if err != nil {
			panic("can't encode config")
		}
		os.Mkdir(PATH_TLS, 0777)
		createRandomCertificate(2048, cfg)
		utils.WriteFile(cfgname, string(cfgJSON))
	}
	cfgJSON := utils.ReadFile(cfgname)
	err := json.Unmarshal([]byte(cfgJSON), &CFG)
	if err != nil {
		panic("can't decode config")
	}
}

// Create RSA certificates.
func createRandomCertificate(bits int, cfg *models.Config) error {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(int64(gopeer.GenerateRandomIntegers(1)[0])),
		Subject: pkix.Name{
			Organization:  []string{"HIDDEN_LAKE"},
			Country:       []string{"NEW_COUNTRY"},
			Province:      []string{"NEW_PROVINCE"},
			Locality:      []string{"NEW_CITY"},
			StreetAddress: []string{"NEW_ADDRESS"},
			PostalCode:    []string{"NEW_POSTAL_CODE"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	priv, _ := rsa.GenerateKey(rand.Reader, bits)
	pub := &priv.PublicKey
	ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		return err
	}

	certOut, err := os.Create(cfg.Host.Http.Tls.Crt)
	if err != nil {
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: ca_b})
	certOut.Close()

	keyOut, err := os.Create(cfg.Host.Http.Tls.Key)
	if err != nil {
		return err
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()

	return nil
}

func newConfig() *models.Config {
	return &models.Config{
		Host: models.Host{
			Http: models.Http{
				Ipv4: "localhost",
				Port: ":7545",
				Tls: models.Tls{
					Crt: "tls/cert.crt",
					Key: "tls/cert.key",
				},
			},
			Tcp: models.Tcp{
				Ipv4: "localhost",
				Port: ":8080",
			},
		},
	}
}
