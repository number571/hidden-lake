package request

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/payload/joiner"
)

const (
	tcHost   = "test_host"
	tcPath   = "test_path"
	tcMethod = "test_method"
)

var (
	tgHead = map[string]string{
		"test_header1": "test_value1",
		"test_header2": "test_value2",
		"test_header3": "test_value3",
	}
	tgBody     = []byte("test_data")
	tgBRequest = `{"method":"test_method","host":"test_host","path":"test_path","head":{"test_header1":"test_value1","test_header2":"test_value2","test_header3":"test_value3"},"body":"dGVzdF9kYXRh"}`
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SRequestError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestInvalidRequest(t *testing.T) {
	t.Parallel()

	if _, err := LoadRequest([]byte{123}); err == nil {
		t.Fatal("success load invalid request bytes")
	}

	bytesJoiner := joiner.NewBytesJoiner32([][]byte{
		{byte(123)},
		{byte(111)},
	})
	if _, err := LoadRequest(bytesJoiner); err == nil {
		t.Fatal("success load invalid request bytes joiner")
	}

	if _, err := LoadRequest("123"); err == nil {
		t.Fatal("success load invalid request string")
	}

	if _, err := LoadRequest(struct{}{}); err == nil {
		t.Fatal("success load invalid request type")
	}
}

func TestRequest(t *testing.T) {
	t.Parallel()

	request := NewRequestBuilder().
		WithMethod(tcMethod).
		WithHost(tcHost).
		WithPath(tcPath).
		WithHead(tgHead).
		WithBody(tgBody).
		Build()

	if request.GetHost() != tcHost {
		t.Fatal("host is not equals")
	}

	if request.GetPath() != tcPath {
		t.Fatal("path is not equals")
	}

	if request.GetMethod() != tcMethod {
		t.Fatal("method is not equals")
	}

	for k, v := range request.GetHead() {
		v1, ok := tgHead[k]
		if !ok {
			t.Fatalf("header undefined '%s'", k)
		}
		if v != v1 {
			t.Fatalf("header is invalid '%s'", v1)
		}
	}

	if !bytes.Equal(request.GetBody(), tgBody) {
		t.Fatal("body is not equals")
	}
}

func TestLoadRequest(t *testing.T) {
	t.Parallel()

	brequest := NewRequestBuilder().
		WithMethod(tcMethod).
		WithHost(tcHost).
		WithPath(tcPath).
		WithHead(tgHead).
		WithBody(tgBody).
		Build().
		ToBytes()

	request1, err := LoadRequest(brequest)
	if err != nil {
		t.Fatal(err)
	}

	request2, err := LoadRequest(tgBRequest)
	if err != nil {
		t.Fatal(err)
	}

	reqStr := request2.ToString()
	if reqStr != tgBRequest {
		t.Fatal("string request is invalid")
	}

	if request1.GetHost() != request2.GetHost() {
		t.Fatal("host is not equals")
	}

	if request1.GetPath() != request2.GetPath() {
		t.Fatal("path is not equals")
	}

	if request1.GetMethod() != request2.GetMethod() {
		t.Fatal("method is not equals")
	}

	for k, v := range request1.GetHead() {
		v1, ok := request2.GetHead()[k]
		if !ok {
			t.Fatalf("header undefined '%s'", k)
		}
		if v != v1 {
			t.Fatalf("header is invalid '%s'", v1)
		}
	}

	if !bytes.Equal(request1.GetBody(), request2.GetBody()) {
		t.Fatal("body is not equals")
	}
}
