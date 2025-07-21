// nolint: errcheck, gosec
package database

import (
	"bytes"
	"crypto/ed25519"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SDatabaseError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	db := &sDatabase{
		fKey: [3][]byte{
			[]byte("abcdefghijklmnopqrstuvwxyz123456"),
			[]byte("1234567890abcdefghijklmnopqrstuv"),
			[]byte(""),
		},
	}

	msg := []byte("hello, world!")
	encMsg := db.encryptBytes(msg)

	decMsg, err := db.decryptBytes(encMsg)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(msg, decMsg) {
		t.Fatal("decrypt failed")
	}

	if _, err := db.decryptBytes([]byte{111}); err == nil {
		t.Fatal("success decrypt bytes (1)")
	}

	randBytes := random.NewRandom().GetBytes(hashing.CHasherSize + symmetric.CCipherBlockSize)
	if _, err := db.decryptBytes(randBytes); err == nil {
		t.Fatal("success decrypt bytes (2)")
	}
}

func TestEncodeDecode(t *testing.T) {
	t.Parallel()

	db := &sDatabase{
		fKey: [3][]byte{
			[]byte("abcdefghijklmnopqrstuvwxyz123456"),
			[]byte("1234567890abcdefghijklmnopqrstuv"),
			[]byte(""),
		},
	}

	seed := []byte("111defghijklmnopqrstuvwxyz123456")
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)

	msg := SMessage{
		FSendTime: time.Now(),
		FSender:   pubKey,
		FMessage:  "hello, world!",
	}

	encMsg := db.messageToBytes(msg)
	decMsg, err := db.bytesToMessage(encMsg)
	if err != nil {
		t.Fatal(err)
	}
	if decMsg.FSendTime.Format(time.DateTime) != msg.FSendTime.Format(time.DateTime) {
		t.Fatal("equal send time")
	}
	if !decMsg.FSender.Equal(msg.FSender) {
		t.Fatal("equal sender")
	}
	if decMsg.FMessage != msg.FMessage {
		t.Fatal("equal message")
	}
}

func TestDatabase(t *testing.T) {
	t.Parallel()

	const dbfile = "testdata/test.db"

	os.Remove(dbfile)
	defer os.Remove(dbfile)

	db, err := NewDatabase(
		dbfile,
		[3][]byte{
			[]byte("abcdefghijklmnopqrstuvwxyz123456"),
			[]byte("1234567890abcdefghijklmnopqrstuv"),
			[]byte("67890abcdefghijklmnopqrstuv12345"),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	seed := []byte("111defghijklmnopqrstuvwxyz123456")
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)

	chanPK := asymmetric.NewPrivKey().GetPubKey()

	for i := range 10 {
		msg := SMessage{
			FSendTime: time.Now(),
			FSender:   pubKey,
			FMessage:  fmt.Sprintf("hello, world! (%d)", i),
		}
		if err := db.Insert(chanPK, msg); err != nil {
			t.Fatal(err)
		}
	}

	msgs, err := db.Select(chanPK, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 10 {
		t.Fatal("len msgs != 10")
	}

	for i, msg := range msgs {
		if msg.FMessage != fmt.Sprintf("hello, world! (%d)", i) {
			t.Fatal("equal message")
		}
	}
}
