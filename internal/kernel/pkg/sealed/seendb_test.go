package sealed

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/storage/database"
)

func testSeenKey() []byte {
	return random.NewRandom().GetBytes(cSeenKeySize)
}

func TestSeenDB(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "seen.db")
	key := testSeenKey()
	hash := random.NewRandom().GetBytes(32)

	db, err := NewSeenDB(path, key)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Get(hash); !errors.Is(err, database.ErrNotFound) {
		t.Fatal("unexpected value before set")
	}
	if err := db.Set(hash, []byte{}); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Get(hash); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Get(hash); err == nil {
		t.Fatal("get after close")
	}

	raw, err := os.ReadFile(path) //nolint:gosec
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(raw, hash) {
		t.Fatal("plain hash stored on disk")
	}

	// same key: record is found after reopen
	db, err = NewSeenDB(path, key)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Get(hash); err != nil {
		t.Fatal(err)
	}
	if err := db.Del(hash); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Get(hash); !errors.Is(err, database.ErrNotFound) {
		t.Fatal("value exists after del")
	}
	if err := db.Set(hash, []byte{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	// another key is rejected: a wrong seen_key or a swapped file must not
	// silently pass every anti-replay lookup.
	if _, err := NewSeenDB(path, testSeenKey()); !errors.Is(err, ErrSeenDBKeyMismatch) {
		t.Fatalf("expected key-mismatch error, got %v", err)
	}
}

func TestSeenDBKeyMismatchDetected(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "seen.db")
	key := testSeenKey()
	db, err := NewSeenDB(path, key)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	// Reopen with the right key works; a different key is rejected.
	db2, err := NewSeenDB(path, key)
	if err != nil {
		t.Fatal(err)
	}
	_ = db2.Close()
	if _, err := NewSeenDB(path, testSeenKey()); !errors.Is(err, ErrSeenDBKeyMismatch) {
		t.Fatalf("expected ErrSeenDBKeyMismatch, got %v", err)
	}
}

func TestSeenDBInvalidKey(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "seen.db")
	if _, err := NewSeenDB(path, []byte{1, 2, 3}); !errors.Is(err, ErrInvalidSeenKey) {
		t.Fatal("invalid seen key accepted")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("database created with invalid key")
	}

	if _, err := NewSeenDB(t.TempDir(), testSeenKey()); !errors.Is(err, ErrOpenSeenDB) {
		t.Fatal("open directory as database")
	}
}

func TestSeenDBLocked(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "seen.db")
	db, err := NewSeenDB(path, testSeenKey())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	start := time.Now()
	if _, err := NewSeenDB(path, testSeenKey()); !errors.Is(err, ErrOpenSeenDB) {
		t.Fatal("locked database opened twice")
	}
	if d := time.Since(start); d > 3*time.Second {
		t.Fatalf("open of locked database took %s", d)
	}
}

func TestSeenDBNoParentDir(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "missing", "seen.db")
	if _, err := NewSeenDB(path, testSeenKey()); !errors.Is(err, ErrSeenDBDirNotFound) {
		t.Fatal("missing parent directory is not reported")
	}
	if _, err := os.Stat(filepath.Dir(path)); !os.IsNotExist(err) {
		t.Fatal("parent directory created")
	}
}

func TestSeenDBParentIsFile(t *testing.T) {
	t.Parallel()

	parent := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(parent, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewSeenDB(filepath.Join(parent, "seen.db"), testSeenKey()); !errors.Is(err, ErrSeenDBDirNotFound) {
		t.Fatal("parent file is not reported")
	}
}
