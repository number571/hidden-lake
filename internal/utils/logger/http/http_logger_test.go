package http

import (
	"context"
	"net/http"
	"testing"
)

const (
	tcService = "TST"
	tcFmtLog  = "service=TST method=GET path=/api/index conn=127.0.0.1:55555 message=hello_world"
)

func TestPanicLogger(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()

	logFunc := GetLogFunc()
	_ = logFunc("_")
}

func TestLogger(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"http://localhost:8080/api/index",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1:55555"

	logBuilder := NewLogBuilder(tcService, req).WithMessage("hello_world")
	logFunc := GetLogFunc()

	if l := logFunc(logBuilder); l != tcFmtLog {
		t.Fatal("got invalid format")
	}

	logGetter := logBuilder.Build()
	if logGetter.GetConn() != "127.0.0.1:55555" {
		t.Fatal("got conn != conn")
	}

	if logGetter.GetMessage() != "hello_world" {
		t.Fatal("got message != message")
	}

	if logGetter.GetMethod() != "GET" {
		t.Fatal("got method != method")
	}

	if logGetter.GetPath() != "/api/index" {
		t.Fatal("got path != path")
	}

	if logGetter.GetService() != "TST" {
		t.Fatal("got service != service")
	}
}
