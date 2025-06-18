package response

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/payload/joiner"
)

const (
	tcResponse = `{"code":200,"head":{"key1":"value1","key2":"value2","key3":"value3"},"body":"aGVsbG8sIHdvcmxkIQ=="}`
	tcBody     = "hello, world!"
)

var (
	tgHead = map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SResponseError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestInvalidResponse(t *testing.T) {
	t.Parallel()

	if _, err := LoadResponse([]byte{123}); err == nil {
		t.Fatal("success load invalid response bytes")
	}

	bytesJoiner := joiner.NewBytesJoiner32([][]byte{
		{byte(123)},
		{byte(111)},
	})
	if _, err := LoadResponse(bytesJoiner); err == nil {
		t.Fatal("success load invalid response bytes joiner")
	}

	if _, err := LoadResponse("123"); err == nil {
		t.Fatal("success load invalid response string")
	}

	if _, err := LoadResponse(struct{}{}); err == nil {
		t.Fatal("success load invalid response type")
	}
}

func TestResponse(t *testing.T) {
	t.Parallel()

	resp := NewResponseBuilder().
		WithCode(200).
		WithHead(tgHead).
		WithBody([]byte(tcBody)).
		Build()

	resp1, err := LoadResponse(resp.ToBytes())
	if err != nil {
		t.Fatal(err)
	}

	respStr := resp.ToString()
	if respStr != tcResponse {
		t.Fatal("string response is invalid")
	}

	resp2, err := LoadResponse(respStr)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, resp)
	testResponse(t, resp1)
	testResponse(t, resp2)
}

func testResponse(t *testing.T, resp IResponse) {
	if resp.GetCode() != 200 {
		t.Fatal("resp code is invalid")
	}
	if !bytes.Equal(resp.GetBody(), []byte(tcBody)) {
		t.Fatal("resp body is invalid")
	}
	if len(resp.GetHead()) != 3 {
		t.Fatal("resp head size is invalid")
	}

	for k, v := range resp.GetHead() {
		v1, ok := tgHead[k]
		if !ok {
			t.Fatal("undefined value in orig head")
		}
		if v1 != v {
			t.Fatal("resp head value is invalid")
		}
	}
}
