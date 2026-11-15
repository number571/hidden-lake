package sealed

import (
	"errors"
	"net"
	"slices"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/kernel/pkg/app/config"
	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	"github.com/number571/hidden-lake/pkg/api/kernel/client/scheme"
)

var (
	_ IConfig = &sConfig{}
)

// IConfig is an in-memory kernel config: nothing is written to disk.
type IConfig interface {
	config.IConfig
	config.IWrapper
	config.IEditor

	AddEndpoint(string)
	DelEndpoint(string)
}

type sConfig struct {
	fMutex       sync.RWMutex
	fSettings    *config.SConfigSettings
	fLogging     logger.ILogging
	fAddress     *config.SAddress
	fServiceAddr string
	fEndpoints   []string
	fFriends     map[string]layer2.IParticipantKey
}

// NewConfig builds the config. Empty network key = defaults of hlk.yml (no settings).
// Only hybrid (PubKey) networks are supported; the default network is rejected.
// Endpoints are the relay addresses (TCP host:port), friends are alias -> PubKey strings.
func NewConfig(
	pNetworkKey string,
	pServiceAddr string,
	pEndpoints []string,
	pFriends map[string]string,
) (IConfig, error) {
	settings, err := newConfigSettings(pNetworkKey)
	if err != nil {
		return nil, err
	}

	friends, err := loadParticipantKeys(pFriends)
	if err != nil {
		return nil, err
	}

	for _, e := range pEndpoints {
		if !isValidEndpoint(e) {
			return nil, ErrInvalidEndpoint
		}
	}

	// Only warn/erro: info-level go-peer QB logs reveal which messages are
	// addressed to this node and the size of each sent payload (metadata the
	// hmac-keyed seen db is meant to hide), and stderr may be captured.
	logging, err := logger.LoadLogging([]string{logger.CLogWarn, logger.CLogErro})
	if err != nil {
		return nil, errors.Join(ErrLoadLogging, err)
	}

	cfg := &sConfig{
		fSettings:    settings,
		fLogging:     logging,
		fAddress:     &config.SAddress{},
		fServiceAddr: pServiceAddr,
		fEndpoints:   make([]string, 0, len(pEndpoints)),
		fFriends:     friends,
	}
	for _, e := range pEndpoints {
		cfg.AddEndpoint(e)
	}

	return cfg, nil
}

func newConfigSettings(pNetworkKey string) (*config.SConfigSettings, error) {
	if pNetworkKey == "" {
		return &config.SConfigSettings{}, nil
	}

	if pNetworkKey == build.CDefaultNetwork {
		return nil, ErrUnsupportedNetwork
	}
	network, ok := build.GetNetwork(pNetworkKey)
	if !ok {
		return nil, ErrNetworkNotFound
	}

	settings := &config.SConfigSettings{
		FNetworkKey:       pNetworkKey,
		FCryptoSchemeType: network.FCryptoSchemeType,
		FMessageSizeBytes: network.FMessageSizeBytes,
		FWorkSizeBits:     network.FWorkSizeBits,
		FFetchTimeoutMS:   network.FFetchTimeoutMS,
		FQueuePeriodMS:    network.FQueuePeriodMS,
	}
	if settings.GetCryptoSchemeType() != scheme.CHybridScheme {
		return nil, ErrUnsupportedNetwork
	}

	return settings, nil
}

func loadParticipantKeys(pFriends map[string]string) (map[string]layer2.IParticipantKey, error) {
	result := make(map[string]layer2.IParticipantKey, len(pFriends))
	for name, val := range pFriends {
		if !isValidAliasName(name) {
			return nil, ErrInvalidAliasName
		}
		pKey := loadPubKey(val)
		if pKey == nil {
			return nil, ErrInvalidParticipantKey
		}
		result[name] = pKey
	}
	if hasDuplicateParticipantKeys(result) {
		return nil, ErrDuplicateParticipantKey
	}
	return result, nil
}

// loadPubKey returns nil for an invalid key: go-peer panics on some
// garbage of the correct length (e.g. PubKey{FF...}).
func loadPubKey(pKey string) (rKey layer2.IParticipantKey) {
	defer func() {
		if recover() != nil {
			rKey = nil
		}
	}()
	pubKey := asymmetric.LoadPubKey(pKey)
	if pubKey == nil {
		return nil
	}
	return pubKey
}

// isValidAliasName: trimmed, non-empty, without "/" and control/format characters
// (the alias is sent as the value of the Hlk-Sender-Name header).
func isValidAliasName(pName string) bool {
	if pName == "" || pName != strings.TrimSpace(pName) || !utf8.ValidString(pName) {
		return false
	}
	for _, r := range pName {
		if r == '/' || unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
			return false
		}
	}
	return true
}

func isValidEndpoint(pEndpoint string) bool {
	host, port, err := net.SplitHostPort(pEndpoint)
	if err != nil || host == "" {
		return false
	}
	portNum, err := strconv.ParseUint(port, 10, 16)
	return err == nil && portNum > 0
}

func hasDuplicateParticipantKeys(pFriends map[string]layer2.IParticipantKey) bool {
	mapping := make(map[string]struct{}, len(pFriends))
	for _, pKey := range pFriends {
		pStr := pKey.ToString()
		if _, ok := mapping[pStr]; ok {
			return true
		}
		mapping[pStr] = struct{}{}
	}
	return false
}

func (p *sConfig) GetConfig() config.IConfig {
	return p
}

func (p *sConfig) GetEditor() config.IEditor {
	return p
}

func (p *sConfig) UpdateFriends(pFriends map[string]layer2.IParticipantKey) error {
	for name := range pFriends {
		if !isValidAliasName(name) {
			return ErrInvalidAliasName
		}
	}
	if hasDuplicateParticipantKeys(pFriends) {
		return ErrDuplicateParticipantKey
	}

	friends := make(map[string]layer2.IParticipantKey, len(pFriends))
	for k, v := range pFriends {
		friends[k] = v
	}

	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fFriends = friends
	return nil
}

func (p *sConfig) GetSettings() config.IConfigSettings {
	return p.fSettings
}

func (p *sConfig) GetLogging() logger.ILogging {
	return p.fLogging
}

func (p *sConfig) GetAddress() config.IAddress {
	return p.fAddress
}

func (p *sConfig) GetFriends() map[string]layer2.IParticipantKey {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	result := make(map[string]layer2.IParticipantKey, len(p.fFriends))
	for k, v := range p.fFriends {
		result[k] = v
	}
	return result
}

func (p *sConfig) GetEndpoints() []string {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return slices.Clone(p.fEndpoints)
}

func (p *sConfig) AddEndpoint(pEndpoint string) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if slices.Contains(p.fEndpoints, pEndpoint) {
		return
	}
	p.fEndpoints = append(p.fEndpoints, pEndpoint)
}

func (p *sConfig) DelEndpoint(pEndpoint string) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fEndpoints = slices.DeleteFunc(p.fEndpoints, func(e string) bool { return e == pEndpoint })
}

func (p *sConfig) GetService(pName string) (string, bool) {
	if pName != hls_messenger_settings.CAppShortName {
		return "", false
	}
	return p.fServiceAddr, true
}
