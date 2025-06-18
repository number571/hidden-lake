package database

import (
	"errors"
	"testing"

	"github.com/number571/go-peer/pkg/storage/database"
)

func TestVoidKVDatabase(t *testing.T) {
	db := NewVoidKVDatabase()
	if err := db.Set([]byte("aaa"), []byte("bbb")); err != nil {
		t.Fatal(err)
	}
	_, err := db.Get([]byte("aaa"))
	if err == nil {
		t.Fatal("success get value from void database")
	}
	if !errors.Is(err, database.ErrNotFound) {
		t.Fatal("got unsupported error")
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}
