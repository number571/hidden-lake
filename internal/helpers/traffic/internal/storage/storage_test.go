package storage

import (
	"bytes"
	"errors"
	"strconv"
	"testing"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/cache"
	hlt_database "github.com/number571/hidden-lake/internal/helpers/traffic/internal/database"
)

const (
	tcHead        = uint64(123)
	tcBody        = "hello, world!"
	tcNetworkKey  = "_"
	tcWorkSize    = 10
	tcCapacity    = 16
	tcMessageSize = (8 << 10)
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SStorageError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestStorageLoad(t *testing.T) {
	t.Parallel()

	db := NewMessageStorage(
		net_message.NewSettings(&net_message.SSettings{
			FWorkSizeBits: tcWorkSize,
			FNetworkKey:   tcNetworkKey,
		}),
		hlt_database.NewVoidKVDatabase(),
		cache.NewLRUCache(tcCapacity),
	)

	if _, err := db.Load([]byte("abc")); err == nil {
		t.Error("success load not exist message (incorrect)")
		return
	}

	hash := hashing.NewHasher([]byte{123}).ToBytes()
	_, errLoad := db.Load(hash)
	if errLoad == nil {
		t.Error("success load not exist message (hash)")
		return
	}

	if !errors.Is(errLoad, ErrMessageIsNotExist) {
		t.Error("got incorrect error type (load)")
		return
	}
}

func TestStorageHashes(t *testing.T) {
	t.Parallel()
	const messagesCapacity = 3

	db := NewMessageStorage(
		net_message.NewSettings(&net_message.SSettings{
			FWorkSizeBits: tcWorkSize,
			FNetworkKey:   tcNetworkKey,
		}),
		hlt_database.NewVoidKVDatabase(),
		cache.NewLRUCache(messagesCapacity),
	)

	cl := client.NewClient(
		asymmetric.NewPrivKey(),
		tcMessageSize,
	)

	pushHashes := make([][]byte, 0, messagesCapacity+1)
	for i := 0; i < messagesCapacity+1; i++ {
		msg, err := newNetworkMessageWithData(cl, tcNetworkKey, strconv.Itoa(i))
		if err != nil {
			t.Error(err)
			return
		}
		if err := db.Push(msg); err != nil {
			t.Error(err)
			return
		}
		if db.Pointer() != uint64(i+1)%messagesCapacity {
			t.Error("got invalid pointer")
			return
		}
		pushHashes = append(pushHashes, msg.GetHash())
	}

	for i := uint64(0); i < messagesCapacity+1; i++ {
		hash, err := db.Hash(i)
		if err != nil {
			break
		}
		if bytes.Equal(hash, pushHashes[0]) {
			t.Error("hash not overwritten")
			return
		}
		index := (2 + i) % (messagesCapacity)
		if !bytes.Equal(hash, pushHashes[1:][index]) {
			t.Error("got invalid hash")
			return
		}
	}
}

func TestStoragePush(t *testing.T) {
	t.Parallel()

	db := NewMessageStorage(
		net_message.NewSettings(&net_message.SSettings{
			FWorkSizeBits: tcWorkSize,
			FNetworkKey:   tcNetworkKey,
		}),
		hlt_database.NewVoidKVDatabase(),
		cache.NewLRUCache(1),
	)

	clTest := client.NewClient(
		asymmetric.NewPrivKey(),
		tcMessageSize,
	)

	msgTest, err := newNetworkMessage(clTest, "some-another-key")
	if err != nil {
		t.Error(err)
		return
	}

	if err := db.Push(msgTest); err == nil {
		t.Error("success push message with difference setting")
		return
	}

	cl := client.NewClient(
		asymmetric.NewPrivKey(),
		tcMessageSize,
	)

	msg1, err := newNetworkMessage(cl, tcNetworkKey)
	if err != nil {
		t.Error(err)
		return
	}

	if err := db.Push(msg1); err != nil {
		t.Error(err)
		return
	}

	errPush := db.Push(msg1)
	if errPush == nil {
		t.Error("success push duplicate")
		return
	}

	if !errors.Is(errPush, ErrMessageIsExist) {
		t.Error("got incorrect error type (push)")
		return
	}

	msg2, err := newNetworkMessage(cl, tcNetworkKey)
	if err != nil {
		t.Error(err)
		return
	}

	if err := db.Push(msg2); err != nil {
		t.Error(err)
		return
	}
}

func TestStorage(t *testing.T) {
	t.Parallel()

	db := NewMessageStorage(
		net_message.NewSettings(&net_message.SSettings{
			FWorkSizeBits: tcWorkSize,
			FNetworkKey:   tcNetworkKey,
		}),
		hlt_database.NewVoidKVDatabase(),
		cache.NewLRUCache(4),
	)

	cl := client.NewClient(
		asymmetric.NewPrivKey(),
		tcMessageSize,
	)

	putHashes := make([][]byte, 0, 3)
	for i := 0; i < 3; i++ {
		msg, err := newNetworkMessage(cl, tcNetworkKey)
		if err != nil {
			t.Error(err)
			return
		}
		if err := db.Push(msg); err != nil {
			t.Error(err)
			return
		}
		putHashes = append(putHashes, msg.GetHash())
	}

	getHashes := make([][]byte, 0, 3)
	for i := uint64(0); ; i++ {
		hash, err := db.Hash(i)
		if err != nil {
			break
		}
		getHashes = append(getHashes, hash)
	}

	if len(getHashes) != 3 {
		t.Error("len getHashes != 3")
		return
	}

	for i := range getHashes {
		if !bytes.Equal(getHashes[i], putHashes[i]) {
			t.Errorf("getHashes[%d] != putHashes[%d]", i, i)
			return
		}
	}

	for _, getHash := range getHashes {
		loadNetMsg, err := db.Load(getHash)
		if err != nil {
			t.Error(err)
			return
		}

		msgHash := loadNetMsg.GetHash()
		if !bytes.Equal(getHash, msgHash) {
			t.Errorf("getHash[%s] != msgHash[%s]", getHash, msgHash)
			return
		}

		pubKey, decMsg, err := cl.DecryptMessage(
			asymmetric.NewMapPubKeys(cl.GetPrivKey().GetPubKey()),
			loadNetMsg.GetPayload().GetBody(),
		)
		if err != nil {
			t.Error(err)
			return
		}

		if !bytes.Equal(pubKey.ToBytes(), cl.GetPrivKey().GetPubKey().ToBytes()) {
			t.Error("load public key != init public key")
			return
		}

		pl := payload.LoadPayload64(decMsg)
		if pl.GetHead() != tcHead {
			t.Error("load msg head != init head")
			return
		}

		if !bytes.Equal(pl.GetBody(), []byte(tcBody)) {
			t.Error("load msg body != init body")
			return
		}
	}
}

func newNetworkMessageWithData(cl client.IClient, networkKey, data string) (net_message.IMessage, error) {
	msg, err := cl.EncryptMessage(
		cl.GetPrivKey().GetPubKey(),
		payload.NewPayload64(tcHead, []byte(data)).ToBytes(),
	)
	if err != nil {
		return nil, err
	}
	netMsg := net_message.NewMessage(
		net_message.NewConstructSettings(&net_message.SConstructSettings{
			FSettings: net_message.NewSettings(&net_message.SSettings{
				FNetworkKey:   networkKey,
				FWorkSizeBits: tcWorkSize,
			}),
		}),
		payload.NewPayload32(0, msg),
	)
	return netMsg, nil
}

func newNetworkMessage(cl client.IClient, networkKey string) (net_message.IMessage, error) {
	return newNetworkMessageWithData(cl, networkKey, tcBody)
}
