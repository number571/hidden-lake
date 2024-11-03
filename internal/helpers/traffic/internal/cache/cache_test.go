package cache

import (
	"bytes"
	"testing"
)

func TestLRUCache(t *testing.T) {
	t.Parallel()

	cache := NewLRUCache(2)

	if ok := cache.Set([]byte("key"), []byte("value")); !ok {
		t.Error("cache: !ok = set(key, value)")
		return
	}

	if _, ok := cache.Get([]byte("key")); !ok {
		t.Error("cache: !ok = get(key)")
		return
	}

	if cache.GetIndex() != 1 {
		t.Error("cache: get_index() != 1")
		return
	}

	if k, ok := cache.GetKey(0); !ok || !bytes.Equal([]byte("key"), k) {
		t.Error("cache: !ok = get_key(0)")
		return
	}
}

func TestVoidLRUCache(t *testing.T) {
	t.Parallel()

	voidCache := NewLRUCache(0)

	if ok := voidCache.Set([]byte("key"), []byte("value")); !ok {
		t.Error("void_cache: !ok = set(key, value)")
		return
	}

	if _, ok := voidCache.Get([]byte("key")); ok {
		t.Error("void_cache: ok = get(key)")
		return
	}

	if voidCache.GetIndex() != 0 {
		t.Error("void_cache: get_index() != 0")
		return
	}

	if _, ok := voidCache.GetKey(0); ok {
		t.Error("void_cache: ok = get_key(0)")
		return
	}
}
