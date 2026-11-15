package sealed

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"net/netip"
	"slices"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

const (
	cTokenSize = 32
)

var (
	_ IBoot = &sBoot{}
)

// IBoot is the parsed stdin line of the hlk-sealed contract (v3).
type IBoot interface {
	GetPrivKey() asymmetric.IPrivKey
	IsPrivKeyGenerated() bool
	GetConfig() IConfig
	GetAPIToken() string
	GetPushToken() string
	GetSeenDBPath() string
	GetSeenKey() []byte
}

type sBootJSON struct {
	FPrivKey     string            `json:"priv_key"`
	FNetworkKey  string            `json:"network_key"`
	FFriends     map[string]string `json:"friends"`
	FEndpoints   []string          `json:"endpoints"`
	FServiceAddr string            `json:"service_addr"`
	FAPIToken    string            `json:"api_token"`
	FPushToken   string            `json:"push_token"`
	FSeenDBPath  string            `json:"seen_db_path"`
	FSeenKey     string            `json:"seen_key"`
}

type sBoot struct {
	fPrivKey    asymmetric.IPrivKey
	fGenerated  bool
	fConfig     IConfig
	fAPIToken   string
	fPushToken  string
	fSeenDBPath string
	fSeenKey    []byte
}

// ParseBoot parses and validates one JSON line. An empty priv_key
// generates a new private key (IsPrivKeyGenerated = true).
func ParseBoot(pLine []byte) (IBoot, error) {
	var v sBootJSON
	decoder := json.NewDecoder(bytes.NewReader(pLine))
	decoder.DisallowUnknownFields() // client of the same version generates the line
	if err := decoder.Decode(&v); err != nil {
		// json errors may quote bytes of the line (secrets): not joined
		return nil, ErrDecodeBoot
	}

	if !isValidToken(v.FAPIToken) {
		return nil, ErrInvalidAPIToken
	}
	if !isValidToken(v.FPushToken) {
		return nil, ErrInvalidPushToken
	}
	if v.FSeenDBPath == "" {
		return nil, ErrInvalidSeenDBPath
	}
	if !isValidServiceAddr(v.FServiceAddr) {
		return nil, ErrInvalidServiceAddr
	}

	seenKey, err := hex.DecodeString(v.FSeenKey)
	if err != nil || len(seenKey) != cSeenKeySize {
		return nil, ErrInvalidSeenKey
	}

	privKey, generated, err := loadPrivKey(v.FPrivKey)
	if err != nil {
		return nil, err
	}

	cfg, err := NewConfig(v.FNetworkKey, v.FServiceAddr, v.FEndpoints, v.FFriends)
	if err != nil {
		return nil, err
	}

	return &sBoot{
		fPrivKey:    privKey,
		fGenerated:  generated,
		fConfig:     cfg,
		fAPIToken:   v.FAPIToken,
		fPushToken:  v.FPushToken,
		fSeenDBPath: v.FSeenDBPath,
		fSeenKey:    seenKey,
	}, nil
}

// isValidToken: exactly 64 hex characters (32 bytes).
func isValidToken(pToken string) bool {
	b, err := hex.DecodeString(pToken)
	return err == nil && len(b) == cTokenSize
}

// isValidServiceAddr: loopback ip:port with non-zero port.
func isValidServiceAddr(pAddr string) bool {
	addrPort, err := netip.ParseAddrPort(pAddr)
	return err == nil && addrPort.Addr().IsLoopback() && addrPort.Port() != 0
}

func loadPrivKey(pPrivKey string) (asymmetric.IPrivKey, bool, error) {
	if pPrivKey == "" {
		return asymmetric.NewPrivKey(), true, nil
	}
	privKey := loadPrivKeySafe(pPrivKey)
	if privKey == nil {
		return nil, false, ErrInvalidPrivKey
	}
	return privKey, false, nil
}

// loadPrivKeySafe returns nil for an invalid key: go-peer panics on some
// garbage of the correct length (e.g. PrivKey{00...}).
func loadPrivKeySafe(pPrivKey string) (rKey asymmetric.IPrivKey) {
	defer func() {
		if recover() != nil {
			rKey = nil
		}
	}()
	return asymmetric.LoadPrivKey(pPrivKey)
}

func (p *sBoot) GetPrivKey() asymmetric.IPrivKey {
	return p.fPrivKey
}

func (p *sBoot) IsPrivKeyGenerated() bool {
	return p.fGenerated
}

func (p *sBoot) GetConfig() IConfig {
	return p.fConfig
}

func (p *sBoot) GetAPIToken() string {
	return p.fAPIToken
}

func (p *sBoot) GetPushToken() string {
	return p.fPushToken
}

func (p *sBoot) GetSeenDBPath() string {
	return p.fSeenDBPath
}

func (p *sBoot) GetSeenKey() []byte {
	return slices.Clone(p.fSeenKey)
}
