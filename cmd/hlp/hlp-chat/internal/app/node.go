package app

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"io"
	"time"
	"unicode"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/cmd/hlp/hlp-chat/internal/database"
	"github.com/number571/hidden-lake/pkg/adapters"
	"github.com/number571/hidden-lake/pkg/adapters/tcp"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

func (p *sApp) getHLNode(
	pNetworkKey string,
	pWriter io.Writer,
) network.IHiddenLakeNode {
	networkByKey, ok := build.GetNetwork(pNetworkKey)
	if !ok {
		panic("network key undefined")
	}

	adapterSettings := adapters.NewSettingsByNetworkKey(pNetworkKey)
	connections := networkByKey.FConnections.GetByScheme("tcp")

	node := network.NewHiddenLakeNode(
		network.NewSettings(&network.SSettings{
			FAdapterSettings: adapterSettings,
		}),
		p.fChanKey,
		p.fDB.GetOrigin(),
		tcp.NewTCPAdapter(
			tcp.NewSettings(&tcp.SSettings{
				FAdapterSettings: adapterSettings,
			}),
			cache.NewLRUCache(2<<10),
			func() []string { return connections },
		),
		func(_ context.Context, pk asymmetric.IPubKey, r request.IRequest) (response.IResponse, error) {
			if r.GetHost() != cHiddenLakeProjectHost {
				return nil, nil // ignore message
			}

			pubKey, ok := validateSignature(r.GetHead(), r.GetBody())
			if !ok {
				return nil, nil // ignore message
			}

			msg := database.SMessage{FSendTime: time.Now(), FSender: pubKey, FMessage: string(r.GetBody())}
			if hasNotGraphicCharacters(msg.FMessage) {
				return nil, nil // ignore message
			}

			if err := p.fDB.Insert(pk, msg); err != nil {
				return nil, nil // ignore message
			}

			fmt.Fprintf(
				pWriter,
				cRecvMessageTeamplte,
				pubKey,
				msg.FMessage,
				msg.FSendTime.Format(time.DateTime),
			)
			return nil, nil
		},
	)

	node.GetOriginNode().GetMapPubKeys().SetPubKey(p.fChanKey.GetPubKey())
	return node
}

func validateSignature(head map[string]string, body []byte) (ed25519.PublicKey, bool) {
	pubkHex, ok1 := head["pubk"]
	saltHex, ok2 := head["salt"]
	signHex, ok3 := head["sign"]
	if !ok1 || !ok2 || !ok3 {
		return nil, false
	}

	pubk, err1 := hex.DecodeString(pubkHex)
	salt, err2 := hex.DecodeString(saltHex)
	sign, err3 := hex.DecodeString(signHex)
	if err1 != nil || err2 != nil || err3 != nil {
		return nil, false
	}

	pubKey := ed25519.PublicKey(pubk)

	hash := hashing.NewHMACHasher(salt, body).ToBytes()
	return pubKey, ed25519.Verify(pubKey, hash, sign)
}

func hasNotGraphicCharacters(pS string) bool {
	for _, c := range pS {
		if !unicode.IsGraphic(c) {
			return true
		}
	}
	return false
}
